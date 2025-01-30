package main

import (
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
	"time"
)

// MultiPaxosService is the main server struct
type MultiPaxosService struct {
	multipaxos.UnimplementedMultiPaxosServiceServer

	multiPaxosClient *MultiPaxosClient
	peersCluster     *PeersCluster
	synchronizer     *Synchronizer

	// A map of slot -> PaxosInstance
	instances      map[int64]*PaxosInstance
	currentSlot    int64
	mu             sync.RWMutex
	slotLock       sync.RWMutex
	instancesMu    sync.RWMutex
	recoveryDoneCh chan struct{}
}

type PaxosInstance struct {
	mu            sync.RWMutex
	promisedRound int32
	acceptedRound int32
	acceptedValue *multipaxos.ScooterEvent
	committed     bool
}

func NewMultiPaxosService(synchronizer *Synchronizer, multiPaxosClient *MultiPaxosClient, peersCluster *PeersCluster, recoveryDoneCh chan struct{}) *MultiPaxosService {
	s := &MultiPaxosService{
		synchronizer:     synchronizer,
		multiPaxosClient: multiPaxosClient,
		peersCluster:     peersCluster,
		instances:        make(map[int64]*PaxosInstance),
		currentSlot:      1,
		mu:               sync.RWMutex{},
		slotLock:         sync.RWMutex{},
		recoveryDoneCh:   recoveryDoneCh,
	}

	return s
}

func (s *MultiPaxosService) GetSnapshot(ctx context.Context, req *multipaxos.GetSnapshotRequest) (*multipaxos.GetSnapshotResponse, error) {
	// Return the current snapshot state
	s.synchronizer.mu.RLock()
	defer s.synchronizer.mu.RUnlock()

	// Convert internal scooter state to proto message format
	protoState := make(map[string]*multipaxos.Scooter)
	for id, scooter := range s.synchronizer.snapshot.state {
		protoState[id] = &multipaxos.Scooter{
			Id:                   scooter.Id,
			IsAvailable:          scooter.IsAvailable,
			TotalDistance:        scooter.TotalDistance,
			CurrentReservationId: scooter.CurrentReservationId,
		}
	}

	return &multipaxos.GetSnapshotResponse{
		LastGoodSlot: s.synchronizer.snapshot.lastGoodSlot,
		State:        protoState,
	}, nil
}

func (s *MultiPaxosService) recoverFromSnapshot() error {
	latestSnapshot, err := s.multiPaxosClient.GetLatestSnapshot()
	if err != nil {
		return err
	}
	log.Printf("Recovered snapshot is: %v", latestSnapshot)

	// Update local state with recovered snapshot
	s.synchronizer.mu.Lock()
	defer s.synchronizer.mu.Unlock()

	// Convert proto state back to internal format
	recoveredState := make(map[string]*Scooter)
	for id, protoScooter := range latestSnapshot.State {
		scooter := &Scooter{
			Id:                   protoScooter.Id,
			IsAvailable:          protoScooter.IsAvailable,
			TotalDistance:        protoScooter.TotalDistance,
			CurrentReservationId: protoScooter.CurrentReservationId,
		}
		recoveredState[id] = scooter
		s.synchronizer.state.state[scooter.Id] = scooter
		s.synchronizer.snapshot.state[scooter.Id] = scooter
		log.Printf("Recovered scooter: %v", recoveredState[id])
	}

	s.synchronizer.snapshot.lastGoodSlot = latestSnapshot.LastGoodSlot
	s.synchronizer.state.lastGoodSlot = latestSnapshot.LastGoodSlot

	log.Printf("Snapshot recovered is: %v, lastGoodSlot: %v", s.synchronizer.snapshot.state, s.synchronizer.snapshot.lastGoodSlot)

	// Update current slot to start after the snapshot
	s.mu.Lock()
	s.currentSlot = latestSnapshot.LastGoodSlot + 1
	s.mu.Unlock()

	return nil
}

