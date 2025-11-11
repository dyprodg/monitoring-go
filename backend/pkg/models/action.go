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

// MemorySurgeRequest represents a request to start memory surge
type MemorySurgeRequest struct {
	SizeMB          int `json:"size_mb"`          // Memory size in MB (max 25% of total RAM)
	DurationSeconds int `json:"duration_seconds"` // Duration in seconds (max 60)
}

// DiskStormRequest represents a request to start disk storm
type DiskStormRequest struct {
	Operations int `json:"operations"`   // Number of file operations (max 10000)
	FileSizeKB int `json:"file_size_kb"` // File size in KB (max 1024)
}

// TrafficFloodRequest represents a request to start traffic flood
type TrafficFloodRequest struct {
	RequestsPerSec  int    `json:"requests_per_sec"`  // Requests per second (max 1000)
	DurationSeconds int    `json:"duration_seconds"`  // Duration in seconds (max 60)
	TargetURL       string `json:"target_url"`        // Target URL (optional, defaults to dummy endpoint)
}

// ActionResponse is the response after starting an action
type ActionResponse struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	StartedAt time.Time  `json:"started_at"`
	Message   string     `json:"message,omitempty"`
}
