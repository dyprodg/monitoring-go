# Project Decisions - Final & Binding
## Key Architectural and Scope Decisions

**Date:** 2025-01-09
**Status:** LOCKED - Do not change without discussion
**Purpose:** Single source of truth for all contradictions found during review

---

## 🎯 Core Architectural Decisions

### Decision #1: Real-Time Updates - POLLING (Not WebSocket)

**Decision:** Use HTTP polling for MVP, not WebSocket

**Rationale:**
- ✅ Simpler to implement (1 day vs 3-4 days)
- ✅ Easier to debug
- ✅ 1-second polling is acceptable for demo
- ✅ Can upgrade to WebSocket in v1.1 if needed
- ✅ No connection management complexity

**Implementation:**
```javascript
// Frontend polls every 1 second
setInterval(() => {
    fetch('/api/metrics')
        .then(res => res.json())
        .then(data => updateCharts(data))
}, 1000);
```

**Future:** WebSocket can be added in v1.1 as enhancement

---

### Decision #2: Database - NONE in MVP

**Decision:** NO PostgreSQL in MVP. In-memory storage only.

**Rationale:**
- ✅ Keeps MVP simple
- ✅ No migration/schema management
- ✅ Faster development (no DB setup)
- ✅ 60 seconds of metrics history is enough for demo
- ✅ Can add database in v1.4 when historical data needed

**Implementation:**
```go
// In-memory circular buffer
type MetricsStore struct {
    data []Metrics  // Last 60 data points
    mu   sync.RWMutex
}
```

**Docker Compose:**
```yaml
# MVP - NO PostgreSQL service
services:
  backend:
    # ...
  frontend:
    # ...
  # postgres: NOT INCLUDED
```

**Future:** Add PostgreSQL in v1.4 for historical data

---

### Decision #3: Pre-built Scenarios - NOT in MVP

**Decision:** NO pre-built scenarios in MVP. Manual actions only.

**Rationale:**
- ✅ MVP proves core concept (actions work)
- ✅ Scenarios add complexity (timeline, sequencing)
- ✅ Save 1-2 weeks of development
- ✅ Can validate MVP before adding scenarios

**MVP Scope:**
```
✅ 4 manual actions (CPU, Memory, Disk, Network)
✅ User clicks button → Action executes
✅ Metrics react in real-time

❌ NO pre-built scenarios
❌ NO timeline visualization
❌ NO auto-play sequences
```

**Future:** Add 3 scenarios in v1.3

---

### Decision #4: Project Structure - internal/actions/

**Decision:** Use `internal/actions/` not `internal/loadgen/`

**Rationale:**
- ✅ Clearer naming (matches domain language)
- ✅ User triggers "actions" not "loadgen"
- ✅ Simpler to understand

**Directory Structure:**
```
backend/
  cmd/
    server/
      main.go
  internal/
    metrics/        ← System metrics collection
    actions/        ← Load generation actions
    api/            ← HTTP handlers
  pkg/
    models/         ← Shared types
  go.mod
```

---

### Decision #5: Test Coverage - 70% Minimum

**Decision:** 70% minimum, 75% target for MVP

**Breakdown:**
```
MVP (v1.0):
- Overall: 70% minimum, 75% target
- Safety code: 100% required
- Actions: 80% minimum
- API handlers: 70% minimum
- Utilities: 60% minimum

Post-MVP (v1.1+):
- Overall: 75% minimum, 80% target
```

**Critical Tests (Must be 100%):**
- Safety limit enforcement
- Action cancellation
- Emergency shutdown
- Resource cleanup

---

## 📅 Timeline Decisions

### Decision #6: Realistic Timeline - 4-6 Weeks

**Decision:** Plan for 5-6 weeks, aim for 4 weeks

**Breakdown:**
```
Optimistic (4 weeks):
- Full-time work (40 hrs/week)
- Experience with Go/React
- No major blockers

Realistic (5-6 weeks):
- Part-time work (20 hrs/week)
- Some learning curve
- Normal debugging time

Conservative (8 weeks):
- First time with Go
- Limited time (10 hrs/week)
- Includes learning
```

**Weekly Milestones:**
```
Week 1: Backend + CPU action working
Week 2: Frontend + one button working
Week 3: All 4 actions working
Week 4: Polish + testing
Week 5: Buffer for issues
Week 6: Final testing + documentation
```

---

## 🎨 UI/UX Decisions

### Decision #7: MVP UI - Simple but Functional

**Decision:** Simple UI for MVP, polish in v1.1

**MVP (v1.0) UI:**
```
✅ 4 metric cards (simple numbers + line charts)
✅ 4 action buttons
✅ Basic event log (last 10 events)
✅ Dark theme
✅ Responsive layout

❌ NO animated gauges (v1.1)
❌ NO fancy progress bars (v1.1)
❌ NO active actions widget (v1.1)
❌ NO toast notifications (v1.1)
```

