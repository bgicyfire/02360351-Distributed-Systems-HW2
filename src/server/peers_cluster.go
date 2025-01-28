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
	keyMap      map[string]string // etcd stores servers as '/servers/694d94aec00f4425', this is the key received on delete, need to convert to peer (ip:port)
	//peers       []string          // a cached version to save time of converting map keys to array
	mu sync.Mutex // To protect concurrent access to the connections map.
}

func NewPeersCluster() *PeersCluster {
	return &PeersCluster{
		connections: make(map[string]*grpc.ClientConn),
		keyMap:      make(map[string]string),
		//peers:       make([]string, 0),
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
	//return pc.peers
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
			log.Printf("Failed to connect to server %s: %v", peer, err)
			return nil, err
		}
		pc.connections[peer] = conn
	}
	return pc.connections[peer], nil
}

func (pc *PeersCluster) AddPeer(key, addr string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.keyMap[key] = addr
	if _, exists := pc.connections[addr]; exists {
		return
	}
	pc.connections[addr] = nil
	//pc.peers = append(pc.peers, addr)
}

// RemovePeer accepts key (format '/servers/694d94aec00f4425') not peer (format '172.19.0.5:50052')
func (pc *PeersCluster) RemovePeer(key string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if addr, exists := pc.keyMap[key]; exists {
		delete(pc.connections, addr)
		delete(pc.keyMap, key)
		//pc.peers = make([]string, 0, len(pc.connections))
		//for peer := range pc.connections {
		//	pc.peers = append(pc.peers, peer)
		//}
	}
}
