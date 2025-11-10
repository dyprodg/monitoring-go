package metrics

import (
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	collector := NewCollector()
	if collector == nil {
		t.Fatal("NewCollector() returned nil")
	}
}

func TestCollectorStartAndGetCurrent(t *testing.T) {
	collector := NewCollector()

	// Start collector
	collector.Start(100 * time.Millisecond)

	// Wait a bit for initial collection
	time.Sleep(150 * time.Millisecond)

	// Get current metrics
	metrics := collector.GetCurrent()

	// Verify metrics are collected
	if metrics.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}

	// CPU should be between 0 and 100
	if metrics.CPU < 0 || metrics.CPU > 100 {
		t.Errorf("CPU should be between 0 and 100, got: %f", metrics.CPU)
	}

	// Memory should be between 0 and 100
	if metrics.Memory < 0 || metrics.Memory > 100 {
		t.Errorf("Memory should be between 0 and 100, got: %f", metrics.Memory)
	}
}

func TestCollectorConcurrentAccess(t *testing.T) {
	collector := NewCollector()
	collector.Start(50 * time.Millisecond)

	// Simulate concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = collector.GetCurrent()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we get here without a race condition, the test passes
	metrics := collector.GetCurrent()
	if metrics.CPU < 0 || metrics.CPU > 100 {
		t.Errorf("CPU should be between 0 and 100, got: %f", metrics.CPU)
	}
}

func TestCollectorMetricsUpdate(t *testing.T) {
	collector := NewCollector()
	// Use longer interval since CPU collection takes ~1 second
	collector.Start(2 * time.Second)

	// Get initial metrics (wait for first collection to complete)
	time.Sleep(1500 * time.Millisecond)
	firstMetrics := collector.GetCurrent()
	firstTime := firstMetrics.Timestamp

	// Wait for next update cycle (2s interval + 1s collection time)
	time.Sleep(3500 * time.Millisecond)
	secondMetrics := collector.GetCurrent()
	secondTime := secondMetrics.Timestamp

	// Timestamps should be different
	if firstTime.Equal(secondTime) {
		t.Errorf("Metrics should update over time. First: %v, Second: %v", firstTime, secondTime)
	}

	// Second timestamp should be after first
	if !secondTime.After(firstTime) {
		t.Errorf("Second timestamp should be after first. First: %v, Second: %v", firstTime, secondTime)
	}

	// Verify both metrics are valid
	if firstMetrics.CPU < 0 || firstMetrics.CPU > 100 {
		t.Errorf("First CPU reading invalid: %f", firstMetrics.CPU)
	}
	if secondMetrics.CPU < 0 || secondMetrics.CPU > 100 {
		t.Errorf("Second CPU reading invalid: %f", secondMetrics.CPU)
	}
}
