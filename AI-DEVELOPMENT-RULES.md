# AI Development Rules & Guidelines
## Senior Developer Instructions for AI Agents

**Project:** Interactive System Monitoring Dashboard
**Purpose:** This document contains MANDATORY rules that AI agents MUST follow when building this project
**Status:** ACTIVE - Must be read before every development session

---

## ⚠️ CRITICAL: Read This First

**Before writing ANY code:**

1. **Read DECISIONS.md** - Contains final decisions on all contradictions
   - WebSocket? → NO (use polling)
   - Database? → NO (in-memory only)
   - Scenarios? → NO (not in MVP)
   - Folder structure? → internal/actions/
   - Timeline? → 4-6 weeks realistic

2. **Read MVP-PLAN.md** - Single source of truth for what to build
   - Week-by-week plan
   - Clear success criteria
   - What's in/out of scope

3. **Read this file** - Rules for HOW to build

**If you see ANY contradiction:**
- STOP
- Check DECISIONS.md first
- If still unclear, ask human
- Do NOT make assumptions

---

## 🎯 Core Philosophy

```
"Build it to work, not to impress with complexity"
"Safety first, features second"
"User can always see what's happening"
"Every action must be reversible"
"The demo must never crash"
```

---

## 🚨 CRITICAL SAFETY RULES (NEVER VIOLATE)

### Rule 1: Hardware Protection is MANDATORY
```go
// THESE LIMITS ARE NON-NEGOTIABLE
const (
    MAX_CPU_PERCENT      = 95    // NEVER exceed 95% CPU
    MAX_CPU_DURATION     = 30    // Max 30 seconds per action
    MAX_MEMORY_PERCENT   = 25    // Max 25% of total system RAM
    MAX_MEMORY_DURATION  = 60    // Max 60 seconds
    MAX_DISK_SIZE_MB     = 100   // Max 100MB temp files
    MAX_CONCURRENT       = 5     // Max 5 actions simultaneously

    // Emergency shutdown thresholds
    CRITICAL_CPU         = 98    // Kill action immediately
    CRITICAL_MEMORY      = 95    // Kill action immediately
)
```

**WHY:** Without these, you WILL crash the system during demos.

**IMPLEMENTATION:**
- ✅ Check limits BEFORE starting any action
- ✅ Monitor limits DURING action execution
- ✅ Emergency shutdown if thresholds exceeded
- ✅ Cleanup resources in defer statements
- ✅ Context-based cancellation for all actions

### Rule 2: All Actions Must Be Cancellable
```go
// EVERY action must accept context
func (a *Action) Execute(ctx context.Context) error {
    // MANDATORY: Check context in loops
    for {
        select {
        case <-ctx.Done():
            // MANDATORY: Cleanup before return
            defer cleanup()
            return ctx.Err()
        default:
            // Do work
        }
    }
}
```

### Rule 3: Resource Cleanup is MANDATORY
```go
// ALWAYS use defer for cleanup
func doWork() error {
    resource := allocate()
    defer resource.Release()  // MANDATORY

    tempFile := createTemp()
    defer os.Remove(tempFile)  // MANDATORY

    // ... work ...
}
```

### Rule 4: Never Block the Main Thread
```go
// WRONG - blocks UI
func handleAction(w http.ResponseWriter, r *http.Request) {
    runExpensiveOperation()  // BAD!
    json.NewEncoder(w).Encode(result)
}

// RIGHT - async execution
func handleAction(w http.ResponseWriter, r *http.Request) {
    actionID := generateID()
    go engine.Execute(actionID, action)  // Good!
    json.NewEncoder(w).Encode(Response{ID: actionID})
}
```

---

## 📐 Architecture Rules

### Rule 5: Follow the Established Structure
```
backend/
  cmd/server/main.go           ← Entry point ONLY
  internal/
    metrics/                   ← Metrics collection ONLY
    actions/                   ← Load generation ONLY
    api/                       ← HTTP handlers ONLY
    websocket/                 ← WebSocket ONLY (if implemented)
  pkg/models/                  ← Shared types ONLY

DO NOT:
- Mix concerns (metrics in actions, actions in API)
- Create new top-level directories without approval
- Put business logic in main.go
```

