package metrics

import (
	"sync"
	"time"

	"monitoring-dashboard/pkg/models"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Collector collects system metrics
type Collector struct {
	mu            sync.RWMutex
	currentMetrics models.Metrics
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{}
}

// Start begins collecting metrics at regular intervals
func (c *Collector) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.collectMetrics()
		}
	}()

	// Collect initial metrics immediately
	c.collectMetrics()
}

// collectMetrics gathers current system metrics
func (c *Collector) collectMetrics() {
	metrics := models.Metrics{
		Timestamp: time.Now(),
	}

	// Collect CPU percentage
	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		metrics.CPU = cpuPercent[0]
	}

	// Collect Memory percentage
	if memStats, err := mem.VirtualMemory(); err == nil {
		metrics.Memory = memStats.UsedPercent
	}

	// Placeholder for disk I/O and network (will implement later)
	metrics.DiskIO = 0.0
	metrics.Network = 0.0

	c.mu.Lock()
	c.currentMetrics = metrics
	c.mu.Unlock()
}

// GetCurrent returns the current metrics
func (c *Collector) GetCurrent() models.Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentMetrics
}
