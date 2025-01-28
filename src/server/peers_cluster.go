package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

type PeersCluster struct {
	connections map[string]*grpc.ClientConn
	mu          sync.Mutex // To protect concurrent access to the connections map.
}

func NewPeersCluster() *PeersCluster {
	return &PeersCluster{
		connections: make(map[string]*grpc.ClientConn),
	}
}

func (pc *PeersCluster) GetPeersList() []string {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	peers := make([]string, 0, len(pc.connections))
	for peer := range pc.connections {
		peers = append(peers, peer)
	}
	return peers
}

func (pc *PeersCluster) GetLiveConnection(peer string) (*grpc.ClientConn, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	// check if server exists
	conn, exists := pc.connections[peer]
	if !exists {
		return nil, fmt.Errorf("peer does not exist %s", peer)
	}
	// Check if the connection is healthy.
	if conn == nil || conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		// Reestablish the connection if it's dead or doesn't exist.
		var err error
		conn, err = grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			//responseChan <- ServerResponse{Address: addr, Err: err}
			log.Printf("Failed to connect to server %s: %v", peer, err)
			return nil, err
		}
		pc.connections[peer] = conn
	}
	return pc.connections[peer], nil
}

func (pc *PeersCluster) AddPeer(addr string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if _, exists := pc.connections[addr]; exists {
		return
	}
	pc.connections[addr] = nil
}

func (pc *PeersCluster) RemovePeer(addr string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if _, exists := pc.connections[addr]; exists {
		delete(pc.connections, addr)
	}

}
