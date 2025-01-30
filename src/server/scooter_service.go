package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ScooterService struct {
	etcdClient   *clientv3.Client
	scooters     map[string]*Scooter
	synchronizer *Synchronizer
}

type ReleaseScooterRequest struct {
	ReservationId string `json:"reservation_id"`
	RideDistance  int64  `json:"ride_distance"` // Assuming distance is an integer value; adjust type if needed
}

type ServerInfo struct {
	MemberId string `json:"member_id"`
	LeaderId string `json:"leader_id"`
	IsLeader bool   `json:"is_leader"`
}

func getServerInfo() ServerInfo {
	return ServerInfo{
		MemberId: myCandidateInfo,
		LeaderId: getLeader(),
		IsLeader: amILeader(),
	}
}

func (s *ScooterService) getScooters(c *gin.Context) {
	scootersList := make([]*Scooter, 0, len(s.scooters))
	for _, scooter := range s.scooters {
		scootersList = append(scootersList, scooter)
	}
	c.JSON(http.StatusOK, gin.H{"scooters": scootersList, "responder": getServerInfo()})
}

func (s *ScooterService) updateScooter(c *gin.Context) {
	var newScooter Scooter
	if err := c.ShouldBindJSON(&newScooter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scooter data"})
		return
	}

	s.synchronizer.CreateScooter(newScooter.Id)
	//s.scooters[newScooter.Id] = &newScooter

	c.JSON(http.StatusOK, gin.H{"newScooter": newScooter, "responder": getServerInfo()})
}

// create operation ---- scooter_id=92929 ===== ordernum=9
// reserve operation --- scooter_id=282
// release operation .... scooter_id=222, distance =399

func (s *ScooterService) reserveScooter(c *gin.Context) {
	scooterId := c.Param("scooter_id")
	if scooterId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing scooter ID"})
		return
	}

	// Generate a random reservation ID for simplicity
	rand.Seed(time.Now().UnixNano())
	reservationId := "R" + strconv.Itoa(rand.Intn(1000000))

	s.synchronizer.ReserveScooter(scooterId, reservationId)

	c.JSON(http.StatusOK, gin.H{"reservation_id": reservationId, "responder": getServerInfo()})

}

func (s *ScooterService) releaseScooter(c *gin.Context) {
	scooterId := c.Param("scooter_id")
	if scooterId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing scooter ID"})
		return
	}

	var req ReleaseScooterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if req.ReservationId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing reservation ID"})
		return
	}

	s.synchronizer.ReleaseScooter(scooterId, req.ReservationId, req.RideDistance)

	// TODO: replace with with a call to multipaxos
	// For now, just returning the prepared response
	c.JSON(http.StatusOK, gin.H{"status": "ok", "responder": getServerInfo()})
}

func (s *ScooterService) getServers(c *gin.Context) {
	resp, err := etcdClient.Get(c, "/servers/", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve servers from etcd"})
		log.Printf("Failed to get servers from etcd: %v", err)
		return
	}
	var servers []string
	for _, ev := range resp.Kvs {
		// Assuming server names are stored as values in etcd
		servers = append(servers, string(ev.Value))
	}
	c.JSON(http.StatusOK, gin.H{"servers": servers, "responder": getServerInfo()})
}

func (s *ScooterService) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok", "responder": getServerInfo()})
}

func (s *ScooterService) RegisterRoutes(router *gin.Engine) {
	router.GET("/scooters", s.getScooters)
	router.PUT("/scooters/:id", s.updateScooter)
	router.POST("/scooters/:scooter_id/reservations", s.reserveScooter)
	router.POST("/scooters/:scooter_id/releases", s.releaseScooter)
	router.GET("/servers", s.getServers)
	router.GET("/health", s.health)
}

func startScooterService(stopCh chan struct{}, etcdClient *clientv3.Client, scooters map[string]*Scooter, synchronizer *Synchronizer) {
	router := gin.Default()
	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Allow all origins
		//AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "http://localhost:4200"
		//},
		MaxAge: 12 * time.Hour,
	}))
	scooterService := ScooterService{etcdClient, scooters, synchronizer}
	scooterService.RegisterRoutes(router)

	go func() {
		if err := router.Run(":50053"); err != nil {
			log.Fatalf("Failed to run scooter server: %v", err)
		}
	}()
	log.Println("Scooter http server listening on port 50053")

	// Listen for stop signal
	<-stopCh
	log.Println("Shutting down scooter http server...")
	// Additional shutdown logic goes here, e.g., gracefully stopping HTTP server
}
