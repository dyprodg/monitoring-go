# MVP Implementation Plan - Clean Version
## Interactive System Monitoring Dashboard

**Version:** 2.0 (Cleaned & Consolidated)
**Last Updated:** 2025-01-09
**Status:** READY TO START
**Timeline:** 4-6 weeks realistic
**Goal:** Working local demo with Docker

---

## 🎯 What We're Building (MVP v1.0)

**Core Concept:**
```
User clicks button → Backend generates load → Metrics spike → User sees it happen
```

**Success Criteria:**
```bash
docker-compose up
# Open http://localhost:3000
# Click "🔥 CPU Stress"
# Watch CPU graph spike to ~90%
# System doesn't crash
= SUCCESS!
```

---

## 📦 MVP Scope (LOCKED - See DECISIONS.md)

### ✅ IN SCOPE

**Backend:**
- 4 system metrics (CPU, Memory, Disk I/O, Network)
- 4 load actions (one per metric)
- REST API (JSON endpoints)
- HTTP polling (NO WebSocket)
- In-memory storage (NO database)
- Safety limits (prevent crashes)
- Action cancellation
- Tests (70% minimum coverage)

**Frontend:**
- React 18 + Vite + TailwindCSS
- Simple dashboard layout
- 4 metric cards with line charts
- 4 action buttons
- Basic event log (last 10 events)
- Polling (every 1 second)
- Dark theme
- Responsive layout

**Infrastructure:**
- Docker + Docker Compose
- Local development only
- Go backend
- React frontend

### ❌ OUT OF SCOPE

**Deferred to Later Versions:**
- ❌ WebSocket (use polling - can add in v1.1)
- ❌ PostgreSQL database (in-memory only - add in v1.4)
- ❌ Pre-built scenarios (manual actions only - add in v1.3)
- ❌ Animated gauges (simple charts - add in v1.1)
- ❌ Active actions widget (basic only - enhance in v1.1)
- ❌ Toast notifications (simple log - add in v1.1)
- ❌ User authentication (public access OK)
- ❌ Cloud deployment (local only - add in v1.2+)
- ❌ Historical data export (no persistence)

**See FUTURE-FEATURES.md for full list**

---

## 🗂️ Project Structure

```
monitoring-dashboard/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go           # Entry point
│   ├── internal/
│   │   ├── metrics/              # System metrics collection
│   │   │   ├── collector.go      # Main collector
│   │   │   ├── cpu.go           # CPU metrics
│   │   │   ├── memory.go        # Memory metrics
│   │   │   ├── disk.go          # Disk I/O metrics
│   │   │   └── network.go       # Network metrics
│   │   ├── actions/              # Load generation
│   │   │   ├── engine.go         # Action engine + safety
│   │   │   ├── cpu_stress.go    # CPU load generator
│   │   │   ├── memory_surge.go  # Memory load generator
│   │   │   ├── disk_storm.go    # Disk I/O generator
│   │   │   └── traffic_flood.go # Network traffic generator
│   │   └── api/                  # HTTP handlers
│   │       ├── handlers.go       # API handlers
│   │       ├── routes.go         # Route definitions
│   │       └── middleware.go     # CORS, logging
│   ├── pkg/
│   │   └── models/               # Shared types
│   │       ├── metrics.go
│   │       ├── action.go
│   │       └── event.go
│   ├── tests/                    # Integration tests
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile                # Multi-stage build
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Dashboard.jsx     # Main layout
│   │   │   ├── MetricCard.jsx    # Metric display + chart
│   │   │   ├── ActionButton.jsx  # Action trigger button
│   │   │   └── EventLog.jsx      # Event feed
│   │   ├── services/
│   │   │   └── api.js            # API client (polling)
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── Dockerfile                # Nginx serve
│
├── docker-compose.yml            # Local development (NO database!)
├── .gitignore
├── README.md
│
└── docs/
    ├── DECISIONS.md              # Key decisions (READ THIS!)
    ├── AI-DEVELOPMENT-RULES.md   # Rules for AI
    ├── MASTER-PLAN.md            # Overview
    ├── VERSION-ROADMAP.md        # Version strategy
    └── archive/                  # Old documents
```

