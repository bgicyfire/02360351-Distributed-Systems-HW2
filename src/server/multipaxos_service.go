package main

import (
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
)

type MultiPaxosService struct {
	multipaxos.UnimplementedMultiPaxosServiceServer
	synchronizer *Synchronizer
	etcdClient   *clientv3.Client
}

func startPaxosServer(stopCh chan struct{}, port string, multiPaxosService *MultiPaxosService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	multipaxos.RegisterMultiPaxosServiceServer(s, multiPaxosService)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *MultiPaxosService) Prepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
	log.Printf("Received Prepare request with ID: %s", req.GetId())

	return &multipaxos.PrepareResponse{Ok: true}, nil
}

func (s *MultiPaxosService) start() {
	if !amILeader() {
		// ping leader to initiate a prepare (if not already initiated)
		return
	}

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
