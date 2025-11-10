package api

import (
	"encoding/json"
	"net/http"
	"time"

	"monitoring-dashboard/internal/metrics"
	"monitoring-dashboard/pkg/models"
)

// Handler handles HTTP requests for the API
type Handler struct {
	collector *metrics.Collector
}

// NewHandler creates a new API handler
func NewHandler(collector *metrics.Collector) *Handler {
	return &Handler{
		collector: collector,
	}
}

// HealthHandler returns the health status of the service
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	health := models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// MetricsHandler returns current system metrics
func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := h.collector.GetCurrent()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
