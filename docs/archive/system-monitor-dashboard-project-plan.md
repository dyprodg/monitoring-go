# Interactive System Monitoring & Load Testing Dashboard - Project Plan

## Project Overview
An **interactive** real-time system monitoring dashboard with Go backend, where users can actively trigger different load scenarios and observe the effects live. Perfect for demos, showcases, and performance testing.

## 🎯 Core Value Proposition
**"Watch your system react in real-time!"**
- User clicks button → System reacts instantly visually
- Various load scenarios (CPU Spikes, Memory Leaks, Traffic Surges)
- Live charts show immediate effects
- Perfect for demos, presentations and interactive showcases

---

## 🎯 Feature Overview (Priority: Must-Have → Nice-to-Have)

### ⭐ MUST-HAVE Features (MVP)

#### 1. **Interactive Control Panel** 🎮
**What:** User interface with buttons to trigger load actions
**Why:** The core feature - without this no interactive showcase
**Details:**
- Quick Action Buttons (4-6 items)
- One-click actions without configuration
- Visual feedback on click (button animation)
- Display of active actions

**Buttons:**
```
🔥 CPU Spike      - CPU to 90% for 10 seconds
💾 Memory Surge   - Allocate 500MB memory
💿 Disk Storm     - 1000 File I/O operations
🌐 Traffic Flood  - 1000 HTTP Requests/sec
```

#### 2. **Real-Time Metrics Dashboard** 📊
**What:** Live display of system metrics with instant updates
**Why:** User must see the effects of their actions IMMEDIATELY
**Details:**
- 4 Main Metrics: CPU, Memory, Network, Disk I/O
- Update Frequency: 1 second
- Circular Gauges for current values
- Line Charts with 60 seconds history
- Color coding: Green (OK) → Yellow (Warning) → Red (Critical)

**Metrics:**
- CPU: Percentage + per-core breakdown
- Memory: Used/Total in MB/GB + percentage
- Network: In/Out Traffic in MB/s + packet count
- Disk I/O: Read/Write in MB/s + operations/sec

#### 3. **Live Event Log** 📝
**What:** Chronological list of all events and actions
**Why:** Shows what's happening - storytelling element
**Details:**
- Timestamp for each event
- Color-coded event types (Info, Warning, Critical, Success)
- Auto-scroll to newest events
- Icons for different event types
- Max. 50 events visible (oldest are removed)

**Event Types:**
```
🔥 Action Started    - User triggered action
📊 Metric Alert      - Threshold exceeded
✅ Action Completed  - Action successfully finished
⚠️ Warning           - System under load
❌ Error            - Something went wrong
```

#### 4. **WebSocket Live-Updates** 🔌
**What:** Bidirectional real-time connection
**Why:** Without WebSocket no real live updates
**Details:**
- Server pushes metrics every second
- Client sends action commands to server
- Automatic reconnect on connection loss
- Connection status indicator

#### 5. **Load Action Execution Engine** ⚙️
**What:** Backend system that generates various loads
**Why:** The heart of the system - makes the magic happen
**Details:**
- Modular action implementations
- Concurrent execution (multiple actions in parallel)
- Graceful cancellation possible
- Resource cleanup after actions

**Implemented Actions:**
```go
- CPUStressAction      // Busy-loops in Goroutines
- MemoryLeakAction     // Allocates memory incrementally
- DiskIOAction         // Read/Write/Delete operations
- TrafficSimAction     // HTTP request flood to own server
- NetworkStressAction  // UDP/TCP traffic generation
```

---

### ⭐⭐ SHOULD-HAVE Features (Enhanced Experience)

#### 6. **Pre-built Scenarios** 🎬
**What:** Pre-configured action sequences with story
**Why:** Makes demos more interesting and realistic
**Details:**
- 3-5 predefined scenarios
- Timeline visualization during execution
- Pause/Resume/Stop functions
- Progress indicator

**Scenarios:**
```
🚀 "Startup Launch Day" (3min)
   0:00 - Normal operations
   0:30 - Press release → Traffic +200%
   1:00 - Traffic spike → +500%
   1:30 - Database stress
   2:00 - Auto-scaling kicks in
   2:30 - Recovery phase

🛒 "Black Friday Rush" (5min)
   0:00 - Pre-sale calm
   1:00 - Sale starts → Massive spike
   2:00 - Sustained high load
   3:00 - Payment processing peak
   4:00 - Gradual cooldown

💥 "DDoS Attack Simulation" (2min)
   0:00 - Normal operations
   0:20 - Attack begins
   0:40 - Peak attack intensity
   1:20 - Mitigation kicks in
   2:00 - Back to normal

🐛 "Memory Leak Detection" (5min)
   Gradual memory increase
   Simulates slow memory leak
   Shows GC struggling

⚡ "Morning Rush Hour" (3min)
   Simulates office workers
   logging in at 9am
```

#### 7. **Visual Action Feedback** ✨
**What:** Animations and visual effects during actions
**Why:** Makes it visually impressive - "Wow-Factor"
**Details:**
- Button pulse animation on activation
- Glow effects on affected metrics
- Color transitions (smooth not jumpy)
- Loading bars for action progress
- Alert badges/notifications

**Effects:**
- Chart lines pulsate during high load
- Gauge needles animated (not instant jump)
- Alert banners slide in/out
- Toast notifications for events

#### 8. **Active Actions Widget** 🎯
**What:** Overview of currently running actions
**Why:** User must know what's active
**Details:**
- List of active actions with progress
- Remaining time for each action
- Individual stop buttons
- "Stop All" button

