package models

import "time"

// Metrics represents system metrics at a point in time
type Metrics struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float64   `json:"cpu"`        // Total CPU percentage (0-100)
	Memory    float64   `json:"memory"`     // Memory percentage (0-100)
	DiskIO    float64   `json:"disk_io"`    // Disk operations per second
	Network   float64   `json:"network"`    // Network MB/s
}

// HealthStatus represents the health of the service
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
