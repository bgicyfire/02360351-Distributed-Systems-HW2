package main

import (
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"github.com/google/uuid"
	"sync"
)

// Node represents a node in the linked list
type Node struct {
	Key   string
	Value *multipaxos.ScooterEvent
	prev  *Node
	next  *Node
}

// MultipaxosQueue implements a thread-safe FIFO queue with key-based removal
type MultipaxosQueue struct {
	head    *Node
	tail    *Node
	nodeMap map[string]*Node
	mutex   sync.RWMutex
	length  int
}

// NewMultipaxosQueue creates a new MultipaxosQueue instance
func NewMultipaxosQueue() *MultipaxosQueue {
	return &MultipaxosQueue{
		nodeMap: make(map[string]*Node),
	}
}
func (q *MultipaxosQueue) EnqueueGenerateKey(value *multipaxos.ScooterEvent) string {
	key := uuid.New().String()
	value.EventId = key
	q.Enqueue(key, value)
	return key
}

// Enqueue adds a new item to the queue
// If an item with the same key exists, it will be replaced
func (q *MultipaxosQueue) Enqueue(key string, value *multipaxos.ScooterEvent) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// If key exists, remove the old node first
	if existing, exists := q.nodeMap[key]; exists {
		q.removeNode(existing)
	}

	// Create new node
	newNode := &Node{
		Key:   key,
		Value: value,
	}

	// Add to end of list
	if q.tail == nil {
		q.head = newNode
		q.tail = newNode
	} else {
		newNode.prev = q.tail
		q.tail.next = newNode
		q.tail = newNode
	}

	// Add to map
	q.nodeMap[key] = newNode
	q.length++
}

// Peek returns the oldest item without removing it
// Returns nil if queue is empty
func (q *MultipaxosQueue) Peek() *multipaxos.ScooterEvent {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.head == nil {
		return nil
	}
	return q.head.Value
}

// Remove removes an item by key
// Returns true if item was found and removed, false otherwise
func (q *MultipaxosQueue) Remove(key string) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if node, exists := q.nodeMap[key]; exists {
		q.removeNode(node)
		return true
	}
	return false
}

// removeNode removes a node from the linked list
// Caller must hold the lock
func (q *MultipaxosQueue) removeNode(node *Node) {
	// Update adjacent nodes
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		q.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		q.tail = node.prev
	}

	// Remove from map
	delete(q.nodeMap, node.Key)
	q.length--
}

// Length returns the current number of items in the queue
func (q *MultipaxosQueue) Length() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return q.length
}