func (s *MultiPaxosService) recoverInstanceFromPeers(slot int64) (*PaxosInstance, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the list of peers from etcd
	peers := s.peersCluster.GetPeersList()
	log.Printf("Attempting to recover slot %d from %d peers", slot, len(peers))

	for _, peer := range peers {
		// Skip own address to avoid self-connection
		if peer == myCandidateInfo {
			continue
		}

		// Create connection with non-blocking dial
		conn, err := grpc.Dial(peer,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(2*time.Second),
			grpc.WithBlock(),                 // Wait for connection with timeout
			grpc.WithReturnConnectionError()) // Return concrete error

		if err != nil {
			log.Printf("Failed to connect to peer %s: %v", peer, err)
			continue
		}

		// Ensure connection is closed
		defer conn.Close()

		// Check connection state before proceeding
		state := conn.GetState()
		if state != connectivity.Ready {
			log.Printf("Connection to peer %s not ready, state: %s", peer, state)
			continue
		}
		client := multipaxos.NewMultiPaxosServiceClient(conn)

		// Use a shorter timeout for the RPC call itself
		callCtx, callCancel := context.WithTimeout(ctx, 1*time.Second)
		defer callCancel()

		resp, err := client.GetPaxosState(callCtx, &multipaxos.GetPaxosStateRequest{
			Slot: slot,
		})

		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				log.Printf("Timeout getting state from peer %s for slot %d", peer, slot)
			} else {
				log.Printf("Error getting state from peer %s for slot %d: %v", peer, slot, err)
			}
			continue
		}

		instance := &PaxosInstance{
			promisedRound: resp.GetPromisedRound(),
			acceptedRound: resp.GetAcceptedRound(),
			acceptedValue: resp.GetAcceptedValue(),
			committed:     resp.GetCommitted(),
		}

		log.Printf("Successfully recovered slot %d from peer %s", slot, peer)
		return instance, nil
	}

	log.Printf("No peers had instance for slot %d", slot)
	return nil, nil
}

func (s *MultiPaxosService) recoverAllInstances() {
	if !HasAnyReadyPeers(s.peersCluster) {
		log.Printf("No peers ready for recovery, starting fresh")
		s.mu.Lock()
		s.currentSlot = 1
		s.mu.Unlock()
		return
	}

	// First try to recover from snapshot
	log.Println("[Recovery] Starting snapshot recovery")
	err := s.recoverFromSnapshot()
	var startingSlot int64 = 1
	if err != nil {
		log.Printf("[Recovery] Failed to recover from snapshot: %v. Falling back to instance recovery from slot 1", err)
	} else {
		startingSlot = s.synchronizer.snapshot.lastGoodSlot + 1
		log.Printf("[Recovery] Successfully recovered snapshot up to slot %d. Continuing recovery from slot %d",
			s.synchronizer.snapshot.lastGoodSlot, startingSlot)
	}

	// Recover individual instances starting from the slot after the snapshot
	log.Printf("[Recovery] Starting instance recovery from slot %d", startingSlot)
	var maxSlot int64 = s.synchronizer.snapshot.lastGoodSlot

	s.instancesMu.Lock()
	defer s.instancesMu.Unlock()

	pendingSlots := make([]int64, 0) // Track slots that fail to recover

	for slot := startingSlot; ; slot++ {
		instance, err := s.recoverInstanceFromPeers(slot)
		if err != nil {
			log.Printf("[Recovery] Failed to recover instance %d: %v. Adding to pending list.", slot, err)
			pendingSlots = append(pendingSlots, slot)
			continue
		}

		if instance == nil {
			// No more instances to recover
			log.Printf("[Recovery] No instance data available for slot %d. Stopping recovery.", slot)
			break
		}

		// Update recovered instance
		s.instances[slot] = instance
		val := instance.acceptedValue

		// Update state with committed value
		s.synchronizer.mu.Lock()
		s.synchronizer.approvedEventLog[slot] = val
		s.synchronizer.updateState(s.synchronizer.snapshot, slot)
		s.synchronizer.mu.Unlock()

		if slot > maxSlot {
			maxSlot = slot
		}
	}

	log.Printf("[Recovery] Instance recovery completed. Starting pending slot recovery.")

	// Retry pending slots with exponential backoff
	for _, slot := range pendingSlots {
		s.retryRecoverSlot(slot)
	}

	// Set the next available slot
	s.mu.Lock()
	if maxSlot > 0 {
		s.currentSlot = maxSlot + 1
		log.Printf("[Recovery] Recovery completed up to slot %d. Next available slot is %d.", maxSlot, s.currentSlot)
	} else {
		s.currentSlot = startingSlot
		log.Printf("[Recovery] No additional slots recovered after snapshot. Starting from slot %d.", startingSlot)
	}
	s.mu.Unlock()

	// Update snapshot with any additionally recovered state
	s.synchronizer.mu.Lock()
	s.synchronizer.snapshot.lastGoodSlot = maxSlot
	s.synchronizer.mu.Unlock()

	log.Printf("[Recovery] State is: %v.", s.synchronizer.state)

	log.Printf("[Recovery] Recovery process completed successfully.")
}

