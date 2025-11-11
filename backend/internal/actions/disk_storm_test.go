package actions

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestNewDiskStormAction(t *testing.T) {
	tests := []struct {
		name          string
		operations    int
		fileSizeKB    int
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid small storm",
			operations:  10,
			fileSizeKB:  10,
			expectError: false,
		},
		{
			name:        "valid medium storm",
			operations:  100,
			fileSizeKB:  100,
			expectError: false,
		},
		{
			name:        "maximum operations",
			operations:  10000,
			fileSizeKB:  1,
			expectError: false,
		},
		{
			name:          "zero operations",
			operations:    0,
			fileSizeKB:    10,
			expectError:   true,
			errorContains: "operations must be between",
		},
		{
			name:          "too many operations",
			operations:    10001,
			fileSizeKB:    10,
			expectError:   true,
			errorContains: "operations must be between",
		},
		{
			name:          "zero file size",
			operations:    10,
			fileSizeKB:    0,
			expectError:   true,
			errorContains: "file_size_kb must be between",
		},
		{
			name:          "file size too large",
			operations:    10,
			fileSizeKB:    1025,
			expectError:   true,
			errorContains: "file_size_kb must be between",
		},
		{
			name:          "exceeds disk size limit",
			operations:    1000,
			fileSizeKB:    1024,
			expectError:   true,
			errorContains: "exceeds limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := NewDiskStormAction(tt.operations, tt.fileSizeKB)

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
			}
		})
	}
}

func TestDiskStormAction_Execute(t *testing.T) {
	t.Run("creates and cleans up files", func(t *testing.T) {
		action, err := NewDiskStormAction(5, 10)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx := context.Background()
		err = action.Execute(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify temp directory is cleaned up
		if action.tempDir != "" {
			if _, err := os.Stat(action.tempDir); !os.IsNotExist(err) {
				t.Errorf("temp directory still exists: %s", action.tempDir)
			}
		}

		// Verify no files remain
		if len(action.createdFiles) != 0 {
			t.Errorf("expected no files remaining, got %d", len(action.createdFiles))
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		action, err := NewDiskStormAction(1000, 100)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after 100ms
		go func() {
			time.Sleep(100 * time.Millisecond)
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

		// Verify cleanup happened
		if action.tempDir != "" {
			if _, err := os.Stat(action.tempDir); !os.IsNotExist(err) {
				t.Errorf("temp directory not cleaned up: %s", action.tempDir)
			}
		}
	})

	t.Run("completes all operations", func(t *testing.T) {
		numOps := 10
		action, err := NewDiskStormAction(numOps, 10)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx := context.Background()
		err = action.Execute(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Progress should be 1.0 at completion
		progress := action.GetProgress()
		if progress != 1.0 {
			t.Errorf("expected progress 1.0, got %f", progress)
		}

		// Completed ops should be totalOps (operations * 3 for read+write+delete)
		expectedOps := numOps * 3
		if action.completedOps != expectedOps {
			t.Errorf("expected %d completed ops, got %d", expectedOps, action.completedOps)
		}
	})
}

func TestDiskStormAction_GetProgress(t *testing.T) {
	t.Run("progress starts at 0", func(t *testing.T) {
		action, err := NewDiskStormAction(10, 10)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		progress := action.GetProgress()
		if progress != 0.0 {
			t.Errorf("expected initial progress 0.0, got %f", progress)
		}
	})

	t.Run("progress increases during execution", func(t *testing.T) {
		action, err := NewDiskStormAction(100, 100)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go action.Execute(ctx)

		// Wait for some operations to complete
		time.Sleep(200 * time.Millisecond)

		progress := action.GetProgress()
		if progress <= 0.0 {
			t.Errorf("expected progress > 0.0, got %f", progress)
		}

		if progress > 1.0 {
			t.Errorf("expected progress <= 1.0, got %f", progress)
		}

		cancel()
	})
}

func TestDiskStormAction_Cleanup(t *testing.T) {
	t.Run("cleanup removes all files and directory", func(t *testing.T) {
		action, err := NewDiskStormAction(5, 10)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		// Create temp directory manually
		tempDir, err := os.MkdirTemp("", "test-cleanup-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}

		action.tempDir = tempDir

		// Create some test files
		testFile := tempDir + "/test.txt"
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		action.createdFiles = append(action.createdFiles, testFile)

		// Verify directory exists
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			t.Fatalf("temp directory should exist")
		}

		// Cleanup
		action.cleanup()

		// Verify directory is removed
		if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
			t.Errorf("temp directory should be removed")
		}

		// Verify file list is cleared
		if action.createdFiles != nil {
			t.Errorf("expected createdFiles to be nil")
		}
	})
}
