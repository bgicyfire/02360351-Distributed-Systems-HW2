package main

import (
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"sync"
)

// MultiPaxosService is the main server struct
type MultiPaxosService struct {
	multipaxos.UnimplementedMultiPaxosServiceServer

	// etcd client can be used for leader election/failure detection
	etcdClient       *clientv3.Client
	multiPaxosClient *MultiPaxosClient
	synchronizer     *Synchronizer

	// A map of slot -> PaxosInstance
	instances   map[int32]*PaxosInstance
	currentSlot int32
	mu          sync.RWMutex
	slotLock    sync.RWMutex
	instancesMu sync.RWMutex
}

type PaxosInstance struct {
	mu            sync.RWMutex
	promisedRound int32
	acceptedRound int32
	acceptedValue int32
	committed     bool
}

func NewMultiPaxosService(synchronizer *Synchronizer, etcdClient *clientv3.Client, multiPaxosClient *MultiPaxosClient) *MultiPaxosService {
	s := &MultiPaxosService{
		synchronizer:     synchronizer,
		etcdClient:       etcdClient,
		multiPaxosClient: multiPaxosClient,
		instances:        make(map[int32]*PaxosInstance),
		currentSlot:      0,
		mu:               sync.RWMutex{},
		slotLock:         sync.RWMutex{},
	}

	s.instances[0] = &PaxosInstance{
		promisedRound: 0,
		acceptedRound: 0,
		acceptedValue: 0,
		committed:     false,
		mu:            sync.RWMutex{},
	}

	return s

}

func (s *MultiPaxosService) getNewSlot() int32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentSlot++
	return s.currentSlot
}

// getInstance safely fetches (or creates) the instance for a given slot.
func (s *MultiPaxosService) getInstance(slot int32) *PaxosInstance {
	s.instancesMu.Lock()
	defer s.instancesMu.Unlock()

	if _, exists := s.instances[slot]; !exists {
		s.instances[slot] = &PaxosInstance{
			mu:            sync.RWMutex{},
			promisedRound: 0,
			acceptedRound: 0,
			acceptedValue: 0,
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

// Prepare handler
func (s *MultiPaxosService) Prepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
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
	slot := req.GetSlot()
	round := req.GetRound()
	val := req.GetValue()
	instance := s.getInstance(slot)

	instance.mu.Lock()
	defer instance.mu.Unlock()

	// We only commit if the round and value match what we've accepted
	if round == instance.acceptedRound && val == instance.acceptedValue {
		instance.committed = true
		log.Printf("Slot %d committed value %d in round %d", slot, val, round)
		return &multipaxos.CommitResponse{Ok: true}, nil
	}

	return &multipaxos.CommitResponse{Ok: false}, nil
}

// Commit handler
func (s *MultiPaxosService) TriggerPrepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
	log.Printf("Received TriggerPrepare from %s", req.Id)
	s.start()
	return &multipaxos.PrepareResponse{Ok: true}, nil
}

func (s *MultiPaxosService) start() {
	if !amILeader() {
		s.multiPaxosClient.TriggerPrepare()
		// If I'm not the leader, do nothing or ping the leader
		return
	}

	s.slotLock.Lock()
	slot := s.currentSlot

	instance, exists := s.instances[slot]
	if !exists {
		// create the instance if needed
		instance = &PaxosInstance{}
		s.instances[slot] = instance
	}

	// If the instance at this slot is already committed, move to next slot
	if instance.committed {
		log.Printf("Slot %d is already committed. Moving to slot %d", slot, slot+1)
		s.currentSlot++
		slot = s.currentSlot

		// create a fresh instance for the new slot
		instance = &PaxosInstance{}
		s.instances[slot] = instance
	}

	// Now increment currentSlot so next time we look for a new slot
	s.currentSlot++
	s.slotLock.Unlock()

	ctx := context.TODO()
	peers := fetchOtherServersList(ctx)

	round := 1

	// IMPORTANT: Pass the 'slot' variable to ProposeValue, rather than hardcoding 1
	// e.g. we propose the value "1" at this new 'slot'
	log.Printf("Leader proposing value at slot %d, round %d", slot, round)
	s.ProposeValue(ctx, slot, 1 /* proposed value */, int32(round), peers)

}

// Example function to propose a new value to a specific slot.
// This would run on the leader node.
func (s *MultiPaxosService) ProposeValue(ctx context.Context, slot int32, proposedValue int32, roundStart int32, peers []string) error {
	round := roundStart

	for {
		// 1. Phase 1 (Prepare)
		acceptedRound, acceptedVal, prepareOKCount := s.doPreparePhase(ctx, slot, round, peers)

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
		acceptOKCount := s.doAcceptPhase(ctx, slot, round, finalValue, peers)
		if acceptOKCount <= len(peers)/2 {
			// Failed to get majority, increment round & retry
			round++
			continue
		}

		// 3. Phase 3 (Commit)
		// Once majority accepted, we commit
		s.doCommitPhase(ctx, slot, round, finalValue, peers)

		// Success, break out
		return nil
	}
}

// doPreparePhase sends Prepare to all peers and counts how many OK
// Also returns the highest acceptedRound/value from the responding peers.
func (s *MultiPaxosService) doPreparePhase(ctx context.Context, slot, round int32, peers []string) (int32, int32, int) {
	prepareOKCount := 0
	highestAcceptedRound := int32(0)
	highestAcceptedValue := int32(0)

	// Send to each peer
	for _, peer := range peers {
		// Connect
		conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		client := multipaxos.NewMultiPaxosServiceClient(conn)

		resp, err := client.Prepare(ctx, &multipaxos.PrepareRequest{
			Slot:  slot,
			Round: round,
			Id:    myCandidateInfo, // your ID
		})
		conn.Close()
		if err != nil {
			continue
		}
		if resp.GetOk() {
			prepareOKCount++
			// Track highest acceptedRound/value returned
			if resp.GetAcceptedRound() > highestAcceptedRound {
				highestAcceptedRound = resp.GetAcceptedRound()
				highestAcceptedValue = resp.GetAcceptedValue()
			}
		}
	}

	return highestAcceptedRound, highestAcceptedValue, prepareOKCount
}

// doAcceptPhase sends Accept to all peers with the final chosen value
func (s *MultiPaxosService) doAcceptPhase(ctx context.Context, slot, round, value int32, peers []string) int {
	acceptOKCount := 0
	for _, peer := range peers {
		conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		client := multipaxos.NewMultiPaxosServiceClient(conn)

		resp, err := client.Accept(ctx, &multipaxos.AcceptRequest{
			Slot:  slot,
			Round: round,
			Id:    myCandidateInfo,
			Value: value,
		})
		conn.Close()
		if err != nil {
			continue
		}
		if resp.GetOk() {
			acceptOKCount++
		}
	}
	return acceptOKCount
}

// doCommitPhase sends Commit to all peers.
func (s *MultiPaxosService) doCommitPhase(ctx context.Context, slot, round, value int32, peers []string) {
	for _, peer := range peers {
		conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		client := multipaxos.NewMultiPaxosServiceClient(conn)
		_, _ = client.Commit(ctx, &multipaxos.CommitRequest{
			Slot:  slot,
			Round: round,
			Id:    myCandidateInfo,
			Value: value,
		})
		conn.Close()
	}
}