func (s *MultiPaxosService) retryRecoverSlot(slot int64) {
	maxRetries := 3
	backoff := time.Second

	for retries := 0; retries < maxRetries; retries++ {
		log.Printf("[Recovery] Retrying recovery for slot %d. Attempt %d/%d.", slot, retries+1, maxRetries)
		instance, err := s.recoverInstanceFromPeers(slot)
		if err == nil && instance != nil {
			// Successfully recovered instance
			s.instancesMu.Lock()
			s.instances[slot] = instance
			s.instancesMu.Unlock()
			log.Printf("[Recovery] Successfully recovered slot %d on retry %d.", slot, retries+1)
			return
		}

		// Wait before retrying
		time.Sleep(backoff)
		backoff *= 2
	}

	log.Printf("[Recovery] Failed to recover slot %d after %d retries.", slot, maxRetries)
}

func (s *MultiPaxosService) loadInstance(slot int64) (*PaxosInstance, bool) {
	s.instancesMu.RLock()
	defer s.instancesMu.RUnlock()
	instance, ok := s.instances[slot]
	if !ok {
		return nil, false
	}

	instance.mu.RLock()
	defer instance.mu.RUnlock()
	return instance, true
}

func (s *MultiPaxosService) getNewSlot() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentSlot++
	return s.currentSlot
}

// getInstance safely fetches (or creates) the instance for a given slot.
func (s *MultiPaxosService) getInstance(slot int64) *PaxosInstance {
	s.instancesMu.Lock()
	defer s.instancesMu.Unlock()

	if _, exists := s.instances[slot]; !exists {
		s.instances[slot] = &PaxosInstance{
			mu:            sync.RWMutex{},
			promisedRound: 0,
			acceptedRound: 0,
			acceptedValue: nil,
			committed:     false,
		}
	}
	return s.instances[slot]
}

func startPaxosServer(stopCh chan struct{}, port string, multiPaxosService *MultiPaxosService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	go func() {
		multipaxos.RegisterMultiPaxosServiceServer(s, multiPaxosService)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Printf("Multipaxos gRPC server listening to port " + port)
	<-stopCh
	log.Println("Shutting down MultiPaxos gRPC server...")

}

// gRPC handlers ---------------------- start
// Prepare handler
func (s *MultiPaxosService) Prepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
	log.Printf("Received Prepare message from %s, slot %d, round %d ", req.Id, req.Slot, req.Round)

	slot := req.GetSlot()
	round := req.GetRound()
	instance := s.getInstance(slot)

	instance.mu.Lock()
	defer instance.mu.Unlock()

	if round > instance.promisedRound {
		// Update promised round
		instance.promisedRound = round

		// Return the acceptedRound and acceptedValue (if any)
		return &multipaxos.PrepareResponse{
			Ok:            true,
			AcceptedRound: instance.acceptedRound,
			AcceptedValue: instance.acceptedValue,
		}, nil
	}

	return &multipaxos.PrepareResponse{
		Ok: false,
	}, nil
}

// Accept handler
func (s *MultiPaxosService) Accept(ctx context.Context, req *multipaxos.AcceptRequest) (*multipaxos.AcceptResponse, error) {
	log.Printf("Received Accept message from %s, slot %d, round %d ", req.Id, req.Slot, req.Round)
	slot := req.GetSlot()
	round := req.GetRound()
	val := req.GetValue()
	instance := s.getInstance(slot)

	instance.mu.Lock()
	defer instance.mu.Unlock()

	if round >= instance.promisedRound {
		// Accept the proposal
		instance.acceptedRound = round
		instance.acceptedValue = val
		return &multipaxos.AcceptResponse{Ok: true}, nil
	}

	return &multipaxos.AcceptResponse{Ok: false}, nil
}

// Commit handler
func (s *MultiPaxosService) Commit(ctx context.Context, req *multipaxos.CommitRequest) (*multipaxos.CommitResponse, error) {
	log.Printf("Received Commit message from %s, slot %d, round %d ", req.Id, req.Slot, req.Round)
	slot := req.GetSlot()
	round := req.GetRound()
	val := req.GetValue()
	instance := s.getInstance(slot)

	instance.mu.Lock()
	defer instance.mu.Unlock()

	// We only commit if the round and value match what we've accepted
	if round == instance.acceptedRound /*TODO: do we need to compare values also? && val == instance.acceptedValue*/ {
		instance.committed = true
		log.Printf("Slot %d committed value %s in round %d", slot, val.ScooterId, round)
		s.synchronizer.updateStateWithCommited(slot, val)
		return &multipaxos.CommitResponse{Ok: true}, nil
	}

	return &multipaxos.CommitResponse{Ok: false}, nil
}

// Commit handler
func (s *MultiPaxosService) TriggerLeader(ctx context.Context, req *multipaxos.TriggerRequest) (*multipaxos.TriggerResponse, error) {
	log.Printf("Received TriggerPrepare from %s. event: %v", req.MemberId, req.Event)
	s.synchronizer.myPendingEvents.Enqueue(req.Event.EventId, req.Event)
	s.start()
	return &multipaxos.TriggerResponse{Ok: true}, nil
}

