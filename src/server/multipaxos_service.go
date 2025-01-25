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

func startPaxosServer(stopCh chan struct{}, port string, etcdClient *clientv3.Client, synchronizer *Synchronizer) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	multipaxos.RegisterMultiPaxosServiceServer(s, &MultiPaxosService{synchronizer: synchronizer, etcdClient: etcdClient})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *MultiPaxosService) Prepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
	log.Printf("Received Prepare request with ID: %s", req.GetId())
	return &multipaxos.PrepareResponse{Ok: true}, nil
}
