package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"strconv"
	"time"
)

func getLeader() string {
	leaderMutex.RLock()
	defer leaderMutex.RUnlock()
	return leaderInfo
}

func amILeader() bool {
	return getLeader() == myCandidateInfo
}

func setLeader(info string) {
	leaderMutex.Lock()
	defer leaderMutex.Unlock()
	leaderInfo = info
}

func runLeaderElection(client *clientv3.Client, prefix string, candidateInfo string, ctx context.Context) {
	leaseDurationStr := os.Getenv("ETCD_LEASE_DURATION")
	leaseDuration, err := strconv.ParseInt(leaseDurationStr, 10, 32)
	if err != nil {
		log.Fatalf("Invalid lease duration: %v", err)
	}
	// Create a sessions to keep the lease alive
	sess, err := concurrency.NewSession(client, concurrency.WithTTL(int(leaseDuration)))
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer sess.Close()

	// Create an election instance on the given key prefix
	election := concurrency.NewElection(sess, prefix)

	// Campaign to become the leader
	err = election.Campaign(ctx, candidateInfo)
	if err != nil {
		log.Fatalf("Failed to campaign for leadership: %v", err)
	}

	// Announce leadership if won
	resp, err := election.Leader(ctx)
	if err != nil {
		log.Fatalf("Failed to get leader: %v", err)
	}
	// Set the leader information locally
	setLeader(string(resp.Kvs[0].Value))
	log.Printf("The leader is %s", string(resp.Kvs[0].Value))
	if string(resp.Kvs[0].Value) == candidateInfo {
		log.Printf("I am the leader :)")
	} else {
		log.Printf("I am NOT the leader")
	}

	// Keep checking the leader periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down gracefully...")
			// Resign from leadership before exiting
			if err := election.Resign(context.Background()); err != nil {
				log.Printf("Failed to resign from leadership: %v", err)
			}
			return
		case <-ticker.C:
			// Refresh leader information
			resp, err := election.Leader(ctx)
			if err != nil {
				log.Printf("Error fetching leader: %v", err)
				continue
			}
			currentLeader := string(resp.Kvs[0].Value)
			if getLeader() != currentLeader {
				setLeader(currentLeader)
			}
		case <-sess.Done():
			log.Println("Session expired or canceled, re-campaigning for leadership")
			return
		}
	}
}

func observeLeader(client *clientv3.Client, prefix string, ctx context.Context) {
	sess, err := concurrency.NewSession(client)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer sess.Close()

	election := concurrency.NewElection(sess, prefix)

	for {
		resp, err := election.Leader(ctx)
		if err != nil {
			log.Printf("Error fetching leader: %v", err)
			time.Sleep(5 * time.Second) // wait before retrying
			continue
		}

		currentLeader := string(resp.Kvs[0].Value)
		if getLeader() != currentLeader {
			setLeader(currentLeader)
			log.Printf("New leader observed: %s", currentLeader)
		}

		// Poll every 5 seconds to check for changes in leadership
		time.Sleep(5 * time.Second)
	}
}

// watchServers watches the etcd prefix for changes and updates the local servers map.
func watchServers(cli *clientv3.Client, prefix string, cluster *PeersCluster, ctx context.Context) {
	// Get the initial list of servers under the prefix.
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Failed to get initial list of servers: %v", err)
	}

	for _, kv := range resp.Kvs {
		cluster.AddPeer(string(kv.Value))
	}

	// Start watching for changes.
	watchChan := cli.Watch(ctx, prefix, clientv3.WithPrefix())

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			key := string(event.Kv.Key)
			value := string(event.Kv.Value)

			switch event.Type {
			case clientv3.EventTypePut:
				// New key or updated key.
				cluster.AddPeer(value)
				log.Printf("Server added/updated: %s -> %s\n", key, value)
			case clientv3.EventTypeDelete:
				// Key deleted.
				cluster.RemovePeer(value)
				log.Printf("Server deleted: %s\n", key)
			}
		}
	}
}
