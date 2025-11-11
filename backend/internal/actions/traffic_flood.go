package actions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// TrafficFloodAction generates network traffic by making HTTP requests
// It creates HTTP requests to a dummy endpoint at a specified rate
//
// Safety: Limited request rate and duration
// Cancellation: Responds to context cancellation within 100ms
// Cleanup: Closes all HTTP connections properly
type TrafficFloodAction struct {
	requestsPerSec int
	duration       time.Duration
	targetURL      string
	startTime      time.Time
	completedReqs  atomic.Int64
	totalReqs      int64
	client         *http.Client
	mu             sync.RWMutex
}

// NewTrafficFloodAction creates a new traffic flood action
func NewTrafficFloodAction(requestsPerSec int, durationSeconds int, targetURL string) (*TrafficFloodAction, error) {
	// Validate inputs
	if requestsPerSec < 1 || requestsPerSec > 1000 {
		return nil, fmt.Errorf("requests_per_sec must be between 1 and 1000, got %d", requestsPerSec)
	}

	if durationSeconds < 1 || durationSeconds > 60 {
		return nil, fmt.Errorf("duration must be between 1 and 60 seconds, got %d", durationSeconds)
	}

	if targetURL == "" {
		// Default to localhost dummy endpoint
		targetURL = "http://localhost:8080/api/dummy"
	}

	duration := time.Duration(durationSeconds) * time.Second
	totalReqs := int64(requestsPerSec * durationSeconds)

	return &TrafficFloodAction{
		requestsPerSec: requestsPerSec,
		duration:       duration,
		targetURL:      targetURL,
		totalReqs:      totalReqs,
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}, nil
}

// Execute runs the traffic flood action
func (a *TrafficFloodAction) Execute(ctx context.Context) error {
	a.mu.Lock()
	a.startTime = time.Now()
	a.mu.Unlock()

	// Calculate interval between requests
	interval := time.Second / time.Duration(a.requestsPerSec)

	// Create ticker for request pacing
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Create timer for duration
	endTime := a.startTime.Add(a.duration)

	// Channel to limit concurrent requests
	semaphore := make(chan struct{}, 50) // Max 50 concurrent requests

	// WaitGroup to track in-flight requests
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			// Wait for in-flight requests to complete
			wg.Wait()
			return ctx.Err()

		case <-ticker.C:
			if time.Now().After(endTime) {
				// Wait for in-flight requests to complete
				wg.Wait()
				return nil
			}

			// Send request
			semaphore <- struct{}{} // Acquire semaphore
			wg.Add(1)

			go func() {
				defer wg.Done()
				defer func() { <-semaphore }() // Release semaphore

				a.sendRequest(ctx)
			}()
		}
	}
}

// sendRequest makes a single HTTP request
func (a *TrafficFloodAction) sendRequest(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.targetURL, nil)
	if err != nil {
		return
	}

	resp, err := a.client.Do(req)
	if err != nil {
		// Ignore errors (endpoint might not exist, but we're generating traffic)
		a.completedReqs.Add(1)
		return
	}

	// Read and discard response body to complete the request
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	a.completedReqs.Add(1)
}

// GetProgress returns the current progress (0.0 to 1.0)
func (a *TrafficFloodAction) GetProgress() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.startTime.IsZero() {
		return 0.0
	}

	elapsed := time.Since(a.startTime)
	if elapsed >= a.duration {
		return 1.0
	}

	return float64(elapsed) / float64(a.duration)
}
