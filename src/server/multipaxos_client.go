package main

import (
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"sync"
)

type MultiPaxosClient struct {
	myId string
}

func (c *MultiPaxosClient) TriggerPrepare(event *multipaxos.ScooterEvent) {
	log.Printf("Sending Trigger to leader with event : %v", event)
	serverAddr := getLeader()
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", serverAddr, err)
		return
	}
	defer conn.Close()
	paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

	// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
	triggerReq := &multipaxos.TriggerRequest{MemberId: c.myId, Event: event}
	ctx := context.Background()
	_, err = paxosClient.TriggerLeader(ctx, triggerReq)
	if err != nil {
		log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
	}
}

func (c *MultiPaxosClient) GetMinLastGoodSlot() int64 {
	ctx := context.Background()
	servers := fetchAllServersList(ctx)
	// Map to store gRPC connections.
	connections := make(map[string]*grpc.ClientConn)
	var mu sync.Mutex // To protect concurrent access to the connections map.

	// Channel to collect responses.
	responseChan := make(chan struct {
		*multipaxos.GetLastGoodSlotResponse
		error
	}, len(servers))

	// WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Broadcast message to all servers.
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	for _, addr := range servers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()

			// Check if the connection exists and is healthy.
			mu.Lock()
			conn, exists := connections[addr]
			if !exists || conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
				// Reestablish the connection if it's dead or doesn't exist.
				var err error
				conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					//responseChan <- ServerResponse{Address: addr, Err: err}
					log.Printf("Failed to connect to server %s: %v", addr, err)
					mu.Unlock()
					return
				}
				connections[addr] = conn
			}
			mu.Unlock()

			// Send the gRPC message.
			paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)
			message := &multipaxos.GetLastGoodSlotRequest{MemberId: c.myId}

			response, err := paxosClient.GetLastGoodSlot(ctx, message)

			// Send the response to the channel.
			responseChan <- struct {
				*multipaxos.GetLastGoodSlotResponse
				error
			}{response, err}
		}(addr)
	}

	// Wait for all goroutines to finish in a separate goroutine.
	go func() {
		wg.Wait()
		close(responseChan)
	}()

	// Collect responses and wait until at least half of the servers respond.
	successCount := 0
	failureCount := 0
	requiredResponses := len(servers) / 2

	minSlot := int64(math.MaxInt64)
	for resp := range responseChan {
		if resp.error != nil {
			log.Printf("GetMinLastGoodSlot Error from server %v\n", resp.error)
			failureCount++
		} else {
			minSlot = min(minSlot, resp.LastGoodSlot)
			successCount++
		}

		// Stop waiting once we have enough responses.
		if successCount >= requiredResponses {
			break
		}
	}

	log.Printf("Broadcast completed. Success: %d, Failures: %d\n", successCount, failureCount)
	if successCount < requiredResponses {
		log.Fatalf("GetMinLastGoodSlot: Received less than quorum responses, cannot continue")
	}
	return minSlot
}