// GetPaxosState handler
func (s *MultiPaxosService) GetPaxosState(ctx context.Context, req *multipaxos.GetPaxosStateRequest) (*multipaxos.GetPaxosStateResponse, error) {
	slot := req.GetSlot()
	instanceState, ok := s.loadInstance(slot)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "instance not found for slot %d", slot)
	}

	return &multipaxos.GetPaxosStateResponse{
		PromisedRound: instanceState.promisedRound,
		AcceptedRound: instanceState.acceptedRound,
		AcceptedValue: instanceState.acceptedValue,
		Committed:     instanceState.committed,
	}, nil
}

// GetLastGoodSlot handler
func (s *MultiPaxosService) GetLastGoodSlot(ctx context.Context, req *multipaxos.GetLastGoodSlotRequest) (*multipaxos.GetLastGoodSlotResponse, error) {
	// TODO : add locks
	slot := s.synchronizer.state.lastGoodSlot
	log.Printf("Received GetLastGoodSlot call, responding with %d", slot)
	return &multipaxos.GetLastGoodSlotResponse{
		MemberId:     req.MemberId,
		LastGoodSlot: slot,
	}, nil
}

// gRPC handlers ---------------------- end

func (s *MultiPaxosService) createNextPaxosInstance(event *multipaxos.ScooterEvent) (*PaxosInstance, int64) {
	s.slotLock.Lock()
	defer s.slotLock.Unlock()
	for {
		instance, exists := s.instances[s.currentSlot]
		if exists {
			if instance.committed {
				log.Printf("Paxos: Slot %d is already committed. Checking slot %d", s.currentSlot, s.currentSlot+1)
				s.currentSlot++
			} else {
				return nil, 0 // current instance was not finished, wait for it by stopping now, and the current instance will start a new paxos instance once finished
			}
		} else {
			instance = &PaxosInstance{acceptedValue: event}
			s.instances[s.currentSlot] = instance
			return instance, s.currentSlot
		}
	}
}
func (s *MultiPaxosService) start() {
	// Wait for recoverIfNeeded to finish, the problem is with using and changing the currentSlot variable,
	// because recovery can change its initial value and this can create a miss
	// this is relevant only to leader, and this method (start) is the only one that uses currentSlot
	<-s.recoveryDoneCh
	if !amILeader() {
		s.multiPaxosClient.TriggerPrepare(s.synchronizer.myPendingEvents.Peek())
		// If I'm not the leader, do nothing or ping the leader
		return
	}

	event := s.synchronizer.myPendingEvents.Peek()
	if event == nil {
		// don't have any events in queue, nothing to work with
		log.Printf("Paxos: Pending events queue is empty, no more events to deliver")
		return
	}
	instance, slot := s.createNextPaxosInstance(event)
	if instance == nil {
		log.Printf("Paxos: there is an existing instance in progress, finishing now, the current instance will trigger a new cycle when finished")
		return
	}

	ctx := context.TODO()
	peers := s.peersCluster.GetPeersList()

	round := 1

	log.Printf("Paxos: Leader proposing value at slot %d, round %d, value %v", slot, round, instance.acceptedValue)
	log.Printf("Paxos: sending prepare message to : %v", peers)
	s.ProposeValue(ctx, slot, instance.acceptedValue /* proposed value */, int32(round), peers)
	s.start()
}

// ProposeValue proposes a new value to a specific slot.
// This would run on the leader node.
func (s *MultiPaxosService) ProposeValue(ctx context.Context, slot int64, proposedValue *multipaxos.ScooterEvent, roundStart int32, peers []string) error {
	round := roundStart

	for {
		// 1. Phase 1 (Prepare)
		acceptedRound, acceptedVal, prepareOKCount := s.multiPaxosClient.doPreparePhase(ctx, slot, round, peers)

		// If fewer than majority responded OK, or got a "no," try next round
		if prepareOKCount <= len(peers)/2 {
			round++
			continue
		}

		// The value to propose in Accept phase
		finalValue := proposedValue
		if acceptedRound != 0 {
			// If a node responded with a previously accepted round, adopt that
			finalValue = acceptedVal
		}

		// 2. Phase 2 (Accept)
		acceptOKCount := s.multiPaxosClient.doAcceptPhase(ctx, slot, round, finalValue, peers)
		if acceptOKCount <= len(peers)/2 {
			// Failed to get majority, increment round & retry
			round++
			continue
		}

		// 3. Phase 3 (Commit)
		// Once majority accepted, we commit
		log.Printf("final value in ProposeValue is %s, accepted round %d", finalValue.ScooterId, acceptedRound)
		s.multiPaxosClient.doCommitPhase(ctx, slot, round, finalValue, peers)

		// Success, break out
		return nil
	}
}
