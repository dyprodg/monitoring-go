package models

import "time"

// ActionType represents the type of load action
type ActionType string

const (
	ActionTypeCPUStress    ActionType = "cpu-stress"
	ActionTypeMemorySurge  ActionType = "memory-surge"
	ActionTypeDiskStorm    ActionType = "disk-storm"
	ActionTypeTrafficFlood ActionType = "traffic-flood"
)

// ActionStatus represents the current status of an action
type ActionStatus string

const (
	ActionStatusStarting  ActionStatus = "starting"
	ActionStatusRunning   ActionStatus = "running"
	ActionStatusCompleted ActionStatus = "completed"
	ActionStatusFailed    ActionStatus = "failed"
	ActionStatusStopped   ActionStatus = "stopped"
)

// Action represents a load generation action
type Action struct {
	ID          string       `json:"id"`
	Type        ActionType   `json:"type"`
	Status      ActionStatus `json:"status"`
	StartedAt   time.Time    `json:"started_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	Progress    float64      `json:"progress"` // 0.0 to 1.0
	Error       string       `json:"error,omitempty"`
}

// CPUStressRequest represents a request to start CPU stress
type CPUStressRequest struct {
	TargetPercent   int `json:"target_percent"`   // Target CPU percentage (0-95)
	DurationSeconds int `json:"duration_seconds"` // Duration in seconds (max 30)
}

// ActionResponse is the response after starting an action
type ActionResponse struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	StartedAt time.Time  `json:"started_at"`
	Message   string     `json:"message,omitempty"`
}
