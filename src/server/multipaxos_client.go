package main

import (
	"context"
	"fmt"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"log"
	"math"
	"sync"
)

type MultiPaxosClient struct {
	myId         string
	peersCluster *PeersCluster
}

func (c *MultiPaxosClient) TriggerPrepare(event *multipaxos.ScooterEvent) {
	log.Printf("Sending Trigger to leader with event : %v", event)
	serverAddr := getLeader()
	conn, err := c.peersCluster.GetLiveConnection(serverAddr)
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", serverAddr, err)
		return
	}
	paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

	// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
	triggerReq := &multipaxos.TriggerRequest{MemberId: c.myId, Event: event}
	ctx := context.Background()
	_, err = paxosClient.TriggerLeader(ctx, triggerReq)
	if err != nil {
		log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
	}
}

// doPreparePhase sends Prepare to all peers and counts how many OK
// Also returns the highest acceptedRound/value from the responding peers.
func (c *MultiPaxosClient) doPreparePhase(ctx context.Context, slot int64, round int32, peers []string) (int32, *multipaxos.ScooterEvent, int) {
	prepareOKCount := 0
	highestAcceptedRound := int32(0)
	highestAcceptedValue := (*multipaxos.ScooterEvent)(nil)

	// Send to each peer
	for _, peer := range peers {
		conn, err := c.peersCluster.GetLiveConnection(peer)
		if err != nil {
			continue
		}
		client := multipaxos.NewMultiPaxosServiceClient(conn)

		resp, err := client.Prepare(ctx, &multipaxos.PrepareRequest{
			Slot:  slot,
			Round: round,
			Id:    myCandidateInfo, // your ID
		})
		if err != nil {
			continue
		}
		if resp.GetOk() {
			prepareOKCount++
			// Track highest acceptedRound/value returned
			value := resp.GetAcceptedValue()
			log.Printf("doPreparePhase accepted value %s, is null = %d", value, value == nil)
			if resp.GetAcceptedRound() > highestAcceptedRound && value != nil {
				highestAcceptedRound = resp.GetAcceptedRound()
				highestAcceptedValue = value
			}
		}
	}

	return highestAcceptedRound, highestAcceptedValue, prepareOKCount
}

// doAcceptPhase sends Accept to all peers with the final chosen value
func (c *MultiPaxosClient) doAcceptPhase(ctx context.Context, slot int64, round int32, value *multipaxos.ScooterEvent, peers []string) int {
	acceptOKCount := 0
	for _, peer := range peers {
		conn, err := c.peersCluster.GetLiveConnection(peer)
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
func (c *MultiPaxosClient) doCommitPhase(ctx context.Context, slot int64, round int32, value *multipaxos.ScooterEvent, peers []string) {
	log.Printf("Sending commit message to everyone with value (scooterId) %s", value.ScooterId)
	for _, peer := range peers {
		conn, err := c.peersCluster.GetLiveConnection(peer)
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
	}
}

func (c *MultiPaxosClient) GetMinLastGoodSlot() int64 {
	ctx := context.Background()
	servers := c.peersCluster.GetPeersList()

	// Channel to collect responses.
	responseChan := make(chan struct {
		*multipaxos.GetLastGoodSlotResponse
		error
	}, len(servers))

	message := &multipaxos.GetLastGoodSlotRequest{MemberId: c.myId}
	// WaitGroup to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Broadcast message to all servers.
	for _, addr := range servers {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			conn, peerError := c.peersCluster.GetLiveConnection(addr)
			if peerError != nil {
				responseChan <- struct {
					*multipaxos.GetLastGoodSlotResponse
					error
				}{nil, peerError}
			}

			// Send the gRPC message.
			paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

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

func (c *MultiPaxosClient) GetLatestSnapshot() (*multipaxos.GetSnapshotResponse, error) {
	ctx := context.Background()
	peers := c.peersCluster.GetPeersList()

	var maxLastGoodSlot int64 = -1
	var latestSnapshot *multipaxos.GetSnapshotResponse

	// Get snapshots from all peers and find the most recent one
	for _, peer := range peers {
		conn, err := c.peersCluster.GetLiveConnection(peer)
		if err != nil {
			log.Printf("Failed to connect to peer %s: %v", peer, err)
			continue
		}

		client := multipaxos.NewMultiPaxosServiceClient(conn)
		snapshot, err := client.GetSnapshot(ctx, &multipaxos.GetSnapshotRequest{
			MemberId: myCandidateInfo,
		})

		if err != nil {
			log.Printf("Failed to get snapshot from peer %s: %v", peer, err)
			continue
		}

		log.Printf("Got snapshot from peer %s, snapshot is: %v", peer, snapshot)
		log.Printf("snapshot.LastGoodSlot = %d, maxLastGoodSlot = %d", snapshot.LastGoodSlot, maxLastGoodSlot)

		if snapshot.LastGoodSlot > maxLastGoodSlot {
			maxLastGoodSlot = snapshot.LastGoodSlot
			latestSnapshot = snapshot
		}
	}

	log.Printf("Recovered snapshot is: %v", latestSnapshot)

	if latestSnapshot == nil {

		return nil, fmt.Errorf("failed to recover snapshot from any peer")
	} else {
		return latestSnapshot, nil
	}
}
