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
	synchronizer     *Synchronizer
	etcdClient       *clientv3.Client
	multiPaxosClient *MultiPaxosClient
}

func startPaxosServer(stopCh chan struct{}, port string, multiPaxosService *MultiPaxosService) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	go func() {
		multipaxos.RegisterMultiPaxosServiceServer(s, multiPaxosService)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	log.Printf("Multipaxos gRPC server listening to port " + port)
	<-stopCh
	log.Println("Shutting down MultiPaxos gRPC server...")
}

func (s *MultiPaxosService) Prepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
	log.Printf("Received Prepare request with ID: %s", req.GetId())
	// find if i have something in queue and if not send null

	return &multipaxos.PrepareResponse{Ok: true}, nil
}
