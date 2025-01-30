package main

import (
	"context"
	"fmt"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"sync"
	"time"
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

func (c *MultiPaxosClient) TryGetPaxosInstanceFromPeers(slot int64) (*PaxosInstance, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the list of peers from etcd
	peers := c.peersCluster.GetPeersList()
	log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Attempting to get slot %d from %d peers", slot, len(peers))
	err1 := error(nil)
	for _, peer := range peers {
		// Skip own address to avoid self-connection
		if peer == myCandidateInfo {
			continue
		}

		// Create connection with non-blocking dial
		conn, err := grpc.Dial(peer,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(2*time.Second),
			grpc.WithBlock(),                 // Wait for connection with timeout
			grpc.WithReturnConnectionError()) // Return concrete error

		if err != nil {
			err1 = err
			log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Failed to connect to peer %s: %v", peer, err)
			continue
		}

		// Ensure connection is closed
		defer conn.Close()

		// Check connection state before proceeding
		state := conn.GetState()
		if state != connectivity.Ready {
			log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Connection to peer %s not ready, state: %s", peer, state)
			continue
		}
		client := multipaxos.NewMultiPaxosServiceClient(conn)

		// Use a shorter timeout for the RPC call itself
		callCtx, callCancel := context.WithTimeout(ctx, 1*time.Second)
		defer callCancel()

		resp, err := client.GetPaxosInstance(callCtx, &multipaxos.GetPaxosInstanceRequest{
			Slot: slot,
		})

		if err != nil {
			err1 = err
			if ctx.Err() == context.DeadlineExceeded {
				log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Timeout getting state from peer %s for slot %d", peer, slot)
			} else {
				log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Error getting state from peer %s for slot %d: %v", peer, slot, err)
			}
			continue
		}
		if !resp.InstanceFound {
			log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Peer %s responded and does not have paxos instance of slot %d", peer, slot)
			continue
		}

		instance := &PaxosInstance{
			promisedRound: resp.GetPromisedRound(),
			acceptedRound: resp.GetAcceptedRound(),
			acceptedValue: resp.GetAcceptedValue(),
			committed:     resp.GetCommitted(),
		}

		log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] Successfully received instance for slot %d from peer %s", slot, peer)
		return instance, nil
	}

	log.Printf("[Recovery:TryGetPaxosInstanceFromPeers] No peers had instance for slot %d", slot)
	return nil, err1
}
