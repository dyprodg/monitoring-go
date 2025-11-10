package actions

import (
	"context"
	"errors"
	"testing"
	"time"

	"monitoring-dashboard/internal/metrics"
	"monitoring-dashboard/pkg/models"
)

// MockExecutor is a mock action executor for testing
type MockExecutor struct {
	executeCalled bool
	progress      float64
	shouldError   bool
	duration      time.Duration
}

func (m *MockExecutor) Execute(ctx context.Context) error {
	m.executeCalled = true

	if m.duration > 0 {
		select {
		case <-time.After(m.duration):
			m.progress = 1.0
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if m.shouldError {
		return errors.New("mock error")
	}
	return nil
}

func (m *MockExecutor) GetProgress() float64 {
	return m.progress
}

func TestNewEngine(t *testing.T) {
	collector := metrics.NewCollector()
	engine := NewEngine(collector)

	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}

	if engine.collector == nil {
		t.Error("Engine collector should not be nil")
	}

	if engine.actions == nil {
		t.Error("Engine actions map should not be nil")
	}
}

func TestStartAction(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond) // Wait for initial metrics

	engine := NewEngine(collector)
	executor := &MockExecutor{duration: 100 * time.Millisecond}

	action, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	if action == nil {
		t.Fatal("StartAction() returned nil action")
	}

	if action.ID == "" {
		t.Error("Action ID should not be empty")
	}

	if action.Type != models.ActionTypeCPUStress {
		t.Errorf("Expected action type %s, got %s", models.ActionTypeCPUStress, action.Type)
	}

	if action.Status != models.ActionStatusRunning {
		t.Errorf("Expected status %s, got %s", models.ActionStatusRunning, action.Status)
	}

	if action.StartedAt.IsZero() {
		t.Error("Action StartedAt should not be zero")
	}

	// Wait for action to complete
	time.Sleep(200 * time.Millisecond)

	// Verify executor was called
	if !executor.executeCalled {
		t.Error("Executor Execute() was not called")
	}
}

func TestStartAction_MaxConcurrentReached(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)

	// Start MAX_CONCURRENT actions
	for i := 0; i < MAX_CONCURRENT; i++ {
		executor := &MockExecutor{duration: 2 * time.Second}
		_, err := engine.StartAction(models.ActionTypeCPUStress, executor)
		if err != nil {
			t.Fatalf("Failed to start action %d: %v", i, err)
		}
	}

	// Try to start one more (should fail)
	executor := &MockExecutor{duration: 2 * time.Second}
	_, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if !errors.Is(err, ErrMaxConcurrentReached) {
		t.Errorf("Expected ErrMaxConcurrentReached, got %v", err)
	}
}

func TestStartAction_CPULimitCheck(t *testing.T) {
	// Note: This test may not trigger the limit if actual CPU is low
	// In real testing, we'd need to actually spike the CPU first
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)
	executor := &MockExecutor{duration: 100 * time.Millisecond}

	// Should succeed if CPU is below threshold
	action, err := engine.StartAction(models.ActionTypeCPUStress, executor)

	// If CPU is high, we expect an error
	if err != nil && !errors.Is(err, ErrCPULimitExceeded) {
		t.Errorf("Unexpected error type: %v", err)
	}

	// If no error, action should be valid
	if err == nil && action == nil {
		t.Error("Action should not be nil when no error")
	}
}

func TestStopAction(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)
	executor := &MockExecutor{duration: 5 * time.Second}

	action, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Stop the action
	err = engine.StopAction(action.ID)
	if err != nil {
		t.Errorf("StopAction() error = %v", err)
	}

	// Wait for cancellation to propagate
	time.Sleep(200 * time.Millisecond)

	// Get the action and check status
	stoppedAction, err := engine.GetAction(action.ID)
	if err != nil {
		t.Fatalf("GetAction() error = %v", err)
	}

	if stoppedAction.Status != models.ActionStatusStopped {
		t.Errorf("Expected status %s, got %s", models.ActionStatusStopped, stoppedAction.Status)
	}
}