```
Active Actions:
┌────────────────────────────┐
│ 🔥 CPU Stress              │
│ Progress: ▓▓▓▓▓▓▓░░░ 7/10s│
│ [Stop]                     │
├────────────────────────────┤
│ 🌐 Traffic Flood           │
│ Progress: ▓▓▓▓░░░░░░ 4/30s│
│ [Stop]                     │
└────────────────────────────┘
[⏹ Stop All Actions]
```

#### 9. **System Health Indicators** 💚
**What:** Overall system status
**Why:** Quick-glance status without reading charts
**Details:**
- Overall Health Score (0-100)
- Status Badges (Healthy/Warning/Critical)
- Key Metrics Summary
- Uptime Counter

```
┌─────────────────────┐
│ System Health       │
│ ● Healthy (Score: 92)│
│                     │
│ Uptime: 2h 34m     │
│ Requests: 45.2K    │
│ Errors: 0.02%      │
│ Avg Response: 42ms │
└─────────────────────┘
```

---

### ⭐⭐⭐ NICE-TO-HAVE Features (Polish & Advanced)

#### 10. **Custom Scenario Builder** 🎨
**What:** User can create their own action sequences
**Why:** For power users and custom testing
**Details:**
- Drag-and-drop timeline editor
- Action library to choose from
- Timing and duration configurable
- Save/Load custom scenarios
- Share link generation

```
Timeline Builder:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
0s    30s   60s   90s   120s

[+ Add Action]

Actions:
┌─────────────────┐
│ 🔥 CPU Stress   │ Drag to timeline
│ 💾 Memory Leak  │
│ 🌐 Traffic      │
└─────────────────┘
```

#### 11. **Historical Data & Replay** 📼
**What:** Save and replay past sessions
**Why:** Comparisons and analysis
**Details:**
- Session recording (metrics + events)
- Replay function with playback controls
- Export as JSON/CSV
- Compare mode (two sessions side-by-side)

#### 12. **Alert Configuration** 🔔
**What:** User-definable thresholds and alerts
**Why:** Custom monitoring rules
**Details:**
- Threshold settings per metric
- Alert actions (Log, Notification, Auto-Action)
- Alert history

```
Alert Rules:
┌─────────────────────────────┐
│ CPU > 80% for 10s → Alert   │
│ Memory > 90% → Critical     │
│ Response Time > 500ms → Warn│
└─────────────────────────────┘
```

#### 13. **Multi-Instance Monitoring** 🖥️🖥️
**What:** Monitor multiple servers/containers simultaneously
**Why:** Realistic load-balancing scenarios
**Details:**
- Agent deployment on multiple nodes
- Combined dashboard view
- Per-instance drill-down
- Load distribution visualization

#### 14. **Performance Comparison Mode** ⚖️
**What:** Before/After or side-by-side comparisons
**Why:** Shows impact of optimizations
**Details:**
- Split-screen view
- Snapshot function
- Diff highlighting
- Export comparison report

```
┌─────────────┬─────────────┐
│   Before    │    After    │
├─────────────┼─────────────┤
│ CPU: 45%    │ CPU: 78%    │
│ Response:   │ Response:   │
│   42ms      │   156ms     │
└─────────────┴─────────────┘
```

#### 15. **Gamification Elements** 🎮
**What:** Challenges and scoring system
**Why:** Makes it entertaining and engaging
**Details:**
- Challenge mode with levels
- Scoring system based on system stability
- Leaderboard (optional)
- Achievements/badges

```
Challenge: "Survive the Storm"
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Goal: Keep CPU < 95% for 60s
      while under heavy load

Current: 45s elapsed
CPU Peak: 92% ✅
Score: ⭐⭐⭐⭐☆ (4/5)
```

#### 16. **Export & Reporting** 📄
**What:** Generate reports from sessions
**Why:** Documentation and sharing
**Details:**
- PDF report generation
- Chart screenshots
- CSV export for metrics
- Markdown summary

#### 17. **Dark/Light Theme Toggle** 🌓
**What:** Theme switcher
**Why:** User preference
**Details:**
- Smooth theme transition
- Persisted user preference
- High-contrast mode (accessibility)

#### 18. **Sound Effects** 🔊
**What:** Audio feedback for events
**Why:** Extra immersion (optional)
**Details:**
- Alert sounds for critical events
- Success jingles
- Mute toggle
- Volume control

---

## 📋 Feature Priority Matrix

```
┌────────────────────────────────────────────────────────────┐
│ High Impact + Low Effort → IMPLEMENT FIRST                 │
├────────────────────────────────────────────────────────────┤
│ ✅ Interactive Control Panel (Must-Have #1)                │
│ ✅ Real-Time Metrics Dashboard (Must-Have #2)              │
│ ✅ Live Event Log (Must-Have #3)                           │
│ ✅ WebSocket Live-Updates (Must-Have #4)                   │
│ ✅ Load Action Engine (Must-Have #5)                       │
├────────────────────────────────────────────────────────────┤
│ High Impact + Medium Effort → IMPLEMENT SECOND             │
├────────────────────────────────────────────────────────────┤
│ ⭐ Pre-built Scenarios (Should-Have #6)                    │
│ ⭐ Visual Action Feedback (Should-Have #7)                 │
│ ⭐ Active Actions Widget (Should-Have #8)                  │
├────────────────────────────────────────────────────────────┤
│ Medium Impact + Low Effort → POLISH PHASE                  │
├────────────────────────────────────────────────────────────┤
│ ⭐ System Health Indicators (Should-Have #9)               │
│ ⭐⭐⭐ Dark/Light Theme (Nice-to-Have #17)                  │
│ ⭐⭐⭐ Sound Effects (Nice-to-Have #18)                     │
├────────────────────────────────────────────────────────────┤
│ High Impact + High Effort → FUTURE ENHANCEMENTS            │
├────────────────────────────────────────────────────────────┤
│ ⭐⭐⭐ Custom Scenario Builder (Nice-to-Have #10)           │
│ ⭐⭐⭐ Historical Data & Replay (Nice-to-Have #11)          │
│ ⭐⭐⭐ Multi-Instance Monitoring (Nice-to-Have #13)         │
├────────────────────────────────────────────────────────────┤
│ Low Priority / Optional                                    │
├────────────────────────────────────────────────────────────┤
│ ⭐⭐⭐ Alert Configuration (Nice-to-Have #12)               │
│ ⭐⭐⭐ Performance Comparison (Nice-to-Have #14)            │
│ ⭐⭐⭐ Gamification (Nice-to-Have #15)                      │
│ ⭐⭐⭐ Export & Reporting (Nice-to-Have #16)                │
└────────────────────────────────────────────────────────────┘
```