### Rule 6: API Design Standards
```http
# Good API design
POST /api/actions/cpu-stress
GET  /api/actions/active
DELETE /api/actions/:id/stop

# Bad API design
POST /api/cpu              ← Too vague
GET  /api/getActions       ← Redundant "get"
POST /api/stopCPU/:id      ← Use DELETE, not POST
```

**Rules:**
- ✅ RESTful conventions
- ✅ Plural nouns for collections
- ✅ Kebab-case for URLs
- ✅ JSON for request/response
- ✅ Proper HTTP status codes

### Rule 7: Error Handling Standards
```go
// MANDATORY error handling pattern
func doSomething() error {
    result, err := operation()
    if err != nil {
        // Log the error
        log.Printf("operation failed: %v", err)
        // Return wrapped error with context
        return fmt.Errorf("doSomething: operation failed: %w", err)
    }
    return nil
}

// HTTP handler error pattern
func handler(w http.ResponseWriter, r *http.Request) {
    if err := validate(r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := process()
    if err != nil {
        log.Printf("process error: %v", err)
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(result)
}
```

---

## 🧪 Testing Rules

### Rule 8: Test Coverage Requirements
```
MINIMUM test coverage: 70% for MVP
TARGET test coverage: 80% for production

MANDATORY tests:
- ✅ All safety limits enforced
- ✅ All actions can be cancelled
- ✅ Resource cleanup works
- ✅ Concurrent actions don't crash
- ✅ API endpoints return correct codes
- ✅ Invalid inputs are rejected
```

### Rule 9: Test Organization
```go
// Test file naming
metrics_test.go              ← Tests for metrics.go
cpu_stress_test.go          ← Tests for cpu_stress.go

// Test function naming
func TestCPUMetricsCollected(t *testing.T)           ← Good
func TestCPUStressIncreasesMetrics(t *testing.T)     ← Good
func TestSafetyLimitsEnforced(t *testing.T)          ← CRITICAL
func TestActionStopsGracefully(t *testing.T)         ← CRITICAL

func TestStuff(t *testing.T)                         ← Bad - too vague
```

### Rule 10: CRITICAL Tests Must Pass
```go
// These tests MUST ALWAYS PASS before merging
var criticalTests = []string{
    "TestSafetyLimitsEnforced",
    "TestActionStopsGracefully",
    "TestEmergencyShutdown",
    "TestResourceCleanup",
    "TestConcurrentActionsSafe",
}
```

---

## 🎨 Frontend Rules

### Rule 11: User Feedback is MANDATORY
```javascript
// WRONG - no feedback
function handleClick() {
    fetch('/api/actions/cpu-stress', {method: 'POST'});
}

// RIGHT - immediate feedback
async function handleClick() {
    // 1. Show loading
    setLoading(true);

    try {
        // 2. Make request
        const response = await fetch('/api/actions/cpu-stress', {
            method: 'POST'
        });

        if (!response.ok) {
            throw new Error('Action failed');
        }

        // 3. Show success
        toast.success('CPU stress started!');
    } catch (error) {
        // 4. Show error
        toast.error('Failed to start action');
        console.error(error);
    } finally {
        // 5. Hide loading
        setLoading(false);
    }
}
```

### Rule 12: Performance Standards
```javascript
// MANDATORY performance targets
const PERFORMANCE_TARGETS = {
    buttonClick: 200,        // Button response < 200ms
    metricUpdate: 1000,      // Update interval: 1s
    chartRender: 16,         // 60 FPS = 16ms per frame
    apiTimeout: 5000,        // API timeout: 5s
};

// Monitor performance
function measurePerformance(name, fn) {
    const start = performance.now();
    fn();
    const duration = performance.now() - start;

    if (duration > PERFORMANCE_TARGETS[name]) {
        console.warn(`${name} too slow: ${duration}ms`);
    }
}
```

