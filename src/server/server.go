package main

import (
	"bufio"
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/scooter"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type server struct {
	scooter.UnimplementedScooterServiceServer
}

func getContainerID() string {
	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "/")
		if len(parts) > 2 {
			idPart := parts[len(parts)-1]
			if len(idPart) == 64 { // Typical length of Docker container IDs
				return idPart
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading cgroup info: %v", err)
	}
	return "unknown"
}

// GetScooterStatus implements scooter.ScooterServiceServer
func (s *server) GetScooterStatus(ctx context.Context, in *scooter.ScooterRequest) (*scooter.ScooterResponse, error) {
	containerID := getContainerID()

	log.Printf("Received: %v, handled by: %s", in.GetScooterId(), containerID)
	return &scooter.ScooterResponse{
		Status:   "Available",
		Hostname: containerID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	scooter.RegisterScooterServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