---

## 🎯 MVP Scope (Minimum Viable Product)

**What will be implemented in MVP:**
- ✅ Features #1-5 (all Must-Haves)
- ✅ Feature #6 (Pre-built Scenarios) - 2-3 Scenarios
- ✅ Feature #7 (Visual Feedback) - Basic animations
- ✅ Feature #8 (Active Actions Widget)

**What comes later:**
- Feature #9 and higher as needed

**Goal:** Functional, impressive demo in 4-6 weeks

---

## Phase 1: Project Setup & Basic Structure

### 1.1 Repository Initialization
```bash
# Project structure
system-monitor-dashboard/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── metrics/
│   │   │   ├── collector.go      # Real system metrics
│   │   │   ├── cpu.go
│   │   │   ├── memory.go
│   │   │   ├── disk.go
│   │   │   └── network.go
│   │   ├── loadgen/               # ⭐ NEW: Load generation
│   │   │   ├── engine.go          # Action execution engine
│   │   │   ├── actions.go         # Action interface
│   │   │   ├── cpu_stress.go      # CPU load generator
│   │   │   ├── memory_leak.go     # Memory load generator
│   │   │   ├── disk_io.go         # Disk I/O generator
│   │   │   ├── traffic_sim.go     # HTTP traffic simulator
│   │   │   └── scenarios.go       # Pre-built scenarios
│   │   ├── websocket/
│   │   │   ├── hub.go
│   │   │   ├── client.go
│   │   │   └── handler.go
│   │   ├── api/
│   │   │   ├── handlers.go
│   │   │   ├── routes.go
│   │   │   ├── actions_handler.go  # ⭐ NEW: Action endpoints
│   │   │   └── scenarios_handler.go # ⭐ NEW: Scenario endpoints
│   │   └── config/
│   │       └── config.go
│   ├── pkg/
│   │   └── models/
│   │       ├── metrics.go
│   │       ├── action.go          # ⭐ NEW: Action models
│   │       └── event.go           # ⭐ NEW: Event models
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Dashboard/
│   │   │   │   ├── CPUChart.jsx
│   │   │   │   ├── MemoryChart.jsx
│   │   │   │   ├── DiskChart.jsx
│   │   │   │   └── NetworkChart.jsx
│   │   │   ├── ControlPanel/      # ⭐ NEW: Control components
│   │   │   │   ├── ControlPanel.jsx
│   │   │   │   ├── QuickActions.jsx
│   │   │   │   ├── ScenarioSelector.jsx
│   │   │   │   └── ActiveActions.jsx
│   │   │   ├── EventLog/          # ⭐ NEW: Event log
│   │   │   │   └── EventLog.jsx
│   │   │   └── Shared/
│   │   │       ├── Gauge.jsx
│   │   │       ├── LineChart.jsx
│   │   │       └── ActionButton.jsx
│   │   ├── services/
│   │   │   ├── websocket.js
│   │   │   └── api.js            # ⭐ NEW: API client
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── Dockerfile
├── infrastructure/
│   ├── terraform/
│   │   ├── main.tf
│   │   ├── vpc.tf
│   │   ├── ecs.tf
│   │   ├── alb.tf
│   │   ├── rds.tf
│   │   └── variables.tf
│   └── k8s/ (optional alternative)
├── .github/
│   └── workflows/
│       ├── backend-ci.yml
│       ├── frontend-ci.yml
│       └── deploy.yml
├── docker-compose.yml
└── README.md
```

### 1.2 Technology Stack
**Backend:**
- Go 1.21+
- Gorilla WebSocket
- gopsutil (system metrics)
- Chi Router (HTTP Router)
- Prometheus Client (optional for metrics export)
- sync/atomic (for concurrent action management)
- context (for cancellable actions)

**Frontend:**
- React 18+ with Vite
- Recharts (chart library)
- TailwindCSS
- WebSocket API

**Infrastructure:**
- AWS ECS Fargate
- AWS Application Load Balancer
- AWS RDS PostgreSQL (for historical data)
- AWS CloudWatch
- AWS ECR (Container Registry)
- Terraform for IaC

**CI/CD:**
- GitHub Actions
- Docker
- AWS CodeDeploy (optional)

---

## Phase 2: Backend Development

### 2.1 Core Metrics Collector
**File: `internal/metrics/collector.go`**
```go
// Tasks:
// - Collect system metrics (CPU, Memory, Disk, Network)
// - Sampling rate: every 1-2 seconds
// - Goroutine-based concurrent collection
// - Error handling and logging
```

**Metrics to implement:**
- CPU: Usage %, per-core usage, load average
- Memory: Total, Used, Free, Swap
- Disk: Usage per partition, I/O stats
- Network: Bytes sent/received, packets, errors