---

## 🚨 CRITICAL SAFETY REQUIREMENTS

**These MUST be implemented before any load action:**

```go
// internal/actions/engine.go

const (
    // Maximum limits (NEVER exceed these)
    MAX_CPU_PERCENT      = 95    // 95% max CPU
    MAX_CPU_DURATION     = 30    // 30 seconds max
    MAX_MEMORY_PERCENT   = 25    // 25% of total RAM max
    MAX_MEMORY_DURATION  = 60    // 60 seconds max
    MAX_DISK_SIZE_MB     = 100   // 100MB temp files max
    MAX_CONCURRENT       = 5     // 5 actions max at once

    // Emergency shutdown thresholds
    CRITICAL_CPU         = 98    // Kill action if CPU hits 98%
    CRITICAL_MEMORY      = 95    // Kill action if memory hits 95%
)

// MANDATORY: All actions must implement
type Action interface {
    Execute(ctx context.Context) error  // Must respect ctx.Done()
    Cleanup() error                      // Must cleanup resources
}
```

**Testing Requirements:**
- ✅ All safety limits MUST have 100% test coverage
- ✅ Emergency shutdown MUST be tested
- ✅ Action cancellation MUST work within 1 second
- ✅ Resource cleanup MUST be verified

**See AI-DEVELOPMENT-RULES.md Rules 1-4 for details**

---

## 📅 4-6 Week Implementation Plan

### Week 1: Backend Foundation
**Goal:** Metrics collection + CPU action working safely

#### Day 1-2: Project Setup + Metrics
```bash
# Tasks
[ ] Initialize Go project structure
[ ] Install dependencies (chi, gopsutil)
[ ] Implement CPU metrics collector
[ ] Implement Memory metrics collector
[ ] Write metrics tests
[ ] Create simple REST API endpoint: GET /api/metrics

# Checkpoint
✅ Can run: go run cmd/server/main.go
✅ Can curl: curl http://localhost:8080/api/metrics
✅ Returns real CPU and memory data
✅ Tests pass
```

**Deliverable:** Backend collects and serves metrics

#### Day 3-4: Safety + CPU Action
```bash
# Tasks
[ ] Create action engine with safety limits
[ ] Implement CPU stress action
[ ] Add emergency shutdown logic
[ ] Write safety tests (100% coverage required!)
[ ] Create API endpoint: POST /api/actions/cpu-stress

# Checkpoint
✅ Can trigger CPU stress via API
✅ CPU increases to ~90%
✅ Safety limits prevent >95%
✅ Action stops after duration
✅ Can cancel action mid-execution
✅ System doesn't crash
✅ All safety tests pass
```

**Deliverable:** CPU stress action works safely

#### Day 5: REST API + Testing
```bash
# Tasks
[ ] Complete REST API endpoints
    - GET /api/health
    - GET /api/metrics
    - POST /api/actions/cpu-stress
    - GET /api/actions/active
    - DELETE /api/actions/:id/stop
[ ] Add CORS middleware
[ ] Write API tests
[ ] Integration testing

# Checkpoint
✅ All endpoints work
✅ CORS configured for localhost:3000
✅ Tests >70% coverage
✅ API documented
```

**Week 1 Deliverable:** ✅ Backend works, CPU action safe and functional

---

### Week 2: Frontend + Integration
**Goal:** Working UI with ONE button

#### Day 1-2: Frontend Setup
```bash
# Tasks
[ ] Initialize Vite + React project
[ ] Install dependencies (Recharts, TailwindCSS)
[ ] Create basic dashboard layout
[ ] Create API client service (polling)
[ ] Test API connection

# Checkpoint
✅ Frontend runs: npm run dev
✅ Can call backend API
✅ CORS works
✅ Basic layout shows
```

**Deliverable:** Frontend connects to backend

