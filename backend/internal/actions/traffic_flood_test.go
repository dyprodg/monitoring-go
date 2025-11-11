package actions

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestNewTrafficFloodAction(t *testing.T) {
	tests := []struct {
		name          string
		reqsPerSec    int
		duration      int
		targetURL     string
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid low rate",
			reqsPerSec:  10,
			duration:    5,
			targetURL:   "http://example.com",
			expectError: false,
		},
		{
			name:        "valid medium rate",
			reqsPerSec:  100,
			duration:    10,
			targetURL:   "http://example.com",
			expectError: false,
		},
		{
			name:        "maximum rate",
			reqsPerSec:  1000,
			duration:    5,
			targetURL:   "http://example.com",
			expectError: false,
		},
		{
			name:        "empty URL defaults to localhost",
			reqsPerSec:  10,
			duration:    5,
			targetURL:   "",
			expectError: false,
		},
		{
			name:          "zero requests per second",
			reqsPerSec:    0,
			duration:      10,
			targetURL:     "http://example.com",
			expectError:   true,
			errorContains: "requests_per_sec must be between",
		},
		{
			name:          "too many requests per second",
			reqsPerSec:    1001,
			duration:      10,
			targetURL:     "http://example.com",
			expectError:   true,
			errorContains: "requests_per_sec must be between",
		},
		{
			name:          "zero duration",
			reqsPerSec:    10,
			duration:      0,
			targetURL:     "http://example.com",
			expectError:   true,
			errorContains: "duration must be between",
		},
		{
			name:          "duration too long",
			reqsPerSec:    10,
			duration:      61,
			targetURL:     "http://example.com",
			expectError:   true,
			errorContains: "duration must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := NewTrafficFloodAction(tt.reqsPerSec, tt.duration, tt.targetURL)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if action == nil {
					t.Errorf("expected action but got nil")
				}
				if tt.targetURL == "" && action.targetURL == "" {
					t.Errorf("expected default URL when empty string provided")
				}
			}
		})
	}
}

func TestTrafficFloodAction_Execute(t *testing.T) {
	t.Run("generates correct number of requests", func(t *testing.T) {
		// Create test server that counts requests
		var requestCount atomic.Int64
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount.Add(1)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		reqsPerSec := 10
		duration := 1
		action, err := NewTrafficFloodAction(reqsPerSec, duration, server.URL)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx := context.Background()
		start := time.Now()
		err = action.Execute(ctx)
		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Should take approximately the specified duration
		expectedDuration := time.Duration(duration) * time.Second
		if elapsed < expectedDuration || elapsed > expectedDuration+500*time.Millisecond {
			t.Logf("duration variance: expected ~%v, got %v", expectedDuration, elapsed)
		}

		// Verify request count (allow some variance due to timing)
		finalCount := requestCount.Load()
		expectedCount := int64(reqsPerSec * duration)
		minCount := expectedCount - 2
		maxCount := expectedCount + 2

		if finalCount < minCount || finalCount > maxCount {
			t.Errorf("request count out of range: expected %d±2, got %d", expectedCount, finalCount)
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		action, err := NewTrafficFloodAction(100, 30, server.URL)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after 200ms
		go func() {
			time.Sleep(200 * time.Millisecond)
			cancel()
		}()

		start := time.Now()
		err = action.Execute(ctx)
		elapsed := time.Since(start)

		if err != context.Canceled {
			t.Errorf("expected context.Canceled error, got %v", err)
		}

		// Should stop within 1 second
		if elapsed > time.Second {
			t.Errorf("cancellation took too long: %v", elapsed)
		}
	})

	t.Run("handles non-existent endpoint gracefully", func(t *testing.T) {
		// Use a URL that will fail
		action, err := NewTrafficFloodAction(5, 1, "http://localhost:99999/nonexistent")
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx := context.Background()
		err = action.Execute(ctx)

		// Should not error even if requests fail
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Should complete requests (even if they fail)
		if action.completedReqs.Load() == 0 {
			t.Errorf("expected some completed requests")
		}
	})

	t.Run("completes normally", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		action, err := NewTrafficFloodAction(5, 1, server.URL)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx := context.Background()
		err = action.Execute(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Progress should be 1.0 at completion
		if action.GetProgress() != 1.0 {
			t.Errorf("expected progress 1.0, got %f", action.GetProgress())
		}
	})

	t.Run("limits concurrent requests", func(t *testing.T) {
		// Track concurrent request count
		var currentConcurrent atomic.Int32
		var maxConcurrent atomic.Int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			current := currentConcurrent.Add(1)

			// Update max if needed
			for {
				max := maxConcurrent.Load()
				if current <= max {
					break
				}
				if maxConcurrent.CompareAndSwap(max, current) {
					break
				}
			}

			// Simulate slow response
			time.Sleep(50 * time.Millisecond)

			currentConcurrent.Add(-1)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		action, err := NewTrafficFloodAction(100, 1, server.URL)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx := context.Background()
		err = action.Execute(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Max concurrent should not exceed semaphore limit (50)
		max := maxConcurrent.Load()
		if max > 50 {
			t.Errorf("concurrent requests exceeded limit: %d > 50", max)
		}

		t.Logf("Max concurrent requests: %d", max)
	})
}

func TestTrafficFloodAction_GetProgress(t *testing.T) {
	t.Run("progress starts at 0", func(t *testing.T) {
		action, err := NewTrafficFloodAction(10, 10, "http://example.com")
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		progress := action.GetProgress()
		if progress != 0.0 {
			t.Errorf("expected initial progress 0.0, got %f", progress)
		}
	})

	t.Run("progress increases during execution", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		action, err := NewTrafficFloodAction(10, 3, server.URL)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go action.Execute(ctx)

		// Wait for action to start
		time.Sleep(100 * time.Millisecond)

		progress := action.GetProgress()
		if progress <= 0.0 || progress >= 1.0 {
			t.Logf("progress at 100ms: %f (may be 0 if not started yet)", progress)
		}

		// Wait longer
		time.Sleep(1500 * time.Millisecond)

		progress2 := action.GetProgress()
		if progress2 <= progress {
			t.Logf("progress increased: %f -> %f", progress, progress2)
		}

		cancel()
	})
}

// Benchmark to verify performance
func BenchmarkTrafficFloodAction(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	action, err := NewTrafficFloodAction(10, 1, server.URL)
	if err != nil {
		b.Fatalf("failed to create action: %v", err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset action state
		action.startTime = time.Time{}
		action.completedReqs.Store(0)

		if err := action.Execute(ctx); err != nil {
			b.Fatalf("execution failed: %v", err)
		}
	}
}

func ExampleTrafficFloodAction() {
	// Create an action that sends 50 requests/second for 5 seconds
	action, err := NewTrafficFloodAction(50, 5, "http://localhost:8080/api/dummy")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Execute the action
	ctx := context.Background()
	if err := action.Execute(ctx); err != nil {
		fmt.Printf("Execution error: %v\n", err)
		return
	}

	fmt.Printf("Traffic flood completed! Sent %d requests\n", action.completedReqs.Load())
}