### 2.2 Load Generation Engine ⭐ NEW
**File: `internal/loadgen/engine.go`**
```go
// Core engine for action execution
type ActionEngine struct {
    activeActions map[string]*RunningAction
    mu            sync.RWMutex
    eventChan     chan Event
}

// Action management
func (e *ActionEngine) StartAction(action Action) (string, error)
func (e *ActionEngine) StopAction(actionID string) error
func (e *ActionEngine) GetActiveActions() []ActionStatus
```

**File: `internal/loadgen/cpu_stress.go`**
```go
// CPU load generator
type CPUStressAction struct {
    TargetPercent int           // Target CPU %
    Duration      time.Duration
    Cores         int            // Number of cores to stress
}

func (a *CPUStressAction) Execute(ctx context.Context) error {
    // Busy loops in Goroutines
    // Math operations for CPU load
    // Context-aware for graceful shutdown
}
```

**File: `internal/loadgen/memory_leak.go`**
```go
// Memory load generator
type MemoryLeakAction struct {
    SizeMB        int
    Duration      time.Duration
    GradualSteps  int  // Number of steps
}

func (a *MemoryLeakAction) Execute(ctx context.Context) error {
    // Allocate memory in steps
    // Hold references to prevent GC
    // Cleanup after duration
}
```

**File: `internal/loadgen/disk_io.go`**
```go
// Disk I/O generator
type DiskIOAction struct {
    Operations int
    FileSizeKB int
    Pattern    string  // "sequential", "random"
}

func (a *DiskIOAction) Execute(ctx context.Context) error {
    // Create temp files
    // Read/Write operations
    // Cleanup afterwards
}
```

**File: `internal/loadgen/traffic_sim.go`**
```go
// HTTP traffic simulator
type TrafficSimAction struct {
    RequestsPerSec int
    Duration       time.Duration
    Endpoints      []string
}

func (a *TrafficSimAction) Execute(ctx context.Context) error {
    // HTTP request flood
    // Rate-limited to RequestsPerSec
    // Concurrent requests with worker pool
}
```

**File: `internal/loadgen/scenarios.go`**
```go
// Pre-built scenarios
type Scenario struct {
    ID          string
    Name        string
    Description string
    Steps       []ScenarioStep
}

type ScenarioStep struct {
    Delay  time.Duration
    Action Action
}

// Predefined scenarios
var PrebuiltScenarios = map[string]Scenario{
    "startup-launch": {
        Name: "Startup Launch Day",
        Steps: []ScenarioStep{
            {Delay: 0, Action: NormalOps()},
            {Delay: 30*time.Second, Action: TrafficSurge(200)},
            {Delay: 60*time.Second, Action: TrafficSpike(500)},
            // ...
        },
    },
    "black-friday": {...},
    "ddos-simulation": {...},
}
```

### 2.3 WebSocket Server
**File: `internal/websocket/hub.go`**
```go
// Tasks:
// - Manage client connections
// - Broadcast metrics to all clients
// - Broadcast events (action started/completed)
// - Connection pooling
// - Graceful disconnect handling
```

### 2.4 REST API ⭐ EXTENDED
**Endpoints:**

**Metrics:**
- `GET /api/health` - Health check
- `GET /api/metrics/current` - Current metrics
- `GET /api/metrics/history?duration=1h` - Historical data
- `WS /ws` - WebSocket connection

**Load Actions:** ⭐ NEW
- `POST /api/actions/cpu-stress` - Start CPU load
- `POST /api/actions/memory-leak` - Start memory load
- `POST /api/actions/disk-io` - Start disk I/O
- `POST /api/actions/traffic` - Start traffic load
- `GET /api/actions/active` - List active actions
- `DELETE /api/actions/:id/stop` - Stop action
- `DELETE /api/actions/stop-all` - Stop all actions

**Scenarios:** ⭐ NEW
- `GET /api/scenarios` - List available scenarios
- `GET /api/scenarios/:id` - Scenario details
- `POST /api/scenarios/:id/start` - Start scenario
- `GET /api/scenarios/running` - Running scenarios
- `DELETE /api/scenarios/:id/stop` - Stop scenario

**Request/Response Examples:**
```json
// POST /api/actions/cpu-stress
{
  "target_percent": 90,
  "duration_seconds": 10,
  "cores": 0  // 0 = all cores
}

// Response
{
  "action_id": "cpu-stress-abc123",
  "status": "started",
  "started_at": "2024-01-15T10:30:00Z"
}

// GET /api/actions/active
{
  "actions": [
    {
      "id": "cpu-stress-abc123",
      "type": "cpu-stress",
      "status": "running",
      "progress": 0.7,
      "started_at": "2024-01-15T10:30:00Z",
      "estimated_completion": "2024-01-15T10:30:10Z"
    }
  ]
}
```

### 2.5 Event System ⭐ NEW
**File: `pkg/models/event.go`**
```go
type Event struct {
    ID        string
    Type      EventType
    Message   string
    Severity  Severity
    Timestamp time.Time
    Metadata  map[string]interface{}
}

type EventType string
const (
    ActionStarted   EventType = "action_started"
    ActionCompleted EventType = "action_completed"
    MetricAlert     EventType = "metric_alert"
    SystemWarning   EventType = "system_warning"
)

type Severity string
const (
    Info     Severity = "info"
    Warning  Severity = "warning"
    Critical Severity = "critical"
    Success  Severity = "success"
)
```

### 2.6 Testing (Backend)
```bash
# Unit tests
go test ./... -v -cover

# Integration tests
go test ./tests/integration/... -v

# Load action tests
go test ./internal/loadgen/... -v

# Benchmark tests
go test -bench=. ./internal/metrics/
go test -bench=. ./internal/loadgen/
```

