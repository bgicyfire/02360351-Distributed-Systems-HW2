package main

import (
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"log"
	"sync"
)

type Synchronizer struct {
	mu                sync.RWMutex
	multiPaxosService *MultiPaxosService
	multiPaxosClient  *MultiPaxosClient
	//state             map[string]*Scooter
	//lastGoodSlot     int64
	state             *Snapshot
	approvedEventLog  map[int64]*multipaxos.ScooterEvent
	myPendingEvents   *MultipaxosQueue
	snapshot          *Snapshot
	snapshot_interval int64
}

type Snapshot struct {
	state        map[string]*Scooter
	lastGoodSlot int64
}

// TODO : remove queueSize if not needed
func NewSynchronizer(snapshot_interval int64, state map[string]*Scooter, multiPaxosClient *MultiPaxosClient) *Synchronizer {

	return &Synchronizer{
		mu:               sync.RWMutex{},
		multiPaxosClient: multiPaxosClient,
		//state:            state,
		//lastGoodSlot:     0,
		state:             &Snapshot{state: state, lastGoodSlot: 0},
		snapshot:          &Snapshot{state: make(map[string]*Scooter), lastGoodSlot: 0},
		approvedEventLog:  make(map[int64]*multipaxos.ScooterEvent),
		myPendingEvents:   NewMultipaxosQueue(),
		snapshot_interval: snapshot_interval,
	}
}

func (s *Synchronizer) CreateScooter(scooterId string) {

	createEvent := &multipaxos.ScooterEvent{
		ScooterId: scooterId,
		EventType: &multipaxos.ScooterEvent_CreateEvent{
			CreateEvent: &multipaxos.CreateScooterEvent{},
		},
	}
	s.myPendingEvents.EnqueueGenerateKey(createEvent)

	// run paxos and wait until approved
	// update local state
	// return to customer (rest api)
	s.startPaxos()
}

func (s *Synchronizer) ReserveScooter(scooterId string, reservationId string) {

	reservation := &multipaxos.ScooterEvent{
		ScooterId: scooterId,
		EventType: &multipaxos.ScooterEvent_ReserveEvent{
			ReserveEvent: &multipaxos.ReserveScooterEvent{
				ReservationId: reservationId,
			},
		},
	}
	s.myPendingEvents.EnqueueGenerateKey(reservation)

	// run paxos and wait until approved
	// update local state
	// return to customer (rest api)
	s.startPaxos()
}

func (s *Synchronizer) ReleaseScooter(scooterId string, reservationId string, rideDistance int64) {

	release := &multipaxos.ScooterEvent{
		ScooterId: scooterId,
		EventType: &multipaxos.ScooterEvent_ReleaseEvent{
			ReleaseEvent: &multipaxos.ReleaseScooterEvent{
				ReservationId: reservationId,
				Distance:      rideDistance,
			},
		},
	}
	s.myPendingEvents.EnqueueGenerateKey(release)

	// run paxos and wait until approved
	// update local state
	// return to customer (rest api)
	s.startPaxos()
}

func (s *Synchronizer) updateStateWithCommited(slot int64, event *multipaxos.ScooterEvent) {
	// TODO : add locks here
	if existingEvent, exists := s.approvedEventLog[slot]; exists {
		// an event is already registered with this slot
		if existingEvent.EventId == event.EventId {
			// it's the same event, already considered in state, we can ignore
			return
		} else {
			// this is a different event, edge case
			// TODO: what to do
			log.Printf("[Synchronizer:updateStateWithCommited] unexpected existing slot with different event id")
		}
	}

	s.approvedEventLog[slot] = event
	s.myPendingEvents.Remove(event.EventId)
	s.updateState(s.state, -1)

	// do snapshot if needed
	if s.state.lastGoodSlot%s.snapshot_interval == 0 {
		log.Printf("Snapshot: PRE Approved events log: %v", s.approvedEventLog)
		minSlot := s.multiPaxosClient.GetMinLastGoodSlot()
		log.Printf("Snapshot: s.state.lastGoodSlot: %d , minSlot: %d", s.state.lastGoodSlot, minSlot)
		minSlot = min(minSlot, s.state.lastGoodSlot)
		previousSnapshotSlot := s.snapshot.lastGoodSlot
		log.Printf("Snapshot: pre previousSnapshotSlot : %d, minSlot: %d", previousSnapshotSlot, minSlot)
		s.updateState(s.snapshot, minSlot)
		for i := previousSnapshotSlot + 1; i <= minSlot; i++ {
			delete(s.approvedEventLog, i)
		}
		log.Printf("Snapshot: post previousSnapshotSlot : %d, %d", previousSnapshotSlot, minSlot)
		log.Printf("Snapshot: completed, compacted %d events", minSlot-previousSnapshotSlot)
		log.Printf("Snapshot: Approved events log: %v", s.approvedEventLog)
	}
}

func (s *Synchronizer) updateState(snapshot *Snapshot, maxSlot int64) {
	// TODO : add locks
	for slot := snapshot.lastGoodSlot + 1; ; slot++ {
		// if we are making snapshot (not the actual state), update until slot == maxSlot (that is the minimum that we received from quorum)
		if maxSlot > 0 && slot > maxSlot {
			return
		}
		event, exists := s.approvedEventLog[slot]
		if !exists {
			return
		}
		snapshot.lastGoodSlot++
		switch x := event.EventType.(type) {
		case *multipaxos.ScooterEvent_CreateEvent:
			log.Printf("submitting create event scooter id = %s", event.ScooterId)
			scooter := &Scooter{
				Id:                   event.ScooterId,
				IsAvailable:          true,
				TotalDistance:        0,
				CurrentReservationId: "",
			}
			snapshot.state[event.ScooterId] = scooter
			// Handle create event
		case *multipaxos.ScooterEvent_ReserveEvent:
			log.Printf("reservation %s", x.ReserveEvent)
			if scooter, ok := snapshot.state[event.ScooterId]; ok {
				if scooter.IsAvailable {
					scooter.IsAvailable = false
					scooter.CurrentReservationId = x.ReserveEvent.ReservationId
				}
			} else {
				log.Printf("!!! scooter %s not found", event.ScooterId)
			}
			// Handle reserve event
		case *multipaxos.ScooterEvent_ReleaseEvent:
			if scooter, ok := snapshot.state[event.ScooterId]; ok {
				if !scooter.IsAvailable {
					scooter.IsAvailable = true
					scooter.TotalDistance += x.ReleaseEvent.Distance
					scooter.CurrentReservationId = ""
				}
			} else {
				log.Printf("!!! scooter %s not found", event.ScooterId)
			}
			// Handle release event
		default:
			log.Fatalf("Unknown scooter event type")
		}
	}
}

func (s *Synchronizer) startPaxos() {
	s.multiPaxosService.start()
}
