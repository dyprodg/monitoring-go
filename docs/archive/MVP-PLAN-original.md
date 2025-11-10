# System Monitor Dashboard - MVP Plan
## "Click Button, See Metrics React" - 4 Weeks

**Last Updated:** 2025-01-09
**Status:** READY TO START
**Goal:** Impressive interactive demo running locally with docker-compose

---

## 🎯 MVP Success Criteria

```
USER CLICKS "CPU STRESS" → CPU GRAPH SPIKES TO 90% IN 2 SECONDS
```

**That's it. Everything else supports this core interaction.**

---

## 📦 MVP Scope (ONLY These Features)

### ✅ MUST BUILD (Week 1-4)

#### 1. **Real-Time Metrics Collection**
- CPU usage (percentage + per-core)
- Memory usage (MB/GB + percentage)
- Disk I/O (read/write MB/s)
- Network traffic (in/out MB/s)
- **Update frequency:** 1 second
- **Storage:** In-memory only (60 seconds history)

#### 2. **4 Interactive Load Actions**
```
🔥 CPU Stress    - 90% for 10s (configurable)
💾 Memory Surge  - Allocate 500MB (configurable)
💿 Disk Storm    - 1000 file I/O operations
🌐 Traffic Flood - 100 HTTP requests/sec for 10s
```

#### 3. **Simple Web Dashboard**
- 4 action buttons (instant trigger)
- 4 metric charts (real-time updates)
- Active actions widget (progress bars)
- Event log (last 50 events)
- **Design:** Clean, dark theme, responsive

#### 4. **Live Updates**
- Polling every 1 second (no WebSocket needed for MVP)
- JSON REST API
- CORS enabled for local dev

#### 5. **Safety & Stability**
- ⚠️ **CRITICAL:** Safety limits enforced
- Graceful action cancellation
- Resource cleanup
- Error handling

---

## 🛡️ MANDATORY SAFETY LIMITS

```go
// These MUST be implemented before load actions
const (
    MAX_CPU_PERCENT      = 95    // Never exceed this
    MAX_CPU_DURATION     = 30    // Max 30 seconds
    MAX_MEMORY_PERCENT   = 25    // Max 25% of system RAM
    MAX_MEMORY_DURATION  = 60    // Max 60 seconds
    MAX_DISK_SIZE_MB     = 100   // Max 100MB temp files
    MAX_CONCURRENT       = 5     // Max 5 actions at once
)

// Auto-shutdown if system exceeds:
CRITICAL_CPU_THRESHOLD     = 98%   // Kill action if hit
CRITICAL_MEMORY_THRESHOLD  = 95%   // Kill action if hit
```

**WHY:** Without these, you could crash the system during demo.

---

## 🗂️ Project Structure (Minimal)

```
monitoring-dashboard/
├── backend/
│   ├── cmd/server/main.go           # Entry point
│   ├── internal/
│   │   ├── metrics/                 # Metrics collection
│   │   │   ├── collector.go
│   │   │   ├── cpu.go
│   │   │   └── memory.go
│   │   ├── actions/                 # Load generators
│   │   │   ├── engine.go            # Action engine + safety
│   │   │   ├── cpu_stress.go
│   │   │   ├── memory_surge.go
│   │   │   ├── disk_storm.go
│   │   │   └── traffic_flood.go
│   │   └── api/                     # REST API
│   │       ├── handlers.go
│   │       └── routes.go
│   ├── go.mod
│   └── Dockerfile
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Dashboard.jsx       # Main layout
│   │   │   ├── MetricCard.jsx      # CPU/Memory/Disk/Network
│   │   │   ├── ActionButton.jsx    # Trigger buttons
│   │   │   ├── ActiveActions.jsx   # Running actions
│   │   │   └── EventLog.jsx        # Event feed
│   │   ├── services/
│   │   │   └── api.js              # API client
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── Dockerfile
│
├── docker-compose.yml               # Local dev environment
├── README.md
└── MVP-PLAN.md                      # This file
```

**REMOVED from MVP:**
- ❌ PostgreSQL database
- ❌ WebSocket (use polling)
- ❌ Terraform/AWS
- ❌ GitHub Actions CI/CD
- ❌ Scenarios system
- ❌ Historical data storage
- ❌ Authentication

---

## 📋 Implementation Plan (4 Weeks)

### 🔵 Week 1: Backend Foundation
**Goal:** Metrics collection + ONE working action

