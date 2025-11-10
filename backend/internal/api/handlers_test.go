package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"monitoring-dashboard/internal/actions"
	"monitoring-dashboard/internal/metrics"
	"monitoring-dashboard/pkg/models"
)

func TestHealthHandler(t *testing.T) {
	// Create collector, engine and handler
	collector := metrics.NewCollector()
	engine := actions.NewEngine(collector)
	handler := NewHandler(collector, engine)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	// Call handler
	handler.HealthHandler(rec, req)

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse response
	var health models.HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&health); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response
	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", health.Status)
	}

	if health.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestMetricsHandler(t *testing.T) {
	// Create collector, engine and handler
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(150 * time.Millisecond) // Wait for metrics collection

	engine := actions.NewEngine(collector)
	handler := NewHandler(collector, engine)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	rec := httptest.NewRecorder()

	// Call handler
	handler.MetricsHandler(rec, req)

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse response
	var metrics models.Metrics
	if err := json.NewDecoder(rec.Body).Decode(&metrics); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify metrics
	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	if metrics.CPU < 0 || metrics.CPU > 100 {
		t.Errorf("CPU should be between 0 and 100, got: %f", metrics.CPU)
	}

	if metrics.Memory < 0 || metrics.Memory > 100 {
		t.Errorf("Memory should be between 0 and 100, got: %f", metrics.Memory)
	}
}

func TestNewHandler(t *testing.T) {
	collector := metrics.NewCollector()
	engine := actions.NewEngine(collector)
	handler := NewHandler(collector, engine)

	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}

	if handler.collector == nil {
		t.Error("Handler collector should not be nil")
	}
}
