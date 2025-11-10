package actions

import (
	"context"
	"testing"
	"time"
)

func TestNewCPUStressAction(t *testing.T) {
	tests := []struct {
		name            string
		targetPercent   int
		durationSeconds int
		wantError       bool
	}{
		{
			name:            "valid action",
			targetPercent:   80,
			durationSeconds: 10,
			wantError:       false,
		},
		{
			name:            "minimum valid",
			targetPercent:   1,
			durationSeconds: 1,
			wantError:       false,
		},
		{
			name:            "maximum valid",
			targetPercent:   MAX_CPU_PERCENT,
			durationSeconds: MAX_CPU_DURATION,
			wantError:       false,
		},
		{
			name:            "CPU too high",
			targetPercent:   96,
			durationSeconds: 10,
			wantError:       true,
		},
		{
			name:            "CPU negative",
			targetPercent:   -1,
			durationSeconds: 10,
			wantError:       true,
		},
		{
			name:            "duration too long",
			targetPercent:   80,
			durationSeconds: 31,
			wantError:       true,
		},
		{
			name:            "duration zero",
			targetPercent:   80,
			durationSeconds: 0,
			wantError:       true,
		},
		{
			name:            "duration negative",
			targetPercent:   80,
			durationSeconds: -1,
			wantError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action, err := NewCPUStressAction(tt.targetPercent, tt.durationSeconds)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if action != nil {
					t.Error("Expected nil action on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if action == nil {
					t.Error("Expected action, got nil")
				}
				if action != nil {
					if action.targetPercent != tt.targetPercent {
						t.Errorf("Expected targetPercent %d, got %d", tt.targetPercent, action.targetPercent)
					}
					if action.duration != time.Duration(tt.durationSeconds)*time.Second {
						t.Errorf("Expected duration %ds, got %v", tt.durationSeconds, action.duration)
					}
				}
			}
		})
	}
}

func TestCPUStressAction_Execute(t *testing.T) {
	action, err := NewCPUStressAction(50, 1)
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}

	ctx := context.Background()
	err = action.Execute(ctx)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Verify progress is complete
	if action.GetProgress() < 0.99 {
		t.Errorf("Expected progress ~1.0, got %f", action.GetProgress())
	}
}

func TestCPUStressAction_ExecuteWithCancellation(t *testing.T) {
	action, err := NewCPUStressAction(50, 10)
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start execution in goroutine
	done := make(chan error)
	go func() {
		done <- action.Execute(ctx)
	}()

	// Cancel after 500ms
	time.Sleep(500 * time.Millisecond)
	cancel()

	// Wait for execution to finish
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Execute() did not respond to cancellation within 2 seconds")
	}

	// Verify progress is incomplete
	progress := action.GetProgress()
	if progress >= 1.0 {
		t.Errorf("Expected progress <1.0 after cancellation, got %f", progress)
	}
	if progress <= 0.0 {
		t.Error("Expected some progress >0.0 after 500ms")
	}
}

func TestCPUStressAction_GetProgress(t *testing.T) {
	action, err := NewCPUStressAction(50, 2)
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}

	// Progress should be 0 before execution
	if action.GetProgress() != 0.0 {
		t.Errorf("Expected initial progress 0.0, got %f", action.GetProgress())
	}

	// Start execution
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go action.Execute(ctx)

	// Wait and check progress increases
	time.Sleep(600 * time.Millisecond)
	progress1 := action.GetProgress()

	time.Sleep(600 * time.Millisecond)
	progress2 := action.GetProgress()

	if progress1 <= 0.0 {
		t.Error("Progress should be > 0 after 600ms")
	}

	if progress2 <= progress1 {
		t.Errorf("Progress should increase over time: %f -> %f", progress1, progress2)
	}

	if progress1 > 1.0 || progress2 > 1.0 {
		t.Error("Progress should never exceed 1.0")
	}

	// Cancel to clean up
	cancel()
}

func TestCPUStressAction_StopsWithin1Second(t *testing.T) {
	// CRITICAL SAFETY TEST: Action must stop within 1 second
	action, err := NewCPUStressAction(80, 30)
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start execution
	done := make(chan error)
	go func() {
		done <- action.Execute(ctx)
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Cancel and measure time to stop
	startCancel := time.Now()
	cancel()

	select {
	case <-done:
		stopDuration := time.Since(startCancel)
		if stopDuration > time.Second {
			t.Errorf("CRITICAL: Action took %v to stop (must be ≤1s)", stopDuration)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("CRITICAL: Action did not stop within 2 seconds")
	}
}

func TestCPUStressAction_CompletesOnTime(t *testing.T) {
	durationSeconds := 1
	action, err := NewCPUStressAction(50, durationSeconds)
	if err != nil {
		t.Fatalf("Failed to create action: %v", err)
	}

	ctx := context.Background()
	start := time.Now()
	err = action.Execute(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	expectedDuration := time.Duration(durationSeconds) * time.Second
	// Allow 20% tolerance for timing
	minDuration := expectedDuration - 200*time.Millisecond
	maxDuration := expectedDuration + 500*time.Millisecond

	if elapsed < minDuration || elapsed > maxDuration {
		t.Errorf("Expected duration ~%v, got %v", expectedDuration, elapsed)
	}
}

func TestCPUStressAction_ConcurrentSafety(t *testing.T) {
	// Test that multiple concurrent actions don't interfere
	action1, _ := NewCPUStressAction(40, 2)
	action2, _ := NewCPUStressAction(40, 2)

	ctx := context.Background()

	done1 := make(chan error)
	done2 := make(chan error)

	go func() {
		done1 <- action1.Execute(ctx)
	}()

	go func() {
		done2 <- action2.Execute(ctx)
	}()

	// Both should complete without error
	err1 := <-done1
	err2 := <-done2

	if err1 != nil {
		t.Errorf("Action1 error: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Action2 error: %v", err2)
	}

	// Both should have progress = 1.0
	if action1.GetProgress() < 0.99 {
		t.Errorf("Action1 incomplete: progress %f", action1.GetProgress())
	}
	if action2.GetProgress() < 0.99 {
		t.Errorf("Action2 incomplete: progress %f", action2.GetProgress())
	}
}
