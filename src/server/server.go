package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	leaderInfo      string
	leaderMutex     sync.RWMutex
	myCandidateInfo string
	etcdClient      *clientv3.Client
	scooters        map[string]*Scooter
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

func fetchAllServersList(ctx context.Context) []string {
	// Fetch server list from etcd
	log.Printf("getting servers list form etcd")
	resp, err := etcdClient.Get(ctx, "/servers/", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Failed to get servers from etcd: %v", err)
	}
	log.Printf("received servers list from etcd: %v", resp.Kvs)

	uniqueServers := make(map[string]bool)
	var result []string
	for _, ev := range resp.Kvs {
		serverAddr := string(ev.Value)
		//if serverAddr == myCandidateInfo {
		//	continue // Skip own address
		//}
		if _, exists := uniqueServers[serverAddr]; !exists {
			uniqueServers[serverAddr] = true
			result = append(result, serverAddr)
			//log.Printf("Fetch found server %s", serverAddr)
		}
	}
	//log.Printf("Fetch will return %d servers", len(result))
	return result
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
	paxosPort := os.Getenv("PAXOS_PORT")
	myCandidateInfo = getLocalIP() + ":" + paxosPort // Unique server identification

	log.Printf("My candidate info : %s", myCandidateInfo)

	queueSize, err := strconv.ParseInt(os.Getenv("LOCAL_QUEUE_SIZE"), 10, 64)
	if err != nil {
		log.Fatalf("Invalid local queue size env var: %v", err)
	}
	scooters = make(map[string]*Scooter)
	multiPaxosClient := &MultiPaxosClient{myId: myCandidateInfo}
	synchronizer := NewSynchronizer(int(queueSize), etcdClient, scooters, multiPaxosClient)
	multiPaxosService := NewMultiPaxosService(synchronizer, etcdClient, multiPaxosClient)
	//multiPaxosService := &MultiPaxosService{synchronizer: synchronizer, etcdClient: etcdClient, multiPaxosClient: multiPaxosClient}
	synchronizer.multiPaxosService = multiPaxosService

	log.Printf("Starting server")

	// Channel to signal goroutines to stop
	stopCh := make(chan struct{})
	// Channel to catch system signals for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	initEtcdClient()
	defer etcdClient.Close()

	go runLeaderElection(etcdClient, myCandidateInfo)
	go observeLeader(etcdClient, "/servers")
	// Start the registerHost function in a separate goroutine
	//go registerHost(stopCh)
	log.Printf("Registered to etcd")

	go startPaxosServer(stopCh, paxosPort, multiPaxosService)

	go startScooterService(stopCh, etcdClient, scooters, synchronizer)

	// Waiting for shutdown signal
	<-signalCh
	close(stopCh) // signal all goroutines to stop
	log.Println("Server is shutting down")
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