#### Day 1-2: Setup & Metrics
- [ ] Initialize Go project structure
- [ ] Install dependencies (chi, gopsutil)
- [ ] Implement CPU metrics collector
- [ ] Implement Memory metrics collector
- [ ] **Test:** Can collect metrics every 1 second

#### Day 3-4: Load Engine + Safety
- [ ] Create action engine with safety limits
- [ ] Implement CPU stress action
- [ ] Add emergency shutdown logic
- [ ] **Test:** CPU stress works without crashing system

#### Day 5: REST API
- [ ] Create REST API with Chi router
- [ ] `GET /api/metrics` - Current metrics
- [ ] `POST /api/actions/cpu-stress` - Trigger action
- [ ] `GET /api/actions/active` - List running
- [ ] `DELETE /api/actions/:id` - Stop action
- [ ] **Test:** Can curl API and trigger CPU stress

**Week 1 Checkpoint:** ✅ Backend works, CPU action safe and functional

---

### 🟢 Week 2: Frontend + Integration
**Goal:** Working UI with ONE button

#### Day 1-2: Frontend Setup
- [ ] Initialize React + Vite project
- [ ] Setup TailwindCSS
- [ ] Create basic dashboard layout
- [ ] Create API client service
- [ ] **Test:** Can call backend API from frontend

#### Day 3-4: Interactive UI
- [ ] Create MetricCard component (gauges)
- [ ] Create ActionButton component
- [ ] Implement polling (every 1 second)
- [ ] Connect button to API
- [ ] **Test:** Click button → CPU spikes → UI updates

#### Day 5: Docker Integration
- [ ] Create backend Dockerfile
- [ ] Create frontend Dockerfile
- [ ] Create docker-compose.yml
- [ ] **Test:** `docker-compose up` works end-to-end

**Week 2 Checkpoint:** ✅ ONE button works perfectly (CPU stress)

---

### 🟡 Week 3: Complete All Actions
**Goal:** All 4 actions working

#### Day 1: Memory & Disk Actions
- [ ] Implement Disk I/O metrics collector
- [ ] Implement Network metrics collector
- [ ] Implement Memory surge action
- [ ] Implement Disk storm action
- [ ] Add API endpoints for new actions
- [ ] **Test:** Memory and Disk actions work

#### Day 2: Network Action
- [ ] Implement Traffic flood action
- [ ] Add API endpoint
- [ ] **Test:** All 4 actions work independently

#### Day 3: Frontend - All Metrics
- [ ] Add Disk and Network metric cards
- [ ] Add all 4 action buttons
- [ ] Style with TailwindCSS
- [ ] **Test:** All 4 buttons work

#### Day 4-5: Active Actions Widget + Event Log
- [ ] Create ActiveActions component (progress bars)
- [ ] Create EventLog component
- [ ] Add event system to backend
- [ ] Connect event stream to frontend
- [ ] **Test:** Can see running actions and events

**Week 3 Checkpoint:** ✅ All 4 actions + event log working

---

### 🟣 Week 4: Charts, Polish & Testing
**Goal:** Demo-ready polish

#### Day 1-2: Charts with Recharts
- [ ] Install Recharts
- [ ] Create LineChart component (60s history)
- [ ] Create Gauge component (current value)
- [ ] Add charts to all 4 metric cards
- [ ] **Test:** Charts update smoothly

#### Day 3: Visual Polish
- [ ] Dark theme styling
- [ ] Loading states
- [ ] Error messages
- [ ] Button animations
- [ ] Responsive layout
- [ ] **Test:** UI looks professional

#### Day 4: Testing & Documentation
- [ ] Write backend unit tests (>70% coverage)
- [ ] Write frontend component tests
- [ ] E2E test: Button click → Metric change
- [ ] Write README with setup instructions
- [ ] **Test:** All tests pass

#### Day 5: Demo Preparation
- [ ] Test on fresh machine
- [ ] Fix any bugs
- [ ] Performance check (can handle all 4 actions simultaneously)
- [ ] Record demo video
- [ ] **Test:** Full demo runs smoothly

**Week 4 Checkpoint:** ✅ MVP COMPLETE - DEMO READY

---

## 🔴 Critical Path (Blocker If Not Working)

```
Day 1-2:   Metrics collection MUST work
    ↓
Day 3-4:   CPU stress MUST work safely
    ↓
Day 5:     REST API MUST work
    ↓
Week 2:    ONE button MUST work end-to-end
    ↓
Week 3:    All 4 actions MUST work
    ↓
Week 4:    POLISH & DEMO READY
```

