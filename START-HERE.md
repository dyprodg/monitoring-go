# START HERE - Quick Start Guide
## Interactive System Monitoring Dashboard

**Welcome!** This guide helps you get started quickly.

---

## 🎯 If You're the Human (Project Owner)

### First Time Here?

**Read these in order:**

1. **MASTER-PLAN.md** (5 min)
   - Overview of the entire project
   - What we're building and why
   - Success criteria

2. **DECISIONS.md** (10 min)
   - All key architectural decisions locked in
   - WebSocket? NO (polling)
   - Database? NO (in-memory)
   - Timeline? 4-6 weeks

3. **MVP-PLAN.md** (15 min)
   - Detailed 4-6 week implementation plan
   - Week-by-week breakdown
   - What AI will build

### Quick Summary

**What:** Interactive system monitoring dashboard where you click buttons to generate load and watch metrics spike in real-time.

**Timeline:** 4-6 weeks for working local demo

**Cost:** $0 for MVP (local only)

**Next Steps:**
1. Review the key documents above
2. When ready, tell AI: "Start building - begin with Week 1, Day 1"
3. Check progress every 2 days
4. Make go/no-go decisions at weekly checkpoints

---

## 🤖 If You're the AI (Development Agent)

### CRITICAL: Read Before Writing ANY Code

**Step 1: Read DECISIONS.md** (MANDATORY)
- Contains final decisions on all contradictions
- WebSocket? → NO (use HTTP polling)
- Database? → NO (in-memory only)
- Scenarios? → NO (not in MVP)
- Folder: `internal/actions/` (not loadgen)

**Step 2: Read AI-DEVELOPMENT-RULES.md** (MANDATORY)
- 33 mandatory rules
- Safety requirements (NEVER exceed limits)
- Code standards
- When to ask for help

**Step 3: Read MVP-PLAN.md** (MANDATORY)
- Week-by-week implementation plan
- Current phase and tasks
- Success criteria

### Quick Reference Card

```
WebSocket?          → NO (use polling every 1s)
Database?           → NO (in-memory only, 60s history)
Scenarios?          → NO (not in MVP, defer to v1.3)
Docker from Day 1?  → NO (add in Week 3)
Gauges?             → NO (simple line charts only)
Tests?              → YES (70% minimum, safety 100%)
Folder structure?   → internal/actions/ (not loadgen)
Timeline?           → 4-6 weeks realistic

Safety Limits:
MAX_CPU_PERCENT     = 95
CRITICAL_CPU        = 98 (emergency shutdown)
MAX_MEMORY_PERCENT  = 25
CRITICAL_MEMORY     = 95 (emergency shutdown)

All actions MUST:
- Accept context.Context
- Cancel within 1 second
- Cleanup resources in defer
- Have 100% test coverage for safety
```

### Starting Development

**When human says "start building":**

```bash
# Week 1, Day 1 Tasks:
1. Initialize Go project
   mkdir -p backend/cmd/server
   cd backend
   go mod init monitoring-dashboard

2. Create main.go
   package main with basic HTTP server

3. Implement CPU metrics collector
   Use gopsutil package
   Collect total CPU percentage

4. Write tests
   Test CPU metrics collection

5. Create API endpoint
   GET /api/metrics

6. Checkpoint with human
   "✅ Backend collects CPU metrics
    ✅ API endpoint works
    ✅ Tests pass
    Ready for Day 3?"
```

**IMPORTANT:**
- Report progress after each task
- Ask human if blocked >4 hours
- NEVER build WebSocket (use polling)
- NEVER setup database (in-memory only)

---

## 📁 Document Navigation

### Planning Documents (Read First)
```
MASTER-PLAN.md         → Overview & navigation
DECISIONS.md           → Key decisions (LOCKED)
MVP-PLAN.md            → What to build (4-6 weeks)
AI-DEVELOPMENT-RULES.md → How to build (mandatory rules)
```

### Reference Documents
```
VERSION-ROADMAP.md     → Version strategy (v1.0 → v3.0)
RELEASE-PLAN.md        → Deployment options
FEATURE-SUMMARY.md     → Visual examples
FUTURE-FEATURES.md     → Post-MVP ideas
```

### Generated During Project
```
README.md              → Project documentation
tests/                 → Test files
docs/                  → Additional docs
```

---

## ✅ Pre-Flight Checklist

### Human Checklist
- [ ] Read MASTER-PLAN.md
- [ ] Read DECISIONS.md
- [ ] Read MVP-PLAN.md
- [ ] Understand success criteria (working demo in 4-6 weeks)
- [x] Know cost ($0 for MVP, ~$10/month for showcase)
- [x] Have 4-6 weeks available (part-time OK)
- [ ] Ready to review progress every 2 days
- [x] Confirmed: Go 1.21+, Node 18+, Docker installed

