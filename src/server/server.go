package main

import (
	"bufio"
	"context"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/multipaxos"
	"github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/github.com/bgicyfire/02360351-Distributed-Systems-HW2/src/server/scooter"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	leaderInfo  string
	leaderMutex sync.RWMutex
)

func getLeader() string {
	leaderMutex.RLock()
	defer leaderMutex.RUnlock()
	return leaderInfo
}

func getElectedLeader(ctx context.Context) string {
	resp, err := election.Leader(ctx)
	if err != nil {
		// This usually means no leader is set yet or an etcd error
		log.Printf("Error fetching leader: %v", err)
		return ""
	}
	return string(resp.Kvs[0].Value)
}
func setLeader(info string) {
	leaderMutex.Lock()
	defer leaderMutex.Unlock()
	leaderInfo = info
}

var etcdClient *clientv3.Client
var scooters map[string]*Scooter

// election is the global pointer to the concurrency.Election we created
var election *concurrency.Election

type server struct {
	scooter.UnimplementedScooterServiceServer
	multipaxos.UnimplementedMultiPaxosServiceServer
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

func (s *server) Prepare(ctx context.Context, req *multipaxos.PrepareRequest) (*multipaxos.PrepareResponse, error) {
	log.Printf("Received Prepare request with ID: %s", req.GetId())
	return &multipaxos.PrepareResponse{Ok: true}, nil
}

// GetScooterStatus implements scooter.ScooterServiceServer
func (s *server) GetScooterStatus(ctx context.Context, in *scooter.ScooterRequest) (*scooter.ScooterResponse, error) {
	containerID := getContainerID()

	log.Printf("Received: %v, handled by: %s", in.GetScooterId(), containerID)

	// Fetch server list from etcd
	log.Printf("getting servers list form etcd")
	resp, err := etcdClient.Get(ctx, "/servers/", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Failed to get servers from etcd: %v", err)
	}
	log.Printf("received servers list from etcd: %v", resp.Kvs)

	for _, ev := range resp.Kvs {
		serverAddr := string(ev.Value)
		if serverAddr == getLocalIP()+":50052" {
			continue // Skip own address
		}
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Failed to connect to server %s: %v", serverAddr, err)
			continue
		}
		defer conn.Close()
		paxosClient := multipaxos.NewMultiPaxosServiceClient(conn)

		// Assuming Prepare takes an ID and returns a *multipaxos.PrepareResponse
		prepareReq := &multipaxos.PrepareRequest{Id: containerID}
		_, err = paxosClient.Prepare(ctx, prepareReq)
		if err != nil {
			log.Printf("Failed to prepare Paxos on server %s: %v", serverAddr, err)
			continue
		}
	}
	return &scooter.ScooterResponse{
		Status:   "Available",
		Hostname: containerID,
		Myleader: getLeader(), // getElectedLeader(ctx),
	}, nil
}

// Initialize etcd client
func initEtcdClient() {
	var err error
	etcdServer := os.Getenv("ETCD_SERVER")
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdServer},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to initialize etcd client: %v", err)
	}
}

func main() {
	log.Printf("Starting server")
	scooters = make(map[string]*Scooter)

	// Channel to signal goroutines to stop
	stopCh := make(chan struct{})
	// Channel to catch system signals for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	initEtcdClient()
	defer etcdClient.Close()

	candidateInfo := getLocalIP() + ":50052" // Unique server identification
	go runLeaderElection(etcdClient, candidateInfo)

	// Start the registerHost function in a separate goroutine
	//go registerHost(stopCh)
	log.Printf("Registered to etcd")

	go startPaxosServer(stopCh)

	log.Printf("Multipaxos server listening to port 50052")

	go startScooterService(stopCh, etcdClient, scooters)
	go startScooterServer(stopCh) //grpc
	log.Printf("Scooter server listening to port 50051")

	// Waiting for shutdown signal
	<-signalCh
	close(stopCh) // signal all goroutines to stop
	log.Println("Server is shutting down")
}

func startScooterServer(stopCh chan struct{}) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	scooter.RegisterScooterServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func startPaxosServer(stopCh chan struct{}) {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	multipaxos.RegisterMultiPaxosServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// registerHost manages etcd registration and lease renewal.
func registerHost(stopCh <-chan struct{}) {
	etcdServer := os.Getenv("ETCD_SERVER")
	leaseDurationStr := os.Getenv("ETCD_LEASE_DURATION")
	localIP := getLocalIP()
	serverAddress := localIP + ":50052"
	// Convert lease duration from string to int64
	leaseDuration, err := strconv.ParseInt(leaseDurationStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid lease duration: %v", err)
	}

	// Establish a new etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdServer},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer client.Close()

	// Calculated interval for trying keep-alive renewals
	keepAliveInterval := time.Duration(leaseDuration/2) * time.Second

	for {
		select {
		case <-stopCh:
			log.Println("Stopping etcd registration.")
			return
		default:
			// Grant a new lease
			lease, err := client.Grant(context.Background(), leaseDuration)
			if err != nil {
				log.Printf("Failed to create etcd lease: %v", err)
				time.Sleep(1 * time.Second) // Simple backoff before retrying
				continue
			}

			containerID := getContainerID()
			key := "/servers/" + containerID

			// Put a key with the lease
			_, err = client.Put(context.Background(), key, serverAddress, clientv3.WithLease(lease.ID))
			if err != nil {
				log.Printf("Failed to set etcd key: %v", err)
				time.Sleep(1 * time.Second) // Simple backoff before retrying
				continue
			}

			// Start keep-alive for the lease
			keepAliveChan, err := client.KeepAlive(context.Background(), lease.ID)
			if err != nil {
				log.Printf("Failed to keep etcd lease alive: %v", err)
				continue
			}

			// Handle the keep-alive responses
			for {
				select {
				case ka, ok := <-keepAliveChan:
					if !ok {
						log.Println("KeepAlive channel closed. Re-establishing lease...")
						break // Exit this inner loop to re-grant the lease
					}
					log.Printf("Lease keep-alive for key %s at revision %d", key, ka.Revision)
					time.Sleep(keepAliveInterval) // Wait before the next renewal attempt
				case <-stopCh:
					log.Println("Stopping etcd registration.")
					return
				}
			}
		}
	}
}

// getLocalIP attempts to determine the local IP address of the host running the container
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}
	for _, address := range addrs {
		// Check the address type and if it is not a loopback type, return it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

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
	election = concurrency.NewElection(sess, "/servers")

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