**Rationale:**
- Proves concept faster
- Can add polish later
- Still looks professional

---

### Decision #8: Charts - Line Charts Only for MVP

**Decision:** Simple line charts for MVP, add gauges in v1.1

**MVP Charts:**
```javascript
// Simple Recharts LineChart
<LineChart data={metrics}>
  <Line type="monotone" dataKey="cpu" stroke="#8884d8" />
  <XAxis dataKey="timestamp" />
  <YAxis />
</LineChart>
```

**Future (v1.1):**
- Animated circular gauges
- Fancy progress indicators
- Glow effects

---

## 🔧 Technical Decisions

### Decision #9: Metrics - Simple First

**Decision:** Start with total percentages, add details in v1.1

**MVP Metrics:**
```go
type Metrics struct {
    Timestamp time.Time
    CPU       float64  // Total CPU %
    Memory    float64  // Memory %
    DiskIO    float64  // Disk operations/sec
    Network   float64  // Network MB/s
}
```

**Future (v1.1):**
```go
type DetailedMetrics struct {
    CPU struct {
        Total   float64
        PerCore []float64  // ← Add later
    }
    Memory struct {
        Percent float64
        Used    uint64     // ← Add later
        Total   uint64     // ← Add later
    }
}
```

---

### Decision #10: Docker - Add in Week 3

**Decision:** Develop locally without Docker, add Docker in Week 3

**Rationale:**
- ✅ Faster iteration during development
- ✅ Easier debugging
- ✅ Hot reload works better
- ✅ Add Docker once code works

**Week 1-2:**
```bash
# Run locally
cd backend && go run cmd/server/main.go
cd frontend && npm run dev
```

**Week 3:**
```bash
# Create Dockerfiles
# Test docker-compose
# Verify deployment
```

---

## 🚀 Deployment Decisions

### Decision #11: Showcase Deployment - Railway.app

**Decision:** Use Railway.app for v1.2 showcase deployment

**Rationale:**
- ✅ Easiest Docker Compose support
- ✅ Free trial available
- ✅ Auto-deploy from GitHub
- ✅ Built-in SSL/HTTPS
- ✅ Simple setup (~30 minutes)

**Cost:** $5-10/month (or free with trial)

**Alternatives considered:**
- Render.com (free tier but sleeps)
- Fly.io (more complex setup)

---

### Decision #12: AWS - Optional After Showcase

**Decision:** AWS deployment is OPTIONAL, not required for success

**Path A (Recommended):**
```
v1.0 Local → v1.1 Polish → v1.2 Showcase → DONE
Cost: ~$10/month
Good for: Portfolio, job search
```

**Path B (If needed):**
```
v1.0 Local → v1.2 Showcase → v2.0 AWS Dev → v2.1 AWS Prod
Cost: ~$150-200/month
Good for: AWS experience, DevOps roles
```

---

## 📚 Documentation Decisions

### Decision #13: Single Source of Truth

**Decision:** MVP-PLAN.md is the ONLY source for "what to build"

**Active Documents:**
```
✅ MVP-PLAN.md              ← What to build (detailed)
✅ AI-DEVELOPMENT-RULES.md  ← How to build (rules)
✅ MASTER-PLAN.md           ← Overview & navigation
✅ VERSION-ROADMAP.md       ← Version strategy
✅ DECISIONS.md (this file) ← Key decisions
✅ RELEASE-PLAN.md          ← Deployment only (Phase 1.5+)
```

**Archived:**
```
📦 system-monitor-dashboard-project-plan.md
   - Moved to /docs/archive/
   - Keep for reference only
   - DO NOT USE as primary source
```

---

## 🔒 Safety Decisions

### Decision #14: Safety Limits - Strict and Enforced

**Decision:** Hardware protection is MANDATORY with these exact limits

**Limits:**
```go
const (
    // Local development
    MAX_CPU_PERCENT      = 95
    MAX_CPU_DURATION     = 30  // seconds
    MAX_MEMORY_PERCENT   = 25  // of total RAM
    MAX_MEMORY_DURATION  = 60  // seconds
    MAX_DISK_SIZE_MB     = 100
    MAX_CONCURRENT       = 5

    // Emergency shutdown
    CRITICAL_CPU         = 98
    CRITICAL_MEMORY      = 95
)
```

**Showcase (public demo) - Stricter:**
```go
const (
    MAX_CPU_PERCENT      = 80  // Lower
    MAX_CPU_DURATION     = 20  // Shorter
    MAX_MEMORY_PERCENT   = 20  // Lower
    MAX_CONCURRENT       = 3   // Fewer
)
```

**Testing:** All safety limits MUST have 100% test coverage

---

### Decision #15: Action Cancellation - Required

**Decision:** ALL actions MUST be cancellable within 1 second

**Implementation:**
```go
func (a *Action) Execute(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            cleanup()
            return ctx.Err()
        default:
            // Do work
        }
    }
}
```

