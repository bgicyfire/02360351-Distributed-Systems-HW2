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

func runLeaderElection(client *clientv3.Client, candidateInfo string) {
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
	election := concurrency.NewElection(sess, "/servers")

	ctx := context.TODO()

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

func observeLeader(client *clientv3.Client, prefix string) {
	sess, err := concurrency.NewSession(client)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer sess.Close()

	election := concurrency.NewElection(sess, prefix)

	ctx := context.Background()

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