#### Day 3-4: First Working Button
```bash
# Tasks
[ ] Create MetricCard component with line chart
[ ] Create ActionButton component
[ ] Implement polling (every 1 second)
[ ] Connect CPU button to API
[ ] Display metrics in real-time

# Checkpoint
✅ Open http://localhost:3000
✅ Click "🔥 CPU Stress"
✅ CPU chart updates
✅ Line shows spike to ~90%
✅ Returns to normal after 10s
```

**Deliverable:** ONE button works end-to-end! 🎉

#### Day 5: Docker Setup
```bash
# Tasks
[ ] Create backend Dockerfile
[ ] Create frontend Dockerfile
[ ] Create docker-compose.yml (backend + frontend ONLY)
[ ] Test docker-compose up
[ ] Verify end-to-end in Docker

# docker-compose.yml
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

# NO PostgreSQL in MVP!

# Checkpoint
✅ docker-compose up works
✅ Can access http://localhost:3000
✅ Button still works in Docker
```

**Week 2 Deliverable:** ✅ ONE button works perfectly in Docker

---

### Week 3: Complete All Actions
**Goal:** All 4 actions working

#### Day 1-2: Remaining Metrics + Actions
```bash
# Tasks
[ ] Implement Disk I/O metrics collector
[ ] Implement Network metrics collector
[ ] Implement Memory surge action
[ ] Implement Disk storm action
[ ] Implement Traffic flood action
[ ] Add API endpoints for new actions
[ ] Write tests for all actions

# Checkpoint
✅ All 4 metrics collecting
✅ All 4 actions implemented
✅ All actions have safety limits
✅ All actions can be cancelled
✅ Tests pass
```

**Deliverable:** All backend actions work

#### Day 3: Frontend - All Metrics
```bash
# Tasks
[ ] Add Disk I/O metric card
[ ] Add Network metric card
[ ] Update Memory and CPU cards
[ ] Style with TailwindCSS
[ ] Ensure responsive layout

# Checkpoint
✅ All 4 metric cards visible
✅ All show real-time data
✅ Charts update smoothly
```

**Deliverable:** All metrics displayed

#### Day 4: Frontend - All Buttons
```bash
# Tasks
[ ] Add Memory surge button
[ ] Add Disk storm button
[ ] Add Traffic flood button
[ ] Test all buttons
[ ] Add loading states
[ ] Add error handling

# Checkpoint
✅ All 4 buttons work
✅ All trigger correct actions
✅ Loading states show
✅ Errors display nicely
```

**Deliverable:** All buttons work

#### Day 5: Event Log
```bash
# Tasks
[ ] Create event system in backend
[ ] Add events to API (GET /api/events)
[ ] Create EventLog component
[ ] Poll for events
[ ] Style event log

# Checkpoint
✅ Events show action lifecycle
✅ Timestamps correct
✅ Auto-scroll to newest
✅ Color-coded by type
```

**Week 3 Deliverable:** ✅ All 4 actions + event log working

---

### Week 4: Charts, Polish & Testing
**Goal:** Demo-ready application

#### Day 1-2: Charts with Recharts
```bash
# Tasks
[ ] Implement proper LineChart for all metrics
[ ] Add 60-second sliding window
[ ] Configure chart styling
[ ] Add axis labels
[ ] Optimize chart performance

# Checkpoint
✅ Charts look professional
✅ Updates smooth (no flicker)
✅ 60s of history visible
✅ Auto-scaling works
```

**Deliverable:** Professional charts

#### Day 3: Visual Polish
```bash
# Tasks
[ ] Dark theme polish
[ ] Improve spacing and layout
[ ] Add loading states everywhere
[ ] Improve error messages
[ ] Button hover effects
[ ] Responsive testing (desktop/laptop)

# Checkpoint
✅ Looks professional
✅ Dark theme consistent
✅ No UI glitches
✅ Works on different screen sizes
```

**Deliverable:** Professional UI