### Rule 13: Accessibility Requirements
```jsx
// MANDATORY accessibility
<button
    onClick={handleClick}
    aria-label="Start CPU stress test"     // REQUIRED
    disabled={isLoading}                   // REQUIRED
    className={`btn ${isLoading ? 'btn-loading' : ''}`}
>
    {isLoading ? 'Starting...' : '🔥 CPU Stress'}
</button>

// WRONG - no aria-label, no disabled state
<button onClick={handleClick}>
    🔥 CPU
</button>
```

---

## 📝 Code Quality Rules

### Rule 14: Naming Conventions
```go
// Go Backend
type CPUStressAction struct {}      // PascalCase for exported
func (a *Action) Execute() {}       // PascalCase for exported
var maxCPUPercent = 95              // camelCase for private

// JavaScript Frontend
const API_BASE_URL = '...';         // UPPER_CASE for constants
function handleCPUStress() {}       // camelCase for functions
const MetricCard = () => {}         // PascalCase for components
```

### Rule 15: Comments & Documentation
```go
// MANDATORY: Public API must have comments
// CPUStressAction generates CPU load for testing.
// It uses busy loops in goroutines to consume CPU cycles.
//
// Safety: Respects MAX_CPU_PERCENT and MAX_CPU_DURATION limits.
// Cancellation: Responds to context cancellation within 100ms.
type CPUStressAction struct {
    TargetPercent int           // Target CPU utilization (0-95)
    Duration      time.Duration // How long to run (max 30s)
}

// WRONG - no documentation
type CPUStressAction struct {
    TargetPercent int
    Duration      time.Duration
}
```

### Rule 16: Magic Numbers are FORBIDDEN
```go
// WRONG - magic numbers
if cpu > 95 {
    return errors.New("too high")
}

// RIGHT - named constants
const MAX_CPU_PERCENT = 95

if cpu > MAX_CPU_PERCENT {
    return fmt.Errorf("CPU %d%% exceeds limit %d%%", cpu, MAX_CPU_PERCENT)
}
```

---

## 🚀 Development Workflow Rules

### Rule 17: Commit Standards
```bash
# Good commit messages
git commit -m "feat: Add CPU stress action with safety limits"
git commit -m "fix: Prevent memory leak in action cleanup"
git commit -m "test: Add emergency shutdown tests"
git commit -m "docs: Update API documentation"

# Bad commit messages
git commit -m "updates"
git commit -m "fix bug"
git commit -m "wip"
```

**Format:**
```
<type>: <description>

Types:
- feat: New feature
- fix: Bug fix
- test: Add/update tests
- docs: Documentation
- refactor: Code restructuring
- perf: Performance improvement
```

### Rule 18: Incremental Development
```
DO:
✅ Build one feature completely before starting next
✅ Test each feature before moving on
✅ Commit working code frequently
✅ Ask for review at checkpoints

DON'T:
❌ Start multiple features simultaneously
❌ Commit broken code
❌ Skip testing "to save time"
❌ Make massive changes without checkpoints
```

### Rule 19: When to Ask for Human Help
```
🔴 STOP and Ask Human:
- Safety limits don't prevent system crash
- Architecture decision needed
- Performance is terrible (>2x targets)
- Security concern discovered
- Breaking changes required

🟡 Ask When Convenient:
- UI/UX design decisions
- Feature prioritization
- Color scheme choices

🟢 AI Can Decide:
- Implementation details
- File organization
- Variable names
- Code formatting
- Test cases
```

---

## 🎯 MVP-Specific Rules

### Rule 20: MVP Scope is SACRED
```
✅ IN SCOPE (Build These):
- 4 metric types (CPU, Memory, Disk, Network)
- 4 load actions (one per metric)
- Simple dashboard with charts
- Action buttons + event log
- Docker setup for local dev
- Basic tests (>70% coverage)

❌ OUT OF SCOPE (Do NOT Build):
- Database persistence
- WebSocket (use polling)
- User authentication
- Cloud deployment
- Pre-built scenarios
- Custom scenario builder
- Mobile app
- Advanced visualizations

IF TEMPTED TO ADD FEATURE:
1. Check if it's in MVP scope
2. If NO → Add to FUTURE-FEATURES.md
3. Keep building MVP
```

