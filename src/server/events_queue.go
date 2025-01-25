package main

import (
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
)

type MultipaxosQueue struct {
	channel chan *multipaxos.ScooterEvent
}

func NewMultipaxosQueue(size int) *MultipaxosQueue {
	return &MultipaxosQueue{channel: make(chan *multipaxos.ScooterEvent, size)}
}

func (q *MultipaxosQueue) Enqueue(item *multipaxos.ScooterEvent) {
	q.channel <- item
}

func (q *MultipaxosQueue) Dequeue() *multipaxos.ScooterEvent {
	return <-q.channel
}