#### Day 4: Testing & Bug Fixes
```bash
# Tasks
[ ] Write frontend component tests
[ ] E2E test: Button click → Metric change
[ ] Run all backend tests
[ ] Check test coverage (>70%)
[ ] Fix any failing tests
[ ] Fix any bugs found

# Checkpoint
✅ All tests pass
✅ Coverage >70% overall
✅ Safety tests at 100%
✅ No critical bugs
```

**Deliverable:** Tests pass, bugs fixed

#### Day 5: Documentation & Demo Prep
```bash
# Tasks
[ ] Write README.md with:
    - Project description
    - Setup instructions
    - How to run
    - How to demo
[ ] Test on fresh machine
[ ] Practice 5-minute demo
[ ] Record demo video (optional)
[ ] Final testing

# Checkpoint
✅ README complete
✅ Setup works from scratch
✅ Can do 5-min demo
✅ No crashes during demo
```

**Week 4 Deliverable:** ✅ MVP COMPLETE - DEMO READY! 🎉

---

## 🧪 Testing Strategy

### Unit Tests (Backend)
```go
// Minimum required tests

// Metrics
TestCPUMetricsCollected()
TestMemoryMetricsCollected()
TestDiskMetricsCollected()
TestNetworkMetricsCollected()

// Actions
TestCPUStressIncreasesMetrics()
TestMemorySurgeIncreasesMemory()
TestDiskStormGeneratesIO()
TestTrafficFloodGeneratesRequests()

// Safety (CRITICAL - 100% coverage required)
TestSafetyLimitsEnforced()           // ← CRITICAL
TestActionStopsGracefully()          // ← CRITICAL
TestEmergencyShutdownWorks()         // ← CRITICAL
TestCPUNeverExceeds95Percent()       // ← CRITICAL
TestMemoryNeverExceeds25Percent()    // ← CRITICAL
TestActionCancelsWithin1Second()     // ← CRITICAL
TestResourceCleanupAfterAction()     // ← CRITICAL
TestMultipleConcurrentActionsSafe()

// API
TestHealthEndpoint()
TestMetricsEndpoint()
TestActionEndpoints()
TestCORSConfigured()
```

### Integration Tests
```bash
# E2E test scenario
1. Start backend
2. Trigger CPU stress action
3. Poll metrics endpoint
4. Verify CPU increased
5. Wait for completion
6. Verify CPU returned to normal
7. Check no temp files left
```

### Frontend Tests
```javascript
// Component tests
TestMetricCardDisplaysData()
TestActionButtonTriggersAPI()
TestChartUpdatesOnNewData()
TestEventLogDisplaysEvents()
TestErrorHandling()
```

### Coverage Targets
```
Overall:        70% minimum, 75% target
Safety code:    100% required (no exceptions)
Actions:        80% minimum
Metrics:        75% minimum
API handlers:   70% minimum
```

---

## 📡 API Specification

### Metrics
```http
GET /api/health
Response: {"status": "healthy", "timestamp": "2025-01-09T10:00:00Z"}

GET /api/metrics
Response: {
  "timestamp": "2025-01-09T10:00:00Z",
  "cpu": 45.2,
  "memory": 62.5,
  "disk_io": 150.3,
  "network": 2.4
}
```

### Actions
```http
POST /api/actions/cpu-stress
Body: {
  "target_percent": 90,
  "duration_seconds": 10
}
Response: {
  "id": "cpu-abc123",
  "status": "started",
  "started_at": "2025-01-09T10:00:00Z"
}

POST /api/actions/memory-surge
Body: {
  "size_mb": 500,
  "duration_seconds": 60
}
Response: {
  "id": "mem-def456",
  "status": "started"
}

POST /api/actions/disk-storm
Body: {
  "operations": 1000,
  "file_size_kb": 10
}

POST /api/actions/traffic-flood
Body: {
  "requests_per_sec": 100,
  "duration_seconds": 10
}

GET /api/actions/active
Response: {
  "actions": [
    {
      "id": "cpu-abc123",
      "type": "cpu-stress",
      "progress": 0.7,
      "started_at": "2025-01-09T10:00:00Z"
    }
  ]
}

DELETE /api/actions/:id/stop
Response: {"status": "stopped"}
```

