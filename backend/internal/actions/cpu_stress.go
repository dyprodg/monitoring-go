package actions

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// CPUStressAction generates CPU load for testing
// It uses busy loops in goroutines to consume CPU cycles
//
// Safety: Respects MAX_CPU_PERCENT and MAX_CPU_DURATION limits
// Cancellation: Responds to context cancellation within 100ms
type CPUStressAction struct {
	targetPercent int
	duration      time.Duration
	startTime     time.Time
	mu            sync.RWMutex
}

// NewCPUStressAction creates a new CPU stress action
func NewCPUStressAction(targetPercent int, durationSeconds int) (*CPUStressAction, error) {
	// Validate inputs
	if targetPercent < 0 || targetPercent > MAX_CPU_PERCENT {
		return nil, fmt.Errorf("target_percent must be between 0 and %d, got %d", MAX_CPU_PERCENT, targetPercent)
	}

	if durationSeconds < 1 || durationSeconds > MAX_CPU_DURATION {
		return nil, fmt.Errorf("duration must be between 1 and %d seconds, got %d", MAX_CPU_DURATION, durationSeconds)
	}

	return &CPUStressAction{
		targetPercent: targetPercent,
		duration:      time.Duration(durationSeconds) * time.Second,
	}, nil
}

// Execute runs the CPU stress action
func (a *CPUStressAction) Execute(ctx context.Context) error {
	a.mu.Lock()
	a.startTime = time.Now()
	a.mu.Unlock()

	// Calculate number of workers based on target percentage and CPU cores
	numCPU := runtime.NumCPU()
	numWorkers := (numCPU * a.targetPercent) / 100
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Create wait group for workers
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			a.stressWorker(ctx)
		}()
	}

	// Wait for workers to complete or context to be cancelled
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		// Wait a bit for workers to stop
		time.Sleep(100 * time.Millisecond)
		return ctx.Err()
	}
}

// stressWorker performs CPU-intensive work
func (a *CPUStressAction) stressWorker(ctx context.Context) {
	endTime := a.startTime.Add(a.duration)

	for {
		// Check if we should stop
		select {
		case <-ctx.Done():
			return
		default:
			// Check if duration has elapsed
			if time.Now().After(endTime) {
				return
			}

			// Perform CPU-intensive work
			// This is a busy loop that consumes CPU cycles
			sum := 0
			for i := 0; i < 1000000; i++ {
				sum += i
			}
		}
	}
}

// GetProgress returns the current progress (0.0 to 1.0)
func (a *CPUStressAction) GetProgress() float64 {
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