**Test coverage goal: >80%**

**Important test cases:**
- CPU Stress: Verify CPU increase
- Memory Leak: Verify memory increase
- Concurrent Actions: Multiple actions simultaneously
- Action Cancellation: Graceful shutdown
- Event Broadcasting: Events reach all clients

---

## Phase 3: Frontend Development

### 3.1 WebSocket Service
**File: `src/services/websocket.js`**
```javascript
// Tasks:
// - Establish WebSocket connection
// - Reconnect logic
// - Event handlers for incoming data (metrics + events)
// - Command sender for actions
// - Error handling
```

### 3.2 API Client ⭐ NEW
**File: `src/services/api.js`**
```javascript
// REST API client
class APIClient {
  // Actions
  async startCPUStress(params) { ... }
  async startMemoryLeak(params) { ... }
  async startDiskIO(params) { ... }
  async startTraffic(params) { ... }
  async stopAction(actionId) { ... }
  async getActiveActions() { ... }

  // Scenarios
  async getScenarios() { ... }
  async startScenario(scenarioId) { ... }
  async stopScenario(scenarioId) { ... }
}
```

### 3.3 Interactive Control Panel ⭐ NEW
**File: `src/components/ControlPanel/ControlPanel.jsx`**
```jsx
// Main control panel
function ControlPanel() {
  const [activeActions, setActiveActions] = useState([]);
  const [scenarios, setScenarios] = useState([]);

  return (
    <div className="control-panel">
      <h2>🎮 Load Control Center</h2>
      <QuickActions onActionTriggered={handleActionStart} />
      <ScenarioSelector scenarios={scenarios} />
      <ActiveActions actions={activeActions} />
    </div>
  );
}
```

**File: `src/components/ControlPanel/QuickActions.jsx`**
```jsx
// Quick action buttons
function QuickActions({ onActionTriggered }) {
  const actions = [
    {
      id: 'cpu-stress',
      icon: '🔥',
      title: 'CPU Spike',
      description: '90% for 10s',
      color: 'red',
      params: { target_percent: 90, duration_seconds: 10 }
    },
    {
      id: 'memory-leak',
      icon: '💾',
      title: 'Memory Surge',
      description: '500MB',
      color: 'orange',
      params: { size_mb: 500, duration_seconds: 60 }
    },
    // ...
  ];

  return (
    <div className="quick-actions-grid">
      {actions.map(action => (
        <ActionButton
          key={action.id}
          {...action}
          onClick={() => triggerAction(action)}
        />
      ))}
    </div>
  );
}
```

**File: `src/components/ControlPanel/ScenarioSelector.jsx`**
```jsx
// Scenario selection & control
function ScenarioSelector({ scenarios }) {
  const [running, setRunning] = useState(null);

  return (
    <div className="scenario-selector">
      <h3>📋 Pre-built Scenarios</h3>
      {scenarios.map(scenario => (
        <ScenarioCard
          key={scenario.id}
          scenario={scenario}
          isRunning={running?.id === scenario.id}
          onStart={() => startScenario(scenario)}
          onStop={() => stopScenario(scenario)}
        />
      ))}
    </div>
  );
}
```

**File: `src/components/ControlPanel/ActiveActions.jsx`**
```jsx
// Active actions widget
function ActiveActions({ actions }) {
  return (
    <div className="active-actions">
      <h3>⚡ Running Actions</h3>
      {actions.length === 0 ? (
        <p className="no-actions">No active actions</p>
      ) : (
        <>
          {actions.map(action => (
            <ActionProgressCard
              key={action.id}
              action={action}
              onStop={() => stopAction(action.id)}
            />
          ))}
          <button
            className="stop-all-btn"
            onClick={stopAllActions}
          >
            ⏹ Stop All Actions
          </button>
        </>
      )}
    </div>
  );
}
```

### 3.4 Chart Components ⭐ EXTENDED
**Components:**
- Real-time line charts for CPU/Memory/Network
- Gauge charts for current values (with animations!)
- Donut charts for disk usage
- Alert badges when thresholds exceeded
- Glow effects during high load
- Smooth color transitions

**File: `src/components/Shared/Gauge.jsx`**
```jsx
// Animated circular gauge
function Gauge({ value, max, label, status }) {
  // Smooth needle animation
  // Color coding: green → yellow → red
  // Glow effect when critical
  // Pulsing animation during high load
}
```

**File: `src/components/Shared/LineChart.jsx`**
```jsx
// Real-time line chart with Recharts
function LineChart({ data, metric }) {
  // Auto-scaling Y-axis
  // 60 seconds sliding window
  // Annotation markers for events
  // Smooth line interpolation
}
```

### 3.5 Event Log Component ⭐ NEW
**File: `src/components/EventLog/EventLog.jsx`**
```jsx
// Live event log
function EventLog({ events }) {
  const eventLogRef = useRef(null);

  // Auto-scroll to newest events
  useEffect(() => {
    if (eventLogRef.current) {
      eventLogRef.current.scrollTop = eventLogRef.current.scrollHeight;
    }
  }, [events]);

  return (
    <div className="event-log" ref={eventLogRef}>
      <h3>🎬 Event Log</h3>
      {events.map(event => (
        <EventItem
          key={event.id}
          event={event}
          icon={getEventIcon(event.type)}
          color={getSeverityColor(event.severity)}
        />
      ))}
    </div>
  );
}
```

