package main

import (
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type Synchronizer struct {
	etcdClient        *clientv3.Client
	multiPaxosService *MultiPaxosService
	multiPaxosClient  *MultiPaxosClient
	state             map[string]*Scooter
	approvedEventLog  []multipaxos.ScooterEvent
	myPendingEvents   *MultipaxosQueue
	nextOrderId       int64
}

func NewSynchronizer(queueSize int, etcdClient *clientv3.Client, state map[string]*Scooter, multiPaxosClient *MultiPaxosClient) *Synchronizer {

	return &Synchronizer{
		etcdClient:       etcdClient,
		multiPaxosClient: multiPaxosClient,
		state:            state,
		approvedEventLog: make([]multipaxos.ScooterEvent, queueSize),
		myPendingEvents:  NewMultipaxosQueue(queueSize),
		nextOrderId:      0,
	}
}

func (s *Synchronizer) CreateScooter(scooterId string) {

	createEvent := &multipaxos.ScooterEvent{
		ScooterId: scooterId,
		OrderId:   s.nextOrderId + 1,
		EventType: &multipaxos.ScooterEvent_CreateEvent{
			CreateEvent: &multipaxos.CreateScooterEvent{},
		},
	}
	s.myPendingEvents.Enqueue(createEvent)

	// run paxos and wait until approved
	// update local state
	// return to customer (rest api)
	s.multiPaxosService.start()
}

func (s *Synchronizer) updateStateWithCommited(event *multipaxos.ScooterEvent) {
	switch x := event.EventType.(type) {
	case *multipaxos.ScooterEvent_CreateEvent:
		log.Printf("submitting create event scooter id = %s", event.ScooterId)
		scooter := &Scooter{
			Id:                   event.ScooterId,
			IsAvailable:          true,
			TotalDistance:        0,
			CurrentReservationId: "",
		}
		s.state[event.ScooterId] = scooter
		// Handle create event
	case *multipaxos.ScooterEvent_ReserveEvent:
		log.Printf("reservation %s", x.ReserveEvent)
		s.state[event.ScooterId].IsAvailable = false
		s.state[event.ScooterId].CurrentReservationId = x.ReserveEvent.ReservationId
		// Handle reserve event
	case *multipaxos.ScooterEvent_ReleaseEvent:
		s.state[event.ScooterId].IsAvailable = true
		s.state[event.ScooterId].TotalDistance += x.ReleaseEvent.Distance
		s.state[event.ScooterId].CurrentReservationId = ""
		// Handle release event
	default:
		log.Fatalf("Unknown scooter event type")
	}
}
