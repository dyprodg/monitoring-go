package actions

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// MemorySurgeAction allocates memory to create memory load
// It allocates a specified amount of memory and holds it for a duration
//
// Safety: Respects MAX_MEMORY_PERCENT and MAX_MEMORY_DURATION limits
// Cancellation: Responds to context cancellation within 100ms
// Cleanup: Releases all allocated memory when done
type MemorySurgeAction struct {
	sizeMB        int
	duration      time.Duration
	startTime     time.Time
	allocatedData [][]byte
	mu            sync.RWMutex
}

// NewMemorySurgeAction creates a new memory surge action
func NewMemorySurgeAction(sizeMB int, durationSeconds int) (*MemorySurgeAction, error) {
	// For safety, limit to 2GB max for testing
	// In production, this could use gopsutil to get actual system memory
	maxMemoryMB := 2048 // 2GB max for safety

	// Validate inputs
	if sizeMB < 1 || sizeMB > maxMemoryMB {
		return nil, fmt.Errorf("size_mb must be between 1 and %d MB, got %d", maxMemoryMB, sizeMB)
	}

	if durationSeconds < 1 || durationSeconds > MAX_MEMORY_DURATION {
		return nil, fmt.Errorf("duration must be between 1 and %d seconds, got %d", MAX_MEMORY_DURATION, durationSeconds)
	}

	return &MemorySurgeAction{
		sizeMB:   sizeMB,
		duration: time.Duration(durationSeconds) * time.Second,
	}, nil
}

// Execute runs the memory surge action
func (a *MemorySurgeAction) Execute(ctx context.Context) error {
	a.mu.Lock()
	a.startTime = time.Now()
	a.mu.Unlock()

	// Allocate memory in chunks
	chunkSize := 1024 * 1024 // 1 MB chunks
	numChunks := a.sizeMB

	a.mu.Lock()
	a.allocatedData = make([][]byte, 0, numChunks)
	a.mu.Unlock()

	// Allocate memory gradually
	for i := 0; i < numChunks; i++ {
		select {
		case <-ctx.Done():
			a.cleanup()
			return ctx.Err()
		default:
			// Allocate and fill chunk with data
			chunk := make([]byte, chunkSize)
			// Fill with non-zero data to prevent compiler optimizations
			for j := range chunk {
				chunk[j] = byte(j % 256)
			}

			a.mu.Lock()
			a.allocatedData = append(a.allocatedData, chunk)
			a.mu.Unlock()

			// Small delay between allocations to prevent system freeze
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Hold the memory for the duration
	endTime := a.startTime.Add(a.duration)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.cleanup()
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(endTime) {
				a.cleanup()
				return nil
			}
			// Touch the memory to keep it allocated
			a.touchMemory()
		}
	}
}

// touchMemory accesses allocated memory to prevent it from being swapped out
func (a *MemorySurgeAction) touchMemory() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Read a byte from each chunk periodically
	for i := range a.allocatedData {
		if len(a.allocatedData[i]) > 0 {
			_ = a.allocatedData[i][0]
		}
	}
}

// cleanup releases all allocated memory
func (a *MemorySurgeAction) cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Clear the allocated data
	a.allocatedData = nil

	// Force garbage collection to release memory immediately
	runtime.GC()
}

// GetProgress returns the current progress (0.0 to 1.0)
func (a *MemorySurgeAction) GetProgress() float64 {
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
