package actions

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DiskStormAction generates disk I/O load by performing file operations
// It creates temporary files, writes/reads data, and cleans up
//
// Safety: Respects MAX_DISK_SIZE_MB limit
// Cancellation: Responds to context cancellation within 100ms
// Cleanup: Removes all temporary files when done
type DiskStormAction struct {
	operations    int
	fileSizeKB    int
	startTime     time.Time
	totalOps      int
	completedOps  int
	tempDir       string
	createdFiles  []string
	mu            sync.RWMutex
}

// NewDiskStormAction creates a new disk storm action
func NewDiskStormAction(operations int, fileSizeKB int) (*DiskStormAction, error) {
	// Validate inputs
	if operations < 1 || operations > 10000 {
		return nil, fmt.Errorf("operations must be between 1 and 10000, got %d", operations)
	}

	if fileSizeKB < 1 || fileSizeKB > 1024 {
		return nil, fmt.Errorf("file_size_kb must be between 1 and 1024 KB, got %d", fileSizeKB)
	}

	// Calculate total disk usage and enforce limit
	totalSizeMB := (operations * fileSizeKB) / 1024
	if totalSizeMB > MAX_DISK_SIZE_MB {
		return nil, fmt.Errorf("total disk usage would be %d MB, exceeds limit of %d MB", totalSizeMB, MAX_DISK_SIZE_MB)
	}

	return &DiskStormAction{
		operations:   operations,
		fileSizeKB:   fileSizeKB,
		totalOps:     operations * 3, // read + write + delete = 3 ops per file
		createdFiles: make([]string, 0),
	}, nil
}

// Execute runs the disk storm action
func (a *DiskStormAction) Execute(ctx context.Context) error {
	a.mu.Lock()
	a.startTime = time.Now()
	a.mu.Unlock()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "disk-storm-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	a.mu.Lock()
	a.tempDir = tempDir
	a.mu.Unlock()

	// Ensure cleanup on exit
	defer a.cleanup()

	// Generate random data buffer (reuse for efficiency)
	dataSize := a.fileSizeKB * 1024
	data := make([]byte, dataSize)
	if _, err := rand.Read(data); err != nil {
		return fmt.Errorf("failed to generate random data: %w", err)
	}

	// Perform file operations
	for i := 0; i < a.operations; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := a.performFileOperation(i, data); err != nil {
				return err
			}
		}
	}

	return nil
}

// performFileOperation writes, reads, and deletes a file
func (a *DiskStormAction) performFileOperation(index int, data []byte) error {
	filename := filepath.Join(a.tempDir, fmt.Sprintf("test-file-%d.dat", index))

	// Write operation
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	a.mu.Lock()
	a.createdFiles = append(a.createdFiles, filename)
	a.completedOps++
	a.mu.Unlock()

	// Read operation
	if _, err := os.ReadFile(filename); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	a.mu.Lock()
	a.completedOps++
	a.mu.Unlock()

	// Delete operation
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	a.mu.Lock()
	a.completedOps++
	// Remove from created files list since it's deleted
	for i, f := range a.createdFiles {
		if f == filename {
			a.createdFiles = append(a.createdFiles[:i], a.createdFiles[i+1:]...)
			break
		}
	}
	a.mu.Unlock()

	return nil
}

// cleanup removes all temporary files and directory
func (a *DiskStormAction) cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Remove any remaining files
	for _, file := range a.createdFiles {
		os.Remove(file) // Ignore errors during cleanup
	}

	// Remove temp directory
	if a.tempDir != "" {
		os.RemoveAll(a.tempDir) // Ignore errors during cleanup
	}

	a.createdFiles = nil
}

// GetProgress returns the current progress (0.0 to 1.0)
func (a *DiskStormAction) GetProgress() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.totalOps == 0 {
		return 0.0
	}

	progress := float64(a.completedOps) / float64(a.totalOps)
	if progress > 1.0 {
		return 1.0
	}

	return progress
}