### 3.6 Dashboard Layout ⭐ EXTENDED
**File: `src/App.jsx`**
```jsx
function App() {
  return (
    <div className="dashboard">
      {/* Header */}
      <Header>
        <h1>📊 System Monitor Dashboard</h1>
        <StatusIndicator />
        <ThemeToggle />
      </Header>

      {/* Main Content */}
      <div className="dashboard-grid">
        {/* Left: Control Panel */}
        <aside className="control-sidebar">
          <ControlPanel />
        </aside>

        {/* Center: Metrics */}
        <main className="metrics-main">
          <div className="metrics-grid">
            <CPUChart />
            <MemoryChart />
            <NetworkChart />
            <DiskChart />
          </div>
        </main>

        {/* Right: Events & Info */}
        <aside className="info-sidebar">
          <SystemHealth />
          <EventLog />
        </aside>
      </div>
    </div>
  );
}
```

**Layout:**
- Responsive grid layout
- Dark/Light theme toggle
- 3-column layout (Control | Metrics | Events)
- Mobile: Stacked layout

### 3.7 Visual Effects & Animations ⭐ NEW
```css
/* Button pulse animation */
@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

/* Glow effect */
.metric-critical {
  box-shadow: 0 0 20px rgba(255, 0, 0, 0.5);
  animation: glow 2s infinite;
}

/* Smooth transitions */
.gauge-needle {
  transition: transform 0.5s ease-out;
}

/* Alert slide in */
@keyframes slideIn {
  from { transform: translateY(-100%); }
  to { transform: translateY(0); }
}
```

### 3.8 Testing (Frontend)
```bash
# Unit tests (Vitest)
npm run test

# Component tests
npm run test:component

# E2E tests (Playwright)
npm run test:e2e

# Visual regression tests
npm run test:visual
```

**Test cases:**
- Button clicks trigger API calls
- WebSocket updates render correctly
- Action progress updates
- Event log scrolling
- Chart animations
- Responsive layout

---

## Phase 4: Containerization

### 4.1 Backend Dockerfile
```dockerfile
# Multi-stage build:
# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### 4.2 Frontend Dockerfile
```dockerfile
# Multi-stage build:
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Serve
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 4.3 docker-compose.yml (for local development)
```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - ENV=development

  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: monitoring
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

---

## Phase 5: CI/CD Pipeline

### 5.1 Backend CI Pipeline
**File: `.github/workflows/backend-ci.yml`**
```yaml
name: Backend CI

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'backend/**'
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        working-directory: ./backend
        run: |
          go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.txt

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          working-directory: ./backend

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-central-1

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      - name: Build and push Docker image
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: monitor-backend
          IMAGE_TAG: ${{ github.sha }}
        run: |
          cd backend
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG .
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
          docker tag $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG $ECR_REGISTRY/$ECR_REPOSITORY:latest
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:latest
```

### 5.2 Frontend CI Pipeline
**File: `.github/workflows/frontend-ci.yml`**
```yaml
name: Frontend CI

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'frontend/**'
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci

      - name: Run tests
        working-directory: ./frontend
        run: npm test

      - name: Run linter
        working-directory: ./frontend
        run: npm run lint

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-central-1

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      - name: Build and push Docker image
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: monitor-frontend
          IMAGE_TAG: ${{ github.sha }}
        run: |
          cd frontend
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG .
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
          docker tag $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG $ECR_REGISTRY/$ECR_REPOSITORY:latest
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:latest
```

### 5.3 Deployment Pipeline
**File: `.github/workflows/deploy.yml`**
```yaml
name: Deploy to AWS

on:
  workflow_run:
    workflows: ["Backend CI", "Frontend CI"]
    types:
      - completed
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    if: ${{ github.event.workflow_run.conclusion == 'success' }}
    steps:
      - uses: actions/checkout@v3

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: eu-central-1

      - name: Update ECS service
        run: |
          aws ecs update-service \
            --cluster monitor-cluster \
            --service monitor-backend-service \
            --force-new-deployment

          aws ecs update-service \
            --cluster monitor-cluster \
            --service monitor-frontend-service \
            --force-new-deployment

      - name: Wait for deployment
        run: |
          aws ecs wait services-stable \
            --cluster monitor-cluster \
            --services monitor-backend-service monitor-frontend-service
```

---

## Phase 6: AWS Infrastructure (Terraform)

### 6.1 VPC & Networking
**File: `infrastructure/terraform/vpc.tf`**
```hcl
# VPC with Public/Private Subnets
# NAT Gateway
# Internet Gateway
# Security Groups
```

### 6.2 ECS Cluster
**File: `infrastructure/terraform/ecs.tf`**
```hcl
resource "aws_ecs_cluster" "main" {
  name = "monitor-cluster"
}

# Task Definitions for Backend & Frontend
# ECS Services with Auto-Scaling
# CloudWatch Log Groups
```

### 6.3 Load Balancer
**File: `infrastructure/terraform/alb.tf`**
```hcl
# Application Load Balancer
# Target Groups
# Listener Rules
# SSL/TLS Certificate (ACM)
```

### 6.4 RDS PostgreSQL
**File: `infrastructure/terraform/rds.tf`**
```hcl
# RDS Instance for historical data
# Automated Backups
# Multi-AZ for Production
```

### 6.5 Terraform Commands
```bash
cd infrastructure/terraform

# Initialize
terraform init

# Plan
terraform plan -out=tfplan

# Apply
terraform apply tfplan