### Events
```http
GET /api/events
Response: {
  "events": [
    {
      "timestamp": "2025-01-09T10:00:00Z",
      "type": "action_started",
      "severity": "info",
      "message": "CPU stress action started (90%, 10s)"
    },
    {
      "timestamp": "2025-01-09T10:00:10Z",
      "type": "action_completed",
      "severity": "success",
      "message": "CPU stress action completed"
    }
  ]
}
```

---

## 🛠️ Technology Stack

### Backend
- **Language:** Go 1.21+
- **Router:** Chi (lightweight HTTP router)
- **Metrics:** gopsutil (cross-platform system metrics)
- **Testing:** Go standard testing
- **No WebSocket:** Using HTTP polling (simpler)
- **No Database:** In-memory storage only

### Frontend
- **Framework:** React 18
- **Build Tool:** Vite (faster than CRA)
- **Styling:** TailwindCSS
- **Charts:** Recharts (React + D3)
- **HTTP Client:** Fetch API (no axios needed)
- **State:** React Hooks (useState, useEffect)

### Infrastructure
- **Containerization:** Docker
- **Orchestration:** Docker Compose
- **Development:** Local machine
- **No Cloud:** MVP runs locally only

---

## 🚀 Quick Start (After Build)

```bash
# Clone repo
git clone <repo>
cd monitoring-dashboard

# Start with Docker
docker-compose up

# OR run locally (development)
# Terminal 1 - Backend
cd backend
go run cmd/server/main.go

# Terminal 2 - Frontend
cd frontend
npm install
npm run dev

# Access
Frontend: http://localhost:3000
Backend:  http://localhost:8080/api/health
```

---

## 🎬 5-Minute Demo Script

**Minute 1: Introduction (0:30)**
```
"This is an interactive system monitoring dashboard.
Unlike traditional monitoring tools, you can actively
trigger different types of system load and watch the
metrics react in real-time."
```

**Minute 2: Single Action Demo (1:30)**
```
"Let me show you CPU stress..."

→ Click "🔥 CPU Stress" button
→ Watch CPU gauge and chart
→ Point out event log entry
→ Watch it complete after 10 seconds

"Notice how it safely went to 90%, held for the
configured duration, then returned to normal.
The system enforces safety limits to prevent crashes."
```

**Minute 3: Multiple Actions (1:30)**
```
"Now let's trigger multiple types of load simultaneously..."

→ Click CPU, Memory, and Network buttons quickly
→ Watch all metrics spike
→ Show event log tracking all actions

"The system handles concurrent load safely with
built-in limits and proper resource management."
```

**Minute 4: Stop Action (0:30)**
```
→ Click stop on one action
→ Show it stops within 1 second
→ Metrics for that type return to normal

"All actions can be cancelled at any time with
proper cleanup."
```

**Minute 5: Technology Overview (1:00)**
```
"Built with:
- Go backend using goroutines for concurrent load generation
- React frontend with real-time updates via polling
- Docker for easy deployment
- Comprehensive tests including safety limits
- All open source and on GitHub"
```

---

## ✅ MVP Complete Checklist

### Functionality
- [ ] `docker-compose up` starts successfully
- [ ] Frontend accessible at http://localhost:3000
- [ ] All 4 metrics display real data
- [ ] All 4 action buttons work
- [ ] Metrics update every 1 second
- [ ] CPU stress reaches ~90%
- [ ] Memory surge allocates memory
- [ ] Disk storm shows I/O activity
- [ ] Traffic flood shows network activity
- [ ] Event log shows action lifecycle
- [ ] Can stop actions individually
- [ ] Charts display and update smoothly
- [ ] Safety limits prevent crashes

### Quality
- [ ] No errors in browser console
- [ ] No crashes when running all 4 actions
- [ ] Backend tests pass (>70% coverage)
- [ ] Safety tests pass (100% coverage)
- [ ] Frontend works in Chrome/Firefox
- [ ] Responsive layout works
- [ ] UI looks professional

