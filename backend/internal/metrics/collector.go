package metrics

import (
	"sync"
	"time"

	"monitoring-dashboard/pkg/models"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// Collector collects system metrics
type Collector struct {
	mu             sync.RWMutex
	currentMetrics models.Metrics
	prevDiskIO     disk.IOCountersStat
	prevNetIO      []net.IOCountersStat
	prevTime       time.Time
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	c := &Collector{
		prevTime: time.Now(),
	}
	// Initialize baseline metrics
	c.initializeBaseline()
	return c
}

// initializeBaseline initializes baseline measurements for rate calculations
func (c *Collector) initializeBaseline() {
	// Get initial disk I/O stats
	if diskStats, err := disk.IOCounters(); err == nil {
		for _, stat := range diskStats {
			c.prevDiskIO = stat
			break // Use first disk
		}
	}

	// Get initial network stats
	if netStats, err := net.IOCounters(false); err == nil {
		c.prevNetIO = netStats
	}
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

	// Collect Disk I/O (operations per second)
	metrics.DiskIO = c.collectDiskIO()

	// Collect Network (MB/s)
	metrics.Network = c.collectNetwork()

	c.mu.Lock()
	c.currentMetrics = metrics
	c.mu.Unlock()
}

// collectDiskIO calculates disk I/O operations per second
func (c *Collector) collectDiskIO() float64 {
	diskStats, err := disk.IOCounters()
	if err != nil || len(diskStats) == 0 {
		return 0.0
	}

	// Get current disk stats (use first disk)
	var currentDisk disk.IOCountersStat
	for _, stat := range diskStats {
		currentDisk = stat
		break
	}

	// Calculate time delta
	now := time.Now()
	timeDelta := now.Sub(c.prevTime).Seconds()
	if timeDelta <= 0 {
		return 0.0
	}

	// Calculate operations per second (read + write operations)
	readDelta := currentDisk.ReadCount - c.prevDiskIO.ReadCount
	writeDelta := currentDisk.WriteCount - c.prevDiskIO.WriteCount
	opsPerSec := float64(readDelta+writeDelta) / timeDelta

	// Update previous values
	c.prevDiskIO = currentDisk
	c.prevTime = now

	return opsPerSec
}

// collectNetwork calculates network throughput in MB/s
func (c *Collector) collectNetwork() float64 {
	netStats, err := net.IOCounters(false)
	if err != nil || len(netStats) == 0 {
		return 0.0
	}

	currentNet := netStats[0]

	// Calculate time delta
	now := time.Now()
	timeDelta := now.Sub(c.prevTime).Seconds()
	if timeDelta <= 0 || len(c.prevNetIO) == 0 {
		c.prevNetIO = netStats
		return 0.0
	}

	prevNet := c.prevNetIO[0]

	// Calculate bytes per second (sent + received)
	bytesSentDelta := currentNet.BytesSent - prevNet.BytesSent
	bytesRecvDelta := currentNet.BytesRecv - prevNet.BytesRecv
	bytesPerSec := float64(bytesSentDelta+bytesRecvDelta) / timeDelta

	// Convert to MB/s
	mbPerSec := bytesPerSec / (1024 * 1024)

	// Update previous values
	c.prevNetIO = netStats

	return mbPerSec
}

// GetCurrent returns the current metrics
func (c *Collector) GetCurrent() models.Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentMetrics
}