**⚠️ RULE:** If any step fails, STOP and fix before moving forward.

---

## 🧪 Testing Requirements

### Minimum Tests (Must Have)

**Backend:**
```go
✅ TestCPUMetricsCollected()
✅ TestMemoryMetricsCollected()
✅ TestCPUStressIncreasesMetrics()
✅ TestMemorySurgeIncreasesMemory()
✅ TestSafetyLimitsEnforced()           // CRITICAL
✅ TestActionStopsGracefully()          // CRITICAL
✅ TestEmergencyShutdownWorks()         // CRITICAL
✅ TestMultipleConcurrentActions()
```

**Frontend:**
```javascript
✅ TestMetricsDisplayUpdates()
✅ TestButtonTriggersAction()
✅ TestActionProgressShown()
✅ TestEventLogUpdates()
```

**Integration:**
```
✅ E2E: docker-compose up works
✅ E2E: Click CPU button → CPU spikes
✅ E2E: All 4 actions work simultaneously
✅ E2E: Can stop actions mid-execution
```

**Coverage Target:** 70% (not 80% for MVP)

---

## 🛠️ Technology Stack (Minimal)

**Backend:**
```
- Go 1.21+
- Chi Router (HTTP)
- gopsutil (system metrics)
- Standard library (no extra dependencies)
```

**Frontend:**
```
- React 18
- Vite (build tool)
- TailwindCSS (styling)
- Recharts (charts)
- Fetch API (no axios needed)
```

**Infrastructure:**
```
- Docker + docker-compose
- No cloud (local only)
- No database
```

---

## 📡 API Specification

### Metrics Endpoints

```http
GET /api/health
Response: { "status": "healthy" }

GET /api/metrics
Response: {
  "cpu": { "percent": 45.2, "cores": [23.1, 67.8, ...] },
  "memory": { "total_mb": 16384, "used_mb": 8192, "percent": 50.0 },
  "disk": { "read_mbps": 12.3, "write_mbps": 5.6 },
  "network": { "in_mbps": 1.2, "out_mbps": 0.8 }
}
```

### Action Endpoints

```http
POST /api/actions/cpu-stress
Body: { "target_percent": 90, "duration_seconds": 10 }
Response: { "id": "cpu-123", "status": "started" }

POST /api/actions/memory-surge
Body: { "size_mb": 500, "duration_seconds": 60 }
Response: { "id": "mem-456", "status": "started" }

POST /api/actions/disk-storm
Body: { "operations": 1000, "file_size_kb": 10 }
Response: { "id": "disk-789", "status": "started" }

POST /api/actions/traffic-flood
Body: { "requests_per_sec": 100, "duration_seconds": 10 }
Response: { "id": "traffic-012", "status": "started" }

GET /api/actions/active
Response: {
  "actions": [
    {
      "id": "cpu-123",
      "type": "cpu-stress",
      "progress": 0.7,
      "started_at": "2025-01-09T10:00:00Z"
    }
  ]
}

DELETE /api/actions/:id
Response: { "status": "stopped" }

GET /api/events
Response: {
  "events": [
    {
      "timestamp": "2025-01-09T10:00:00Z",
      "type": "action_started",
      "severity": "info",
      "message": "CPU stress action started"
    }
  ]
}
```

---

## 🚀 Quick Start (After Build)

```bash
# Clone and start
git clone <repo>
cd monitoring-dashboard
docker-compose up

# Access
Frontend: http://localhost:3000
Backend:  http://localhost:8080/api/health

# Demo
1. Open http://localhost:3000
2. Click "🔥 CPU Stress"
3. Watch CPU graph spike
4. Check event log
5. Stop action
```

---

## 🎬 Demo Script (5 Minutes)

```
1. Show dashboard (0:30)
   "Here's our system monitoring dashboard"

2. Explain metrics (0:30)
   "We track CPU, Memory, Disk I/O, and Network in real-time"

3. CPU Stress Demo (1:00)
   "Click CPU Stress → See it spike to 90%"
   "Progress bar shows time remaining"
   "Event log tracks what's happening"

4. Memory Surge Demo (1:00)
   "Click Memory Surge → Watch memory allocation"

5. Multiple Actions (1:00)
   "Run all 4 actions simultaneously"
   "System handles concurrent load safely"

6. Stop Actions (0:30)
   "Can stop any action individually or all at once"

7. Q&A (0:30)
```

---

## 🐛 Known Limitations (By Design)