### Documentation
- [ ] README has setup instructions
- [ ] README has demo guide
- [ ] Code has basic comments
- [ ] API endpoints documented

### Demo Readiness
- [ ] Can complete 5-minute demo without issues
- [ ] Actions produce visible effects
- [ ] UI is impressive
- [ ] No obvious bugs

---

## 🚦 Go/No-Go Decision Points

### After Week 1
**GO if:**
- ✅ Backend collects metrics reliably
- ✅ CPU stress action works
- ✅ Safety limits enforced
- ✅ Tests pass

**NO-GO if:**
- ❌ Can't collect metrics
- ❌ CPU stress crashes system
- ❌ Safety limits don't work

**Action if NO-GO:** Debug for 2 days, then get human help

### After Week 2
**GO if:**
- ✅ Frontend displays metrics
- ✅ One button works end-to-end
- ✅ Docker-compose functional

**NO-GO if:**
- ❌ Frontend can't connect to backend
- ❌ Button doesn't trigger action
- ❌ Major architectural issues

**Action if NO-GO:** Get human help, may need architecture change

### After Week 3
**GO if:**
- ✅ All 4 actions work
- ✅ Event log functional
- ✅ Can handle concurrent actions

**NO-GO if:**
- ❌ Actions crash system
- ❌ Major bugs in core functionality

**Action if NO-GO:** Reduce scope (remove problematic actions)

---

## 📞 When to Ask for Human Help

### 🔴 STOP and Ask Human:
- System crashes during testing (safety limits failed)
- Can't stop actions (cancellation broken)
- Can't decide between two technical approaches
- Performance terrible (>5s to show metrics)
- Docker issues after 4 hours of debugging
- Major security concern discovered

### 🟡 Ask When Convenient:
- UI/UX design decisions
- Color scheme choices
- Feature prioritization questions
- Testing strategy validation

### 🟢 AI Can Decide:
- Implementation details
- File organization
- Variable names
- Code formatting
- Test case specifics
- Documentation style

---

## 🎯 Next Steps After MVP

**After v1.0 Complete:**

1. **Demo and Validate**
   - Show to friends/colleagues
   - Get feedback
   - Identify improvements

2. **Decide Next Step:**
   - Option A: Polish (v1.1) → Showcase (v1.2) → DONE
   - Option B: Polish (v1.1) → AWS (v2.0+)
   - See VERSION-ROADMAP.md for details

3. **Update Portfolio**
   - Add to GitHub
   - Record demo video
   - Write blog post

---

## 📚 Reference Documents

**Before Starting:**
- [ ] Read DECISIONS.md (key architectural decisions)
- [ ] Read AI-DEVELOPMENT-RULES.md (how to build)
- [ ] Read this file (what to build)

**During Development:**
- Check DECISIONS.md for any uncertainty
- Follow AI-DEVELOPMENT-RULES.md standards
- Update progress in this file

**After MVP:**
- See VERSION-ROADMAP.md for next steps
- See FUTURE-FEATURES.md for enhancement ideas
- See RELEASE-PLAN.md for deployment options

---

## 🎯 Critical Reminders

**For AI:**
1. Read DECISIONS.md first - it locks all contradictions
2. NO WebSocket - use polling
3. NO database - in-memory only
4. NO scenarios - manual actions only
5. Safety limits are MANDATORY with 100% test coverage
6. Ask human if blocked >4 hours

**For Human:**
1. This is the ONLY plan to follow
2. Timeline is 4-6 weeks realistic (not guaranteed 4)
3. Review progress every 2 days
4. Make go/no-go decisions at checkpoints
5. Refer to DECISIONS.md if AI asks questions

---

**Status:** READY TO START
**Timeline:** 4-6 weeks
**First Task:** Week 1, Day 1 - Initialize Go project + CPU metrics
**Next Checkpoint:** After Day 2 - Backend metrics working

**Let's build something impressive! 🚀**