# Destroy (when needed)
terraform destroy
```

---

## Phase 7: Monitoring & Observability

### 7.1 CloudWatch Integration
- Container logs to CloudWatch
- Custom metrics
- Alarms for:
  - High CPU/Memory
  - Service Health
  - Response Time
  - Error Rate

### 7.2 Health Checks
```go
// Backend Health Check Endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
    health := map[string]string{
        "status": "healthy",
        "version": "1.0.0",
        "timestamp": time.Now().Format(time.RFC3339),
    }
    json.NewEncoder(w).Encode(health)
}
```

### 7.3 Metrics Export
- Prometheus-compatible metrics
- Grafana dashboard (optional)

---

## Phase 8: Security & Best Practices

### 8.1 Security Checklist
- [ ] HTTPS/TLS for all connections
- [ ] Environment variables for secrets
- [ ] AWS Secrets Manager for credentials
- [ ] Rate limiting on API
- [ ] CORS configuration
- [ ] Input validation
- [ ] SQL injection protection (Prepared Statements)
- [ ] Security headers (CSP, HSTS, etc.)

### 8.2 AWS IAM Policies
- Least privilege principle
- Separate roles for CI/CD
- ECS Task Execution Role
- RDS Access Policies

---

## Phase 9: Testing Strategy

### 9.1 Backend Tests
```bash
# Unit tests
cd backend
go test ./internal/... -v

# Integration tests
go test ./tests/integration/... -v

# Load tests (with vegeta or k6)
vegeta attack -duration=60s -rate=1000 \
  -targets=targets.txt | vegeta report
```

### 9.2 Frontend Tests
```bash
cd frontend

# Unit tests
npm run test

# Component tests
npm run test:component

# E2E tests
npm run test:e2e
```

### 9.3 Performance Tests
- Load testing with k6 or Artillery
- WebSocket connection tests
- Latency measurements
- Stress testing

---

## Phase 10: Documentation

### 10.1 README.md
- Project overview
- Features
- Setup instructions
- API documentation
- Deployment guide

### 10.2 API Documentation
- OpenAPI/Swagger Spec
- Request/Response examples
- WebSocket protocol

### 10.3 Architecture Decision Records (ADRs)
- Technology decisions
- Architecture patterns
- Trade-offs

---

## Implementation Order for Claude Code

### Sprint 1 (Week 1): Foundation
**Goal:** Basic Backend + Metrics Collection
1. ✅ Set up repository structure
2. ✅ Backend: Basic HTTP server with Chi Router
3. ✅ Backend: Implement metrics collector (CPU, Memory)
4. ✅ Backend: Unit tests for metrics
5. ✅ Docker setup (Backend)
6. ✅ Local docker-compose for development

**Deliverable:** Backend runs locally and collects metrics

---

### Sprint 2 (Week 2): Load Generation Engine ⭐
**Goal:** Implement action system
7. ✅ Backend: Action interface & engine (`loadgen/engine.go`)
8. ✅ Backend: Implement CPU stress action
9. ✅ Backend: Implement memory leak action
10. ✅ Backend: Implement disk I/O action
11. ✅ Backend: Implement traffic sim action
12. ✅ Backend: Action API endpoints
13. ✅ Backend: Implement event system
14. ✅ Backend: Unit tests for all actions

**Deliverable:** Backend can generate load via API

---

### Sprint 3 (Week 3): WebSocket & Frontend Foundation
**Goal:** Real-time updates + Basic UI
15. ✅ Backend: Implement WebSocket hub
16. ✅ Backend: Metrics broadcasting over WebSocket
17. ✅ Backend: Events broadcasting over WebSocket
18. ✅ Backend: Integration tests
19. ✅ Frontend: Set up React app (Vite + Tailwind)
20. ✅ Frontend: Implement WebSocket service
21. ✅ Frontend: Implement API client
22. ✅ Frontend: Basic dashboard layout
23. ✅ Docker setup (Frontend)

**Deliverable:** Frontend receives live data

---

### Sprint 4 (Week 4): Interactive Control Panel ⭐
**Goal:** User can trigger actions
24. ✅ Frontend: ControlPanel component
25. ✅ Frontend: QuickActions component with buttons
26. ✅ Frontend: ActionButton component (styled & animated)
27. ✅ Frontend: Active Actions widget
28. ✅ Frontend: Action triggering works
29. ✅ Backend: Complete Disk & Network metrics
30. ✅ Frontend: Toast notifications for events
31. ✅ Frontend: Button click → Backend → Metric change (E2E test)

**Deliverable:** User can click buttons and see effects live!

---

### Sprint 5 (Week 5): Charts & Visual Polish ⭐
**Goal:** Impressive visualizations
32. ✅ Frontend: Gauge components (animated)
33. ✅ Frontend: Line chart components (Recharts)
34. ✅ Frontend: CPU/Memory/Network/Disk charts
35. ✅ Frontend: Color-coded states (green/yellow/red)
36. ✅ Frontend: Glow effects for critical states
37. ✅ Frontend: Event log component
38. ✅ Frontend: System health widget
39. ✅ Frontend: Responsive design
40. ✅ Frontend: Dark theme styling

**Deliverable:** Visually impressive dashboard

---

### Sprint 6 (Week 6): Scenarios & Testing ⭐
**Goal:** Pre-built scenarios + Production-ready
41. ✅ Backend: Implement scenario system
42. ✅ Backend: 3 pre-built scenarios (Startup Launch, Black Friday, DDoS)
43. ✅ Backend: Scenario API endpoints
44. ✅ Frontend: Scenario selector component
45. ✅ Frontend: Scenario progress visualization
46. ✅ Frontend: Scenario timeline display
47. ✅ Backend: Historical data (PostgreSQL)
48. ✅ Frontend: Tests (Unit + E2E)
49. ✅ Backend: Tests (Coverage >80%)
50. ✅ Finalize documentation

**Deliverable:** Complete MVP with scenarios

---

### Sprint 7 (Week 7): CI/CD Pipeline
**Goal:** Automated deployment
51. ✅ GitHub Actions: Backend CI pipeline
52. ✅ GitHub Actions: Frontend CI pipeline
53. ✅ AWS ECR setup
54. ✅ Automated testing in pipeline
55. ✅ Code coverage reports
56. ✅ Linting & format checks
57. ✅ Docker image building & pushing

**Deliverable:** Automated CI/CD

---

### Sprint 8 (Week 8): AWS Infrastructure
**Goal:** Production deployment
58. ✅ Terraform: VPC & Networking
59. ✅ Terraform: ECS Cluster
60. ✅ Terraform: Load Balancer
61. ✅ Terraform: RDS PostgreSQL
62. ✅ Terraform: CloudWatch & Alarms
63. ✅ Infrastructure testing
64. ✅ GitHub Actions: Deployment pipeline
65. ✅ SSL/TLS setup
66. ✅ Domain configuration

**Deliverable:** Live on AWS!

---

### Sprint 9 (Week 9): Polish & Optimization
**Goal:** Production-quality
67. ✅ Performance optimization (Backend)
68. ✅ Performance optimization (Frontend)
69. ✅ Security hardening
70. ✅ Monitoring & logging setup
71. ✅ Load testing
72. ✅ Bug fixes
73. ✅ UX improvements
74. ✅ Documentation updates

**Deliverable:** Production-ready system

---

## 🎯 Critical Path (Must work for demo)

```
Week 1-2: Backend Foundation + Load Engine
         ↓
