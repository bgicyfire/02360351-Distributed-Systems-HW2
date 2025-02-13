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

const ETCD_SERVERS_PREFIX = "/servers"

var (
	leaderInfo      string
	leaderMutex     sync.RWMutex
	myCandidateInfo string
	etcdClient      *clientv3.Client
	scooters        map[string]*Scooter
)

//func fetchAllServersList(ctx context.Context) []string {
//	// Fetch server list from etcd
//	log.Printf("getting servers list form etcd")
//	resp, err := etcdClient.Get(ctx, "/servers/", clientv3.WithPrefix())
//	if err != nil {
//		log.Fatalf("Failed to get servers from etcd: %v", err)
//	}
//	log.Printf("received servers list from etcd: %v", resp.Kvs)
//
//	uniqueServers := make(map[string]bool)
//	var result []string
//	for _, ev := range resp.Kvs {
//		serverAddr := string(ev.Value)
//		//if serverAddr == myCandidateInfo {
//		//	continue // Skip own address
//		//}
//		if _, exists := uniqueServers[serverAddr]; !exists {
//			uniqueServers[serverAddr] = true
//			result = append(result, serverAddr)
//			//log.Printf("Fetch found server %s", serverAddr)
//		}
//	}
//	//log.Printf("Fetch will return %d servers", len(result))
//	return result
//}

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

	snapshotInterval, err := strconv.ParseInt(os.Getenv("SNAPSHOT_INTERVAL"), 10, 64)
	if err != nil {
		log.Fatalf("Invalid SNAPSHOT_INTERVAL env var: %v", err)
	}

	myCandidateInfo = getLocalIP() + ":" + paxosPort // Unique server identification
	log.Printf("My candidate info : %s", myCandidateInfo)

	scooters = make(map[string]*Scooter)
	initEtcdClient()
	defer etcdClient.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	peersCluster := LoadInitialServersList(etcdClient, ETCD_SERVERS_PREFIX, ctx)
	go watchServers(etcdClient, ETCD_SERVERS_PREFIX, peersCluster, ctx)
	multiPaxosClient := &MultiPaxosClient{myId: myCandidateInfo, peersCluster: peersCluster}
	synchronizer := NewSynchronizer(snapshotInterval, scooters, multiPaxosClient)
	recoveryDoneCh := make(chan struct{})
	multiPaxosService := NewMultiPaxosService(synchronizer, multiPaxosClient, peersCluster, recoveryDoneCh)
	synchronizer.multiPaxosService = multiPaxosService

	log.Printf("Starting server")

	// Channel to signal goroutines to stop
	stopCh := make(chan struct{})
	// Channel to catch system signals for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go runLeaderElection(etcdClient, ETCD_SERVERS_PREFIX, myCandidateInfo, ctx)
	go observeLeader(etcdClient, ETCD_SERVERS_PREFIX, ctx)
	// Start the registerHost function in a separate goroutine
	//go registerHost(stopCh)
	log.Printf("Registered to etcd")

	// start the paxos gRPC server, but dont allow starting new paxos instances (only if is leader) until recovery is completed
	go startPaxosServer(stopCh, paxosPort, multiPaxosService)
	multiPaxosService.recoverAllInstances()
	close(recoveryDoneCh)

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