**Test:** Every action must have test proving it stops within 1s

---

## 🎯 Scope Decisions

### Decision #16: MVP Scope - Locked

**Decision:** This is MVP scope. Do NOT add features.

**IN SCOPE ✅:**
```
Functionality:
✅ 4 metrics (CPU, Memory, Disk, Network)
✅ 4 actions (stress each metric)
✅ Real-time updates (1-second polling)
✅ Simple line charts
✅ Basic event log
✅ Safety limits
✅ Docker setup

Quality:
✅ 70% test coverage minimum
✅ Safety tests at 100%
✅ Works in Chrome/Firefox
✅ Mobile responsive (basic)
```

**OUT OF SCOPE ❌:**
```
❌ WebSocket (use polling)
❌ Database (in-memory only)
❌ Pre-built scenarios
❌ Animated gauges
❌ Active actions widget
❌ User authentication
❌ Custom scenario builder
❌ Historical data export
❌ Multi-instance monitoring
```

**If Tempted to Add:**
1. Stop
2. Add to FUTURE-FEATURES.md
3. Continue with MVP
4. Build it AFTER v1.0 complete

---

## 🌐 Browser Support Decisions

### Decision #17: Browser Support - Modern Only

**Decision:** Support Chrome/Firefox latest, test Safari

**Supported:**
```
✅ Chrome latest (primary target)
✅ Firefox latest (primary target)
✅ Edge latest (Chromium-based)
✅ Safari latest (test, but not blocker)

❌ IE11 (not supported)
❌ Old browsers (not supported)
```

**Mobile:**
```
✅ Basic responsive (works on mobile)
❌ Mobile-optimized UI (defer to v1.1)
❌ Touch gestures (defer to v1.1)
```

---

## 📊 Performance Decisions

### Decision #18: Performance Targets

**Decision:** These are minimum performance requirements

**Response Times:**
```
Button click → Visual feedback:  <200ms
API request → Response:          <500ms
Action start → Metric visible:   <2s
Chart update:                    <100ms
Initial page load:               <3s
```

**Action Effectiveness:**
```
CPU Stress:    Must reach 85-95% CPU
Memory Surge:  Must allocate 90-100% of requested memory
Disk Storm:    Must complete all operations
Traffic Flood: Must generate target requests/sec
```

**If Not Meeting:**
- Debug and optimize first
- If still failing, adjust targets
- Document actual performance

---

## 🎯 Milestone Decisions

### Decision #19: Go/No-Go Decision Points

**Decision:** Check these milestones, stop if failing

**After Week 1:**
```
✅ GO if: Backend collects metrics, CPU action works
❌ NO-GO if: Can't collect metrics or safety limits fail
Action: If NO-GO, get human help, reassess approach
```

**After Week 2:**
```
✅ GO if: Frontend shows metrics, one button works end-to-end
❌ NO-GO if: Frontend can't connect or actions don't work
Action: If NO-GO, debug for 2 days, then get human help
```

**After Week 4:**
```
✅ GO if: All 4 actions work, can do 5-min demo
❌ NO-GO if: Major features broken or crashes
Action: If NO-GO, extend timeline or reduce scope
```

---

## 💰 Cost Decisions

### Decision #20: Budget Awareness

**Decision:** Clear cost expectations at each stage

**Costs:**
```
v1.0 MVP (local):      $0/month
v1.1 Polish (local):   $0/month
v1.2 Showcase:         $5-10/month
v1.3 Scenarios:        $5-10/month
v1.4 Persistence:      $20-25/month
v2.0 AWS Dev:          $75-100/month
v2.1 AWS Prod:         $150-200/month
```

**Recommendation:**
- Start at v1.0 (free)
- Go to v1.2 for portfolio (~$10/month)
- Only go to AWS if budget allows and needed

---

## 🔄 Change Process

**These decisions are LOCKED.**

**To change a decision:**
1. Stop development
2. Discuss with human
3. Update this file
4. Update all affected documents
5. Communicate to AI
6. Resume development

**Do NOT:**
- Change decisions during development
- Make exceptions "just this once"
- Skip updating documentation

---

## 📋 Quick Reference

**For AI starting development:**
- WebSocket? → NO, use polling
- Database? → NO, in-memory only
- Scenarios? → NO, not in MVP
- Docker from start? → NO, add Week 3
- Gauges? → NO, simple charts only
- Timeline? → 4-6 weeks realistic
- Coverage? → 70% minimum

**For Human reviewing progress:**
- Is AI building WebSocket? → STOP, wrong!
- Is AI setting up database? → STOP, wrong!
- Is AI adding scenarios? → STOP, not in MVP!
- Refer back to this document

---

**Version:** 1.0 (Final)
**Last Updated:** 2025-01-09
**Status:** LOCKED - Single source of truth for all key decisions
