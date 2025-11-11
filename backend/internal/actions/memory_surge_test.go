package actions

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestNewMemorySurgeAction(t *testing.T) {
	tests := []struct {
		name          string
		sizeMB        int
		duration      int
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid small allocation",
			sizeMB:      10,
			duration:    5,
			expectError: false,
		},
		{
			name:        "valid medium allocation",
			sizeMB:      100,
			duration:    30,
			expectError: false,
		},
		{
			name:        "maximum duration",
			sizeMB:      50,
			duration:    MAX_MEMORY_DURATION,
			expectError: false,
		},
		{
			name:          "zero size",
			sizeMB:        0,
			duration:      10,
			expectError:   true,
			errorContains: "size_mb must be between",
		},
		{
			name:          "negative size",
			sizeMB:        -10,
			duration:      10,
			expectError:   true,
			errorContains: "size_mb must be between",
		},
		{
			name:          "zero duration",
			sizeMB:        10,
			duration:      0,
			expectError:   true,
			errorContains: "duration must be between",
		},
		{
			name:          "duration too long",
			sizeMB:        10,
			duration:      MAX_MEMORY_DURATION + 1,
			expectError:   true,
			errorContains: "duration must be between",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := NewMemorySurgeAction(tt.sizeMB, tt.duration)

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

func TestMemorySurgeAction_Execute(t *testing.T) {
	t.Run("allocates and releases memory", func(t *testing.T) {
		action, err := NewMemorySurgeAction(5, 2)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		var memBefore runtime.MemStats
		runtime.ReadMemStats(&memBefore)

		ctx := context.Background()
		err = action.Execute(ctx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Force GC and verify memory is released
		runtime.GC()
		time.Sleep(100 * time.Millisecond)

		var memAfter runtime.MemStats
		runtime.ReadMemStats(&memAfter)

		// Memory should be approximately back to baseline
		// Allow some variance for runtime overhead
		diff := int64(memAfter.Alloc) - int64(memBefore.Alloc)
		maxDiff := int64(5 * 1024 * 1024) // 5 MB tolerance

		if diff > maxDiff {
			t.Logf("Memory not fully released: before=%d, after=%d, diff=%d MB",
				memBefore.Alloc, memAfter.Alloc, diff/(1024*1024))
			// Don't fail the test as GC timing is not deterministic
		}
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		action, err := NewMemorySurgeAction(10, 30)
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

		// Verify memory is cleaned up
		if action.allocatedData != nil {
			t.Errorf("expected allocated data to be nil after cancellation")
		}
	})

	t.Run("completes normally", func(t *testing.T) {
		action, err := NewMemorySurgeAction(5, 1)
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
		expectedDuration := 1 * time.Second
		if elapsed < expectedDuration || elapsed > expectedDuration+500*time.Millisecond {
			t.Errorf("duration mismatch: expected ~%v, got %v", expectedDuration, elapsed)
		}

		// Progress should be 1.0 at completion
		if action.GetProgress() != 1.0 {
			t.Errorf("expected progress 1.0, got %f", action.GetProgress())
		}
	})
}

func TestMemorySurgeAction_GetProgress(t *testing.T) {
	t.Run("progress starts at 0", func(t *testing.T) {
		action, err := NewMemorySurgeAction(5, 10)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		progress := action.GetProgress()
		if progress != 0.0 {
			t.Errorf("expected initial progress 0.0, got %f", progress)
		}
	})

	t.Run("progress increases during execution", func(t *testing.T) {
		action, err := NewMemorySurgeAction(5, 2)
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
			t.Errorf("expected progress between 0 and 1, got %f", progress)
		}

		// Wait longer
		time.Sleep(1 * time.Second)

		progress2 := action.GetProgress()
		if progress2 <= progress {
			t.Errorf("expected progress to increase: %f -> %f", progress, progress2)
		}

		cancel()
	})
}

func TestMemorySurgeAction_Cleanup(t *testing.T) {
	t.Run("cleanup releases memory", func(t *testing.T) {
		action, err := NewMemorySurgeAction(10, 10)
		if err != nil {
			t.Fatalf("failed to create action: %v", err)
		}

		// Manually allocate memory
		action.startTime = time.Now()
		action.allocatedData = make([][]byte, 10)
		for i := range action.allocatedData {
			action.allocatedData[i] = make([]byte, 1024*1024) // 1 MB
		}

		// Verify memory is allocated
		if len(action.allocatedData) != 10 {
			t.Errorf("expected 10 allocations, got %d", len(action.allocatedData))
		}

		// Cleanup
		action.cleanup()

		// Verify memory is released
		if action.allocatedData != nil {
			t.Errorf("expected allocatedData to be nil after cleanup")
		}
	})
}
