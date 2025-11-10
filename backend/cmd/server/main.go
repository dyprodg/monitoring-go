package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"monitoring-dashboard/internal/actions"
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

	// Initialize action engine
	engine := actions.NewEngine(collector)
	log.Println("Action engine initialized with safety limits:")
	log.Printf("  - Max CPU: %d%%, Critical: %d%%", actions.MAX_CPU_PERCENT, actions.CRITICAL_CPU)
	log.Printf("  - Max Memory: %d%%, Critical: %d%%", actions.MAX_MEMORY_PERCENT, actions.CRITICAL_MEMORY)
	log.Printf("  - Max concurrent actions: %d", actions.MAX_CONCURRENT)

	// Initialize API handler
	handler := api.NewHandler(collector, engine)
	router := handler.SetupRoutes()

	// Start HTTP server
	addr := fmt.Sprintf(":%s", ServerPort)
	log.Printf("Server starting on http://localhost:%s", ServerPort)
	log.Printf("Health check: http://localhost:%s/api/health", ServerPort)
	log.Printf("Metrics endpoint: http://localhost:%s/api/metrics", ServerPort)
	log.Printf("Actions endpoint: http://localhost:%s/api/actions/cpu-stress", ServerPort)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