### Rule 21: MVP Progress Tracking
```
BEFORE starting work:
1. Read MVP-PLAN.md
2. Identify current phase
3. List tasks for this session
4. Set clear checkpoint

DURING work:
1. Check off completed tasks
2. Note blockers immediately
3. Update progress

AFTER work:
1. Test what was built
2. Document what works
3. Note what's next
4. Report progress
```

---

## 🔒 Security Rules

### Rule 22: Input Validation is MANDATORY
```go
// MANDATORY validation for all inputs
func validateCPUStressRequest(req CPUStressRequest) error {
    if req.TargetPercent < 0 || req.TargetPercent > MAX_CPU_PERCENT {
        return fmt.Errorf("target_percent must be 0-%d", MAX_CPU_PERCENT)
    }

    if req.Duration < 0 || req.Duration > MAX_CPU_DURATION {
        return fmt.Errorf("duration must be 0-%d seconds", MAX_CPU_DURATION)
    }

    return nil
}
```

### Rule 23: No Secrets in Code
```bash
# WRONG - secret in code
DB_PASSWORD=supersecret123

# RIGHT - use environment variables
DB_PASSWORD=${DB_PASSWORD}

# WRONG - commit .env file
git add .env

# RIGHT - .env in .gitignore
echo ".env" >> .gitignore
```

### Rule 24: CORS Configuration
```go
// MVP: Allow localhost only
func configureCORS() {
    allowedOrigins := []string{
        "http://localhost:3000",
        "http://localhost:5173",  // Vite dev server
    }

    // Production: Specific domain only
    if isProduction {
        allowedOrigins = []string{
            "https://yourdomain.com",
        }
    }
}
```

---

## 📊 Monitoring & Observability Rules

### Rule 25: Logging Standards
```go
// MANDATORY log levels
log.Info("Action started", "id", actionID, "type", actionType)
log.Warn("CPU approaching limit", "current", cpu, "limit", MAX_CPU)
log.Error("Action failed", "error", err, "id", actionID)

// WRONG - no context
log.Println("something happened")

// WRONG - too verbose
log.Println("Entering function doWork()")
```

### Rule 26: Metrics That Matter
```go
// MUST track these metrics
var requiredMetrics = []string{
    "active_actions_count",     // How many actions running
    "action_duration_seconds",  // How long actions take
    "action_failures_total",    // How many failures
    "cpu_usage_percent",        // Current CPU
    "memory_usage_percent",     // Current memory
}
```

---

## 🎨 UI/UX Rules

### Rule 27: Visual Feedback Timeline
```
User Action → Visual Response

0ms:     User clicks button
0-50ms:  Button press animation starts
50-200ms: Loading state shown
200ms+:   API call in flight
Response: Update UI immediately
          Show success/error toast
          Update metrics

RULE: User MUST see feedback within 200ms
```

### Rule 28: Error Messages Must Be Helpful
```javascript
// WRONG - useless error
toast.error("Error");

// WRONG - technical jargon
toast.error("HTTP 500: Internal Server Error");

// RIGHT - actionable message
toast.error("Could not start CPU stress. Please try again.");

// BETTER - with help
toast.error(
    "Could not start CPU stress. " +
    "System may already be under load. " +
    "Stop other actions and try again."
);
```

### Rule 29: Loading States Are MANDATORY
```jsx
// WRONG - no loading state
<button onClick={handleClick}>Start</button>

// RIGHT - clear loading state
<button
    onClick={handleClick}
    disabled={isLoading}
>
    {isLoading ? (
        <>
            <Spinner /> Starting...
        </>
    ) : (
        '🔥 Start CPU Stress'
    )}
</button>
```

---

## 🎯 AI Agent-Specific Instructions

### Rule 30: Context Awareness
```
BEFORE writing ANY code:
1. Read relevant existing code
2. Understand the current architecture
3. Check existing patterns
4. Maintain consistency

DON'T:
- Assume you know the structure
- Introduce new patterns without reason
- Ignore existing conventions
```

