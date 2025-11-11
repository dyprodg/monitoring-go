package api

import (
	"encoding/json"
	"net/http"
	"time"

	"monitoring-dashboard/internal/actions"
	"monitoring-dashboard/internal/metrics"
	"monitoring-dashboard/pkg/models"

	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for the API
type Handler struct {
	collector *metrics.Collector
	engine    *actions.Engine
}

// NewHandler creates a new API handler
func NewHandler(collector *metrics.Collector, engine *actions.Engine) *Handler {
	return &Handler{
		collector: collector,
		engine:    engine,
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

// CPUStressHandler starts a CPU stress action
func (h *Handler) CPUStressHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CPUStressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create CPU stress action
	cpuAction, err := actions.NewCPUStressAction(req.TargetPercent, req.DurationSeconds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start action
	action, err := h.engine.StartAction(models.ActionTypeCPUStress, cpuAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.ActionResponse{
		ID:        action.ID,
		Status:    string(action.Status),
		StartedAt: action.StartedAt,
		Message:   "CPU stress action started",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetActiveActionsHandler returns all active actions
func (h *Handler) GetActiveActionsHandler(w http.ResponseWriter, r *http.Request) {
	activeActions := h.engine.GetActiveActions()

	response := map[string]interface{}{
		"actions": activeActions,
		"count":   len(activeActions),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MemorySurgeHandler starts a memory surge action
func (h *Handler) MemorySurgeHandler(w http.ResponseWriter, r *http.Request) {
	var req models.MemorySurgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create memory surge action
	memAction, err := actions.NewMemorySurgeAction(req.SizeMB, req.DurationSeconds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start action
	action, err := h.engine.StartAction(models.ActionTypeMemorySurge, memAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.ActionResponse{
		ID:        action.ID,
		Status:    string(action.Status),
		StartedAt: action.StartedAt,
		Message:   "Memory surge action started",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DiskStormHandler starts a disk storm action
func (h *Handler) DiskStormHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DiskStormRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create disk storm action
	diskAction, err := actions.NewDiskStormAction(req.Operations, req.FileSizeKB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start action
	action, err := h.engine.StartAction(models.ActionTypeDiskStorm, diskAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.ActionResponse{
		ID:        action.ID,
		Status:    string(action.Status),
		StartedAt: action.StartedAt,
		Message:   "Disk storm action started",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// TrafficFloodHandler starts a traffic flood action
func (h *Handler) TrafficFloodHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TrafficFloodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create traffic flood action
	trafficAction, err := actions.NewTrafficFloodAction(req.RequestsPerSec, req.DurationSeconds, req.TargetURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start action
	action, err := h.engine.StartAction(models.ActionTypeTrafficFlood, trafficAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.ActionResponse{
		ID:        action.ID,
		Status:    string(action.Status),
		StartedAt: action.StartedAt,
		Message:   "Traffic flood action started",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// StopActionHandler stops a running action
func (h *Handler) StopActionHandler(w http.ResponseWriter, r *http.Request) {
	actionID := chi.URLParam(r, "id")
	if actionID == "" {
		http.Error(w, "Action ID is required", http.StatusBadRequest)
		return
	}

	err := h.engine.StopAction(actionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]string{
		"status":  "stopped",
		"message": "Action stopped successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
