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
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", serverAddr, err)
		return
	}
	defer conn.Close()
	paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

	// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
	prepareReq := &multipaxos.PrepareRequest{Id: c.myId, Round: 1}
	ctx := context.Background()
	_, err = paxosClient.TriggerPrepare(ctx, prepareReq)
	if err != nil {
		log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
	}
}

//func (c *MultiPaxosClient) start() {
//	if !amILeader() {
//		// ping leader to initiate a prepare (if not already initiated)
//		go c.TriggerPrepare()
//		return
//	}
//
//	// Send prepare to everyone
//	ctx := context.TODO()
//	otherServers := fetchOtherServersList(ctx)
//
//	for _, serverAddr := range otherServers {
//		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
//		if err != nil {
//			log.Printf("Failed to connect to server %s: %v", serverAddr, err)
//			continue
//		}
//		defer conn.Close()
//		paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)
//
//		// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
//		prepareReq := &multipaxos.PrepareRequest{Id: myCandidateInfo, Round: 1}
//		_, err = paxosClient.Prepare(ctx, prepareReq)
//		if err != nil {
//			log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
//			continue
//		}
//	}
//}
