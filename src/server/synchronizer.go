package main

import (
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Synchronizer struct {
	etcdClient        *clientv3.Client
	multiPaxosService *MultiPaxosService
	state             map[string]*Scooter
	approvedEventLog  []multipaxos.ScooterEvent
	myPendingEvents   *MultipaxosQueue
	nextOrderId       int64
}

func NewSynchronizer(queueSize int,
	etcdClient *clientv3.Client,
	state map[string]*Scooter) *Synchronizer {

	return &Synchronizer{
		etcdClient:       etcdClient,
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
