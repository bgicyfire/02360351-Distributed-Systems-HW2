package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type ScooterService struct {
	// Add any fields here if needed, e.g., a database connection
}

type ScooterStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (s *ScooterService) getScooterStatus(c *gin.Context) {
	scooterID := c.Query("id")
	if scooterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	status := ScooterStatus{
		ID:     scooterID,
		Status: "Available",
	}

	c.JSON(http.StatusOK, status)
}

func (s *ScooterService) getServers(c *gin.Context) {
	servers := []string{"a", "b"} // This is the fixed array of strings to be returned
	c.JSON(http.StatusOK, gin.H{"servers": servers})
}

func (s *ScooterService) RegisterRoutes(router *gin.Engine) {
	router.GET("/scooter/status", s.getScooterStatus)
	router.GET("/server", s.getServers)
}

func startScooterService(stopCh chan struct{}) {
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
	scooterService := ScooterService{}
	scooterService.RegisterRoutes(router)

	go func() {
		log.Println("Scooter server listening on port 50053")
		if err := router.Run(":50053"); err != nil {
			log.Fatalf("Failed to run scooter server: %v", err)
		}
	}()

	// Listen for stop signal
	<-stopCh
	log.Println("Shutting down scooter server...")
	// Additional shutdown logic goes here, e.g., gracefully stopping HTTP server
}
