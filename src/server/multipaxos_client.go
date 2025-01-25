package main

import (
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type MultiPaxosClient struct {
	myId string
}

func (c *MultiPaxosClient) TriggerPrepare() {
	serverAddr := getLeader()
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", serverAddr, err)
		return
	}
	defer conn.Close()
	paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

	// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
	prepareReq := &multipaxos.PrepareRequest{Id: c.myId, Round: 1}
	ctx := context.Background()
	_, err = paxosClient.Prepare(ctx, prepareReq)
	if err != nil {
		log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
	}
}

func (c *MultiPaxosClient) start() {
	if !amILeader() {
		// ping leader to initiate a prepare (if not already initiated)
		go c.TriggerPrepare()
		return
	}

	// Send prepare to everyone
	ctx := context.TODO()
	otherServers := fetchOtherServersList(ctx)

	for _, serverAddr := range otherServers {
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Failed to connect to server %s: %v", serverAddr, err)
			continue
		}
		defer conn.Close()
		paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

		// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
		prepareReq := &multipaxos.PrepareRequest{Id: myCandidateInfo, Round: 1}
		_, err = paxosClient.Prepare(ctx, prepareReq)
		if err != nil {
			log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
			continue
		}
	}
}

func (c *MultiPaxosClient) SendPromise(leader string, event *multipaxos.ScooterEvent) {
	ctx := context.TODO()

	conn, err := grpc.NewClient(leader, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", leader, err)
		return
	}
	defer conn.Close()
	paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

	promiseReq := &multipaxos.PromiseRequest{Id: myCandidateInfo, Round: 1, Ack: true, Value: event}
	_, err = paxosClient.Promise(ctx, promiseReq)
	if err != nil {
		log.Printf("Failed to promise Paxos on server %s: %v", leader, err)
		return
	}
}
