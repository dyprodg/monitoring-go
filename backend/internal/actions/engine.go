package actions

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"monitoring-dashboard/internal/metrics"
	"monitoring-dashboard/pkg/models"

	"github.com/google/uuid"
)

// Safety limits (MANDATORY - NEVER VIOLATE)
const (
	// Maximum limits
	MAX_CPU_PERCENT     = 95  // NEVER exceed 95% CPU
	MAX_CPU_DURATION    = 30  // Max 30 seconds per action
	MAX_MEMORY_PERCENT  = 25  // Max 25% of total RAM
	MAX_MEMORY_DURATION = 60  // Max 60 seconds
	MAX_DISK_SIZE_MB    = 100 // Max 100MB temp files
	MAX_CONCURRENT      = 5   // Max 5 actions simultaneously

	// Emergency shutdown thresholds
	CRITICAL_CPU    = 98 // Kill action immediately
	CRITICAL_MEMORY = 95 // Kill action immediately
)

var (
	ErrMaxConcurrentReached = errors.New("maximum concurrent actions reached")
	ErrCPULimitExceeded     = errors.New("CPU limit exceeded")
	ErrMemoryLimitExceeded  = errors.New("memory limit exceeded")
	ErrDurationExceeded     = errors.New("duration limit exceeded")
	ErrActionNotFound       = errors.New("action not found")
)

// ActionExecutor defines the interface for executable actions
type ActionExecutor interface {
	Execute(ctx context.Context) error
	GetProgress() float64
}

// Engine manages action execution with safety limits
type Engine struct {
	mu         sync.RWMutex
	actions    map[string]*actionContext
	collector  *metrics.Collector
	cancelFuncs map[string]context.CancelFunc
}

// actionContext holds the context for a running action
type actionContext struct {
	action   *models.Action
	executor ActionExecutor
	cancel   context.CancelFunc
}

// NewEngine creates a new action engine
func NewEngine(collector *metrics.Collector) *Engine {
	return &Engine{
		actions:     make(map[string]*actionContext),
		collector:   collector,
		cancelFuncs: make(map[string]context.CancelFunc),
	}
}

// StartAction starts a new action with safety checks
func (e *Engine) StartAction(actionType models.ActionType, executor ActionExecutor) (*models.Action, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check concurrent action limit
	if len(e.actions) >= MAX_CONCURRENT {
		return nil, ErrMaxConcurrentReached
	}

	// Check current system metrics for safety
	currentMetrics := e.collector.GetCurrent()
	if currentMetrics.CPU > float64(MAX_CPU_PERCENT-10) {
		return nil, fmt.Errorf("%w: current CPU %.1f%% too high", ErrCPULimitExceeded, currentMetrics.CPU)
	}
	if currentMetrics.Memory > float64(MAX_MEMORY_PERCENT+50) {
		return nil, fmt.Errorf("%w: current memory %.1f%% too high", ErrMemoryLimitExceeded, currentMetrics.Memory)
	}

	// Create action
	action := &models.Action{
		ID:        uuid.New().String(),
		Type:      actionType,
		Status:    models.ActionStatusStarting,
		StartedAt: time.Now(),
		Progress:  0.0,
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Store action context
	e.actions[action.ID] = &actionContext{
		action:   action,
		executor: executor,
		cancel:   cancel,
	}

	// Start action in goroutine
	go e.runAction(ctx, action.ID)

	// Start safety monitor
	go e.monitorSafety(ctx, action.ID)

	action.Status = models.ActionStatusRunning
	return action, nil
}

// runAction executes an action
func (e *Engine) runAction(ctx context.Context, actionID string) {
	e.mu.RLock()
	actionCtx, exists := e.actions[actionID]
	e.mu.RUnlock()

	if !exists {
		return
	}

	// Execute the action
	err := actionCtx.executor.Execute(ctx)

	// Update action status
	e.mu.Lock()
	defer e.mu.Unlock()

	if actionCtx, exists := e.actions[actionID]; exists {
		now := time.Now()
		actionCtx.action.CompletedAt = &now

		if err != nil {
			if errors.Is(err, context.Canceled) {
				actionCtx.action.Status = models.ActionStatusStopped
			} else {
				actionCtx.action.Status = models.ActionStatusFailed
				actionCtx.action.Error = err.Error()
			}
		} else {
			actionCtx.action.Status = models.ActionStatusCompleted
			actionCtx.action.Progress = 1.0
		}
	}
}

// monitorSafety monitors system metrics and performs emergency shutdown if needed
func (e *Engine) monitorSafety(ctx context.Context, actionID string) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := e.collector.GetCurrent()

			// Emergency shutdown conditions
			if metrics.CPU >= CRITICAL_CPU {
				e.StopAction(actionID)
				return
			}
			if metrics.Memory >= CRITICAL_MEMORY {
				e.StopAction(actionID)
				return
			}

			// Update progress
			e.mu.Lock()
			if actionCtx, exists := e.actions[actionID]; exists {
				actionCtx.action.Progress = actionCtx.executor.GetProgress()
			}
			e.mu.Unlock()
		}
	}
}

// StopAction stops a running action
func (e *Engine) StopAction(actionID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	actionCtx, exists := e.actions[actionID]
	if !exists {
		return ErrActionNotFound
	}

	// Cancel the action context
	actionCtx.cancel()

	return nil
}

// GetAction returns an action by ID
func (e *Engine) GetAction(actionID string) (*models.Action, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	actionCtx, exists := e.actions[actionID]
	if !exists {
		return nil, ErrActionNotFound
	}

	return actionCtx.action, nil
}

// GetActiveActions returns all currently active actions
func (e *Engine) GetActiveActions() []*models.Action {
	e.mu.RLock()
	defer e.mu.RUnlock()

	active := make([]*models.Action, 0, len(e.actions))
	for _, actionCtx := range e.actions {
		if actionCtx.action.Status == models.ActionStatusStarting ||
			actionCtx.action.Status == models.ActionStatusRunning {
			active = append(active, actionCtx.action)
		}
	}

	return active
}

// StopAllActions stops all currently running actions
func (e *Engine) StopAllActions() int {
	e.mu.Lock()
	defer e.mu.Unlock()

	count := 0
	for _, actionCtx := range e.actions {
		if actionCtx.action.Status == models.ActionStatusStarting ||
			actionCtx.action.Status == models.ActionStatusRunning {
			actionCtx.cancel()
			count++
		}
	}

	return count
}

// Cleanup removes completed actions
func (e *Engine) Cleanup() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for id, actionCtx := range e.actions {
		if actionCtx.action.Status == models.ActionStatusCompleted ||
			actionCtx.action.Status == models.ActionStatusFailed ||
			actionCtx.action.Status == models.ActionStatusStopped {
			// Only cleanup actions older than 1 minute
			if actionCtx.action.CompletedAt != nil &&
				time.Since(*actionCtx.action.CompletedAt) > time.Minute {
				delete(e.actions, id)
			}
		}
	}
}