### Rule 31: Explain Your Decisions
```
WHEN making architectural decisions:
1. Explain WHY you chose this approach
2. List alternatives considered
3. Explain trade-offs
4. Document assumptions

EXAMPLE:
"I'm using polling (1s interval) instead of WebSocket because:
- MVP requirement: Keep it simple
- 1s latency is acceptable for demo
- Easier to debug
- Can upgrade to WebSocket later
Trade-off: Slightly higher server load vs simpler code"
```

### Rule 32: Progressive Disclosure
```
WHEN explaining to user:
1. Start with high-level summary
2. Provide details if asked
3. Show code examples when relevant
4. Offer to explain further

DON'T:
- Dump all information at once
- Use jargon without explanation
- Assume user knowledge level
```

### Rule 33: Checkpoint Protocol
```
AFTER completing major task:
1. Summarize what was built
2. What works (with proof)
3. What's next
4. Any blockers or concerns

EXAMPLE:
"✅ Completed: CPU stress action with safety limits

 What works:
 - Can trigger via API: POST /api/actions/cpu-stress
 - CPU increases to target %
 - Stops after duration
 - Emergency shutdown at 98%
 - Tests pass (see test output)

 What's next:
 - Memory surge action
 - Then disk and network

 No blockers."
```

---

## 📋 Pre-Development Checklist

Before starting ANY development session:

- [ ] Read these rules
- [ ] Read MVP-PLAN.md
- [ ] Understand current phase
- [ ] Know what to build
- [ ] Know success criteria
- [ ] Know safety requirements
- [ ] Know when to ask for help

---

## 🚨 Red Flags (Stop If You See These)

```
🚩 System crashed during testing
🚩 Can't stop an action once started
🚩 Memory keeps growing
🚩 Tests are failing
🚩 No safety limits implemented
🚩 API returns 500 errors
🚩 Frontend freezes during action
🚩 Breaking changes required
🚩 Security vulnerability found

ACTION: Stop, report to human, wait for guidance
```

---

## ✅ Quality Gates

### Before Committing Code:
- [ ] Code follows project structure
- [ ] Safety limits implemented
- [ ] Error handling present
- [ ] Tests written and passing
- [ ] No console.log in production
- [ ] No commented-out code
- [ ] No TODOs without issue number

### Before Marking Task Complete:
- [ ] Feature works as specified
- [ ] Tests prove it works
- [ ] Documentation updated
- [ ] No known bugs
- [ ] Performance acceptable
- [ ] User feedback implemented

### Before Moving to Next Phase:
- [ ] All phase tasks complete
- [ ] All tests passing
- [ ] No critical bugs
- [ ] Human approval obtained
- [ ] Progress documented

---

## 🎓 Learning From Mistakes

If something goes wrong:

1. **Document it**
   - What happened
   - What was expected
   - What went wrong
   - How it was fixed

2. **Add safeguards**
   - Add test to prevent recurrence
   - Add validation if missing
   - Update documentation

3. **Share learning**
   - Update this document if needed
   - Add to troubleshooting guide

---

## 📚 Reference Quick Links

**Before coding:**
- Read: MVP-PLAN.md → Know what to build
- Read: These rules → Know how to build
- Read: Existing code → Maintain consistency

**During coding:**
- Check: Safety limits enforced?
- Check: Error handling present?
- Check: Tests written?
- Check: User feedback shown?

**After coding:**
- Test: Does it work?
- Test: Can it be stopped?
- Test: Does it handle errors?
- Document: What was built

---

## 🎯 Success Metrics

**Code Quality:**
- Zero crashes during demo
- All safety tests pass
- >70% test coverage
- All API errors handled
- All user actions have feedback

**User Experience:**
- Button response <200ms
- Metrics update every 1s
- Charts animate smoothly
- Errors are helpful
- Loading states clear

**Development Process:**
- Checkpoints every task
- Working code committed
- Tests before features
- Documentation current

---

## 🔄 Version

**Version:** 1.0
**Last Updated:** 2025-01-09
**Status:** Active
**Review:** After MVP complete

---

**REMEMBER: These rules exist to prevent problems, not slow you down. Follow them, and the MVP will succeed. Ignore them, and the demo will crash.**

**When in doubt: ASK. Better to ask than to break things.**