func TestStopAction_NotFound(t *testing.T) {
	collector := metrics.NewCollector()
	engine := NewEngine(collector)

	err := engine.StopAction("non-existent-id")
	if !errors.Is(err, ErrActionNotFound) {
		t.Errorf("Expected ErrActionNotFound, got %v", err)
	}
}

func TestGetAction(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)
	executor := &MockExecutor{duration: 100 * time.Millisecond}

	startedAction, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	retrievedAction, err := engine.GetAction(startedAction.ID)
	if err != nil {
		t.Fatalf("GetAction() error = %v", err)
	}

	if retrievedAction.ID != startedAction.ID {
		t.Errorf("Expected ID %s, got %s", startedAction.ID, retrievedAction.ID)
	}
}

func TestGetAction_NotFound(t *testing.T) {
	collector := metrics.NewCollector()
	engine := NewEngine(collector)

	_, err := engine.GetAction("non-existent-id")
	if !errors.Is(err, ErrActionNotFound) {
		t.Errorf("Expected ErrActionNotFound, got %v", err)
	}
}

func TestGetActiveActions(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)

	// Initially no active actions
	active := engine.GetActiveActions()
	if len(active) != 0 {
		t.Errorf("Expected 0 active actions, got %d", len(active))
	}

	// Start 2 actions
	executor1 := &MockExecutor{duration: 2 * time.Second}
	executor2 := &MockExecutor{duration: 2 * time.Second}

	_, err := engine.StartAction(models.ActionTypeCPUStress, executor1)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	_, err = engine.StartAction(models.ActionTypeMemorySurge, executor2)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	// Should have 2 active actions
	active = engine.GetActiveActions()
	if len(active) != 2 {
		t.Errorf("Expected 2 active actions, got %d", len(active))
	}
}

func TestActionCompletes(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)
	executor := &MockExecutor{duration: 100 * time.Millisecond}

	action, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	// Wait for action to complete
	time.Sleep(300 * time.Millisecond)

	completedAction, err := engine.GetAction(action.ID)
	if err != nil {
		t.Fatalf("GetAction() error = %v", err)
	}

	if completedAction.Status != models.ActionStatusCompleted {
		t.Errorf("Expected status %s, got %s", models.ActionStatusCompleted, completedAction.Status)
	}

	if completedAction.CompletedAt == nil {
		t.Error("CompletedAt should not be nil")
	}

	if completedAction.Progress != 1.0 {
		t.Errorf("Expected progress 1.0, got %f", completedAction.Progress)
	}
}

func TestActionFails(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)
	executor := &MockExecutor{shouldError: true}

	action, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	// Wait for action to fail
	time.Sleep(200 * time.Millisecond)

	failedAction, err := engine.GetAction(action.ID)
	if err != nil {
		t.Fatalf("GetAction() error = %v", err)
	}

	if failedAction.Status != models.ActionStatusFailed {
		t.Errorf("Expected status %s, got %s", models.ActionStatusFailed, failedAction.Status)
	}

	if failedAction.Error == "" {
		t.Error("Error message should not be empty for failed action")
	}
}

func TestCleanup(t *testing.T) {
	collector := metrics.NewCollector()
	collector.Start(100 * time.Millisecond)
	time.Sleep(1500 * time.Millisecond)

	engine := NewEngine(collector)
	executor := &MockExecutor{duration: 50 * time.Millisecond}

	action, err := engine.StartAction(models.ActionTypeCPUStress, executor)
	if err != nil {
		t.Fatalf("StartAction() error = %v", err)
	}

	// Wait for action to complete
	time.Sleep(200 * time.Millisecond)

	// Cleanup shouldn't remove recent actions (< 1 minute old)
	engine.Cleanup()

	_, err = engine.GetAction(action.ID)
	if err != nil {
		t.Error("Action should still exist after cleanup (not old enough)")
	}
}