Week 3-4: Frontend + Interactive Controls
         ↓
Week 5:   Visual Polish (Charts & Animations)
         ↓
Week 6:   Scenarios + Testing
         ↓
DEMO READY! ✨
         ↓
Week 7-9: CI/CD + AWS (optional for local demo)
```

---

## 🚀 Quickstart for Claude Code

**First Session:**
```bash
# Sprint 1, Tasks 1-6
"Create the repository structure, set up a Go HTTP server
with Chi Router, implement CPU and Memory metrics collection
with gopsutil, write unit tests, and create a Dockerfile
plus docker-compose.yml for local development."
```

**Second Session:**
```bash
# Sprint 2, Tasks 7-14
"Implement the load generation system: Create an action
interface, an execution engine, and implement 4 action types
(CPU Stress, Memory Leak, Disk I/O, Traffic Sim). Add REST API
endpoints to start/stop actions. Implement an event system.
Write tests for all actions."
```

**Third Session:**
```bash
# Sprint 3, Tasks 15-23
"Implement WebSocket support in the backend for real-time updates.
Create a React frontend app with Vite and TailwindCSS.
Implement WebSocket and API client services. Create a
basic dashboard layout with grid. Docker setup for frontend."
```

**Fourth Session:**
```bash
# Sprint 4, Tasks 24-31
"Create the interactive control panel in the frontend: QuickActions
component with 4 action buttons (CPU, Memory, Disk, Traffic),
Active Actions widget with progress bars, toast notifications.
Connect frontend actions to backend API. E2E test that
button click leads to metric change."
```

And so on...

---

## Environment Variables

### Backend
```env
PORT=8080
ENV=production
DB_HOST=postgres.xxxxx.eu-central-1.rds.amazonaws.com
DB_PORT=5432
DB_NAME=monitoring
DB_USER=admin
DB_PASSWORD=<from-secrets-manager>
METRICS_INTERVAL=2s
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
```

### Frontend
```env
VITE_API_URL=https://api.monitor.example.com
VITE_WS_URL=wss://api.monitor.example.com/ws
```

---

## Cost Estimation (AWS)

**Monthly costs (approx.):**
- ECS Fargate (2 Tasks): ~$30-50
- Application Load Balancer: ~$20
- RDS t3.micro: ~$15-20
- Data Transfer: ~$5-10
- CloudWatch Logs: ~$5
- ECR Storage: ~$1

**Total: ~$75-110/month** (can be reduced with Reserved Instances)

---

## Success Metrics

### Performance
- API Response Time: <100ms (p95)
- WebSocket Latency: <50ms
- Frontend Load Time: <2s
- Action Execution: Immediate visual feedback (<500ms)
- Chart Update Rate: 1 FPS minimum

### Reliability
- Uptime: 99.5%+
- Test Coverage: >80%
- Zero Critical Security Vulnerabilities
- Graceful Action Cancellation: 100%
- WebSocket Reconnect: <5s

### User Experience ⭐
- Button Click → Visual Response: <200ms
- Action Started → Metric Change Visible: <2s
- Smooth Animations: 60 FPS
- Event Log Updates: Real-time (<1s)
- Mobile Responsive: All features work

### Monitoring
- All metrics collected every 1-2s
- Real-time updates in frontend
- Historical data retention: 30 days
- Event logging: All actions captured

---

## Next Steps

1. **Create repository** on GitHub
2. **Prepare AWS account** (IAM User, Access Keys)
3. **Set up local development environment:**
   ```bash
   # Install Go
   # Install Docker
   # Install Node.js
   # Install Terraform
   # Install AWS CLI
   ```
4. **Start Sprint 1** with Claude Code:
   - Repository structure
   - Backend foundation
   - Metrics collector

---

## Resources & Links

- [Go Documentation](https://go.dev/doc/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [gopsutil](https://github.com/shirou/gopsutil)
- [Recharts](https://recharts.org/)
- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [GitHub Actions](https://docs.github.com/en/actions)

---

**Project Status:** ✅ Ready to Start

This plan can be directly handed to Claude Code. Each sprint is clearly defined with concrete tasks. Implementation can happen step by step with tests and CI/CD from the beginning.