### AI Checklist
- [ ] Read DECISIONS.md (no contradictions!)
- [ ] Read AI-DEVELOPMENT-RULES.md (33 mandatory rules)
- [ ] Read MVP-PLAN.md Week 1 tasks
- [ ] Understand: NO WebSocket, NO database in MVP
- [ ] Know safety limits (95% CPU max, 25% memory max)
- [ ] Know when to ask for help (stuck >4 hours)
- [ ] Ready to report progress frequently

---

## 🚀 Quick Start Commands

### Starting Development (Local)

```bash
# Backend
cd backend
go run cmd/server/main.go

# Frontend (in separate terminal)
cd frontend
npm install
npm run dev

# Access
# Backend:  http://localhost:8080
# Frontend: http://localhost:3000
```

### Starting with Docker (Week 3+)

```bash
docker-compose up

# Access
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
```

---

## 🎯 Success Criteria (Simple)

**MVP is successful when:**
```
1. Run: docker-compose up
2. Open: http://localhost:3000
3. Click: "🔥 CPU Stress" button
4. See: CPU graph spike to ~90%
5. Wait: Action completes after 10 seconds
6. Verify: No crash, system returns to normal

= SUCCESS! You have a working MVP! 🎉
```

---

## 📞 When to Ask for Help

### Human Should Ask:
- Confused about what to do next? → Read MASTER-PLAN.md
- AI is stuck or confused? → Check DECISIONS.md, clarify
- Progress seems slow? → Review MVP-PLAN.md timeline
- Want to add a feature? → Check if it's in FUTURE-FEATURES.md

### AI Should Ask:
- 🔴 STOP IMMEDIATELY if:
  - Safety tests are failing
  - System crashed during testing
  - Can't stop an action once started
  - Architecture decision needed

- 🟡 Ask when convenient:
  - UI/UX design choices
  - Color scheme decisions

- 🟢 AI can decide:
  - Implementation details
  - Variable names
  - Code organization

---

## 🎯 Common Questions

**Q: How long will this take?**
A: 4-6 weeks for MVP, depending on available time.

**Q: Do I need to know Go/React?**
A: AI will write the code, but reviewing helps.

**Q: What if I get stuck?**
A: AI will ask for help. Check DECISIONS.md first.

**Q: Can I change the plan?**
A: Yes, but update DECISIONS.md first.

**Q: Do I need AWS?**
A: No! MVP is local only. AWS is optional later.

**Q: What if MVP fails?**
A: That's OK! You'll learn why. Can pivot or stop.

---

## 📋 Weekly Checkpoints

### Week 1 Checkpoint
**Expected:** Backend collects metrics, CPU action works
**Demo:** `curl http://localhost:8080/api/metrics`
**Decision:** GO if working, NO-GO if major issues

### Week 2 Checkpoint
**Expected:** Frontend shows metrics, one button works
**Demo:** Click CPU button, see graph spike
**Decision:** GO if working, get help if broken

### Week 3 Checkpoint
**Expected:** All 4 actions work
**Demo:** All buttons trigger actions safely
**Decision:** GO to Week 4 polish

### Week 4 Checkpoint
**Expected:** MVP complete, tests pass
**Demo:** Full 5-minute demo
**Decision:** MVP DONE → Choose v1.1 or v1.2

---

## 🎯 What Success Looks Like

**After 4-6 Weeks:**
- ✅ Working demo on your laptop
- ✅ Click button → metrics spike
- ✅ All 4 actions functional
- ✅ Tests passing (>70%)
- ✅ Can do impressive 5-min demo
- ✅ Code on GitHub
- ✅ Docker setup working

**Next Options:**
1. **Stop here** - You have a working demo
2. **v1.1 Polish** - Make it look better (+1 week)
3. **v1.2 Showcase** - Put it online (+3-5 days)
4. **v2.0 AWS** - Add cloud deployment (+6-8 weeks)

---

## 🚀 Ready to Start?

### For Human:
1. Confirm you've read the key docs
2. Tell AI: "Start building MVP - begin Week 1, Day 1"
3. Review progress after Day 2
4. Make go/no-go decisions at each weekly checkpoint

### For AI:
1. Read DECISIONS.md
2. Read AI-DEVELOPMENT-RULES.md
3. Read MVP-PLAN.md Week 1
4. Start with Week 1, Day 1, Task 1
5. Report progress frequently

---

**CURRENT STATUS:** Ready to Start
**NEXT ACTION:** Human confirms, AI begins Week 1, Day 1
**FIRST CHECKPOINT:** After Day 2 (backend metrics working)

**Let's build something impressive! 🚀**
