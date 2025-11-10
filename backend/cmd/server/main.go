package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"monitoring-dashboard/internal/api"
	"monitoring-dashboard/internal/metrics"
)

const (
	// Server configuration
	ServerPort = "8080"

	// Metrics collection interval
	MetricsInterval = 1 * time.Second
)

func main() {
	log.Println("Starting Interactive System Monitoring Dashboard...")

	// Initialize metrics collector
	collector := metrics.NewCollector()
	collector.Start(MetricsInterval)
	log.Printf("Metrics collector started (interval: %v)", MetricsInterval)

	// Initialize API handler
	handler := api.NewHandler(collector)
	router := handler.SetupRoutes()

	// Start HTTP server
	addr := fmt.Sprintf(":%s", ServerPort)
	log.Printf("Server starting on http://localhost:%s", ServerPort)
	log.Printf("Health check: http://localhost:%s/api/health", ServerPort)
	log.Printf("Metrics endpoint: http://localhost:%s/api/metrics", ServerPort)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