- ❌ No historical data persistence (only 60s in memory)
- ❌ No user authentication (demo/development only)
- ❌ No cloud deployment (local only)
- ❌ No pre-built scenarios (just manual actions)
- ❌ No WebSocket (polling is fine for 1s updates)
- ❌ No mobile app (web only)

**These are FEATURES, not BUGS.** They keep the MVP simple.

---

## 💡 Development Tips

### If Something Isn't Working

**Backend not collecting metrics?**
```bash
# Test metrics collection standalone
cd backend
go run cmd/server/main.go
curl http://localhost:8080/api/metrics
```

**Frontend not updating?**
```bash
# Check browser console for errors
# Verify API is accessible
curl http://localhost:8080/api/metrics
```

**Docker not starting?**
```bash
# Check logs
docker-compose logs backend
docker-compose logs frontend

# Rebuild
docker-compose down
docker-compose up --build
```

**Action not stopping?**
- Check context cancellation in Go
- Verify cleanup in defer statements
- Test emergency shutdown

---

## ✅ MVP Complete Checklist

### Functionality
- [ ] Can run `docker-compose up` successfully
- [ ] All 4 metrics display in UI
- [ ] All 4 action buttons work
- [ ] Metrics update every 1 second
- [ ] CPU stress increases CPU to ~90%
- [ ] Memory surge increases memory by ~500MB
- [ ] Disk storm shows I/O activity
- [ ] Traffic flood shows network activity
- [ ] Can see active actions with progress
- [ ] Event log shows action lifecycle
- [ ] Can stop actions individually
- [ ] Can stop all actions at once
- [ ] Charts display and update smoothly
- [ ] Safety limits prevent system crash

### Quality
- [ ] No errors in browser console
- [ ] No crashes when running all 4 actions
- [ ] Backend tests pass (>70% coverage)
- [ ] Frontend tests pass
- [ ] E2E test passes
- [ ] UI looks professional (dark theme)
- [ ] Responsive on laptop screen (mobile optional)

### Documentation
- [ ] README has setup instructions
- [ ] README has demo instructions
- [ ] Code has basic comments
- [ ] API endpoints documented

### Demo Readiness
- [ ] Can complete 5-minute demo without issues
- [ ] Actions produce visible effects
- [ ] UI is impressive
- [ ] No obvious bugs during demo

---

## 🚦 Go/No-Go Decision Points

### After Week 1
**GO if:**
- ✅ Backend collects metrics
- ✅ CPU stress action works
- ✅ REST API functional

**NO-GO if:**
- ❌ Can't collect metrics reliably
- ❌ CPU stress crashes system
- ❌ Major architectural issues

**Action if NO-GO:** Debug for 2 more days, then reassess approach

### After Week 2
**GO if:**
- ✅ Frontend displays metrics
- ✅ One button works end-to-end
- ✅ Docker-compose functional

**NO-GO if:**
- ❌ Frontend can't connect to backend
- ❌ Button doesn't trigger action
- ❌ Docker issues persist

**Action if NO-GO:** Simplify architecture further

### After Week 3
**GO if:**
- ✅ All 4 actions work
- ✅ Event log functional
- ✅ Can handle concurrent actions

**NO-GO if:**
- ❌ Actions crash the system
- ❌ Major bugs in core functionality

**Action if NO-GO:** Cut scope (remove problematic actions)

---

## 📞 When to Ask for Human Help

### 🔴 STOP and Ask Human:
- System keeps crashing during testing
- Safety limits don't prevent system freeze
- Can't decide between two technical approaches
- Performance is terrible (>5s update lag)
- Docker issues after 4 hours of debugging
- Major security concern discovered

### 🟡 Ask Human When Convenient:
- UI/UX feedback needed
- Color scheme decisions
- Feature prioritization
- Testing strategy validation

### 🟢 AI Can Decide:
- Implementation details
- File structure
- Variable names
- Code organization
- Test cases
- Documentation

---

## 🎯 Next Steps

1. **Human Review This Plan**
   - Approve scope
   - Approve safety limits
   - Set timeline expectations

2. **Start Week 1, Day 1**
   - Initialize Go backend
   - Implement CPU metrics
   - Get human approval after Day 2

3. **Weekly Check-ins**
   - End of Week 1: Backend demo
   - End of Week 2: Frontend demo
   - End of Week 3: All actions demo
   - End of Week 4: Final MVP demo

---

**This is the plan AI should follow. No database, no AWS, no extra features. Just the core interaction working perfectly.**

**Ready to start? Begin with Week 1, Day 1.**
