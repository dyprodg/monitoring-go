# System Monitor Dashboard - Release Plan

## 📋 Release Strategy

**Philosophy:** Build locally first, deploy to cloud incrementally
- **Phase 1:** Local MVP with Docker Compose
- **Phase 1.5:** Showcase Deployment (Simple Public Demo) ⭐ NEW
- **Phase 2:** AWS Dev Environment
- **Phase 3:** AWS Production Environment

**Repository:** GitHub (version control only, no CI/CD initially)

---

## 🌐 Deployment Paths

You have TWO paths after completing the local MVP:

### Path A: Quick Showcase (Recommended for Portfolio)
```
Phase 1 (Local MVP) → Phase 1.5 (Showcase) → DONE
Cost: ~$10/month
Time: +2-3 days
Result: Public URL you can share
```

### Path B: Full Production (For Enterprise Portfolio)
```
Phase 1 (Local MVP) → Phase 2 (AWS Dev) → Phase 3 (AWS Prod)
Cost: ~$150-200/month
Time: +6-8 weeks
Result: Production-grade AWS deployment
```

Choose Path A if you want to showcase it online quickly and cheaply.
Choose Path B if you want to demonstrate AWS/production expertise.

---

## 🎯 Phase 1: Local MVP (Docker Compose)

**Goal:** Fully functional system running locally with docker-compose
**Timeline:** 4-6 weeks
**Deliverable:** Demo-ready application

**📋 For detailed implementation plan, see: [MVP-PLAN.md](MVP-PLAN.md)**

### Phase 1 Overview

Phase 1 builds the complete MVP locally:
- ✅ 4 system metrics (CPU, Memory, Disk, Network)
- ✅ 4 load actions with safety limits
- ✅ React dashboard with real-time updates
- ✅ HTTP polling (NO WebSocket in MVP)
- ✅ In-memory storage (NO database in MVP)
- ✅ Docker Compose setup
- ✅ Tests (70% coverage minimum)

**Key Decisions (see DECISIONS.md):**
- Using HTTP polling, NOT WebSocket
- In-memory storage, NO PostgreSQL
- NO pre-built scenarios in MVP
- Folder structure: `internal/actions/`

**Week-by-Week Breakdown:**
```
Week 1: Backend + CPU action working safely
Week 2: Frontend + one button working
Week 3: All 4 actions + event log
Week 4: Charts, polish, testing

= v1.0 MVP Complete
```

**Checkpoint:** ✅ Can run `docker-compose up`, trigger actions, see metrics react

**Next Step:** Choose deployment path (Phase 1.5 or Phase 2)

---

## 🎯 Phase 1.5: Showcase Deployment (Simple Public Demo) ⭐

**Goal:** Deploy to simple cloud hosting for public showcase
**Timeline:** 2-3 days
**Prerequisites:** Phase 1 MVP complete
**Cost:** ~$5-10/month

### Why This Phase?
```
Perfect for:
✅ Portfolio/resume (live demo link)
✅ Sharing with friends/colleagues
✅ Job applications (show real working project)
✅ Quick validation before AWS investment

NOT for:
❌ Production workloads
❌ High traffic
❌ Enterprise deployment
```

### Phase 1.5.1: Choose Simple Hosting
**Goal:** Select and configure hosting platform

**Recommended Options:**

**Option A: Railway.app (RECOMMENDED)**
```
Pros:
✅ Easiest Docker Compose support
✅ Free trial, then $5-10/month
✅ Auto-deploy from GitHub
✅ Built-in SSL/HTTPS
✅ Simple environment variables
✅ Logs and metrics included

Cons:
❌ Limited free tier
❌ US-based servers only

Setup: 30 minutes
```

**Option B: Render.com**
```
Pros:
✅ Free tier available
✅ Good Docker support
✅ Auto-deploy from GitHub
✅ Built-in SSL
✅ Global CDN

Cons:
❌ Free tier sleeps after inactivity
❌ Slower cold starts

Setup: 45 minutes
```

**Option C: Fly.io**
```
Pros:
✅ Generous free tier
✅ Docker-native
✅ Edge deployment
✅ Fast performance

Cons:
❌ Slightly more complex setup
❌ Command-line heavy

Setup: 1 hour
```

### Phase 1.5.2: Deployment Setup
**Goal:** Deploy application to chosen platform

**Tasks:**
- [ ] Create account on chosen platform
- [ ] Install platform CLI (if needed)
- [ ] Configure deployment settings
  - [ ] Environment variables
  - [ ] Resource limits (CPU/Memory)
  - [ ] Port configuration
- [ ] Connect GitHub repository
- [ ] Set up auto-deploy (optional)
- [ ] Configure custom domain (optional)
  - [ ] Buy domain (~$12/year) or use free subdomain
  - [ ] Configure DNS
  - [ ] Enable SSL/HTTPS

**Checkpoint:** ✅ Application accessible via public URL

---

### Phase 1.5.3: Showcase Mode Features
**Goal:** Add features specific to public demo

**Features to Add:**
- [ ] Read-only public mode (optional login for full features)
- [ ] Rate limiting (prevent abuse)
  - [ ] Max 10 actions per IP per hour
  - [ ] Cooldown between actions
- [ ] Welcome banner
  - [ ] "This is a demo application"
  - [ ] Link to GitHub repo
  - [ ] Brief description
- [ ] About page
  - [ ] Project description
  - [ ] Technologies used
  - [ ] Your contact info
  - [ ] Link to portfolio
- [ ] Demo mode button
  - [ ] Auto-runs a quick demo sequence
  - [ ] Shows all features
- [ ] Safety for public access
  - [ ] Stricter resource limits
  - [ ] Auto-stop after 30 seconds
  - [ ] Clear all actions daily

**Implementation Time:** 1 day

**Checkpoint:** ✅ Safe for public access, impressive demo mode

---

### Phase 1.5.4: Configuration for Public Hosting

**Environment Variables:**
```env
# Backend (.env)
PORT=8080
ENV=production
METRICS_INTERVAL=2s

# Safety limits (stricter for public demo)
MAX_CPU_PERCENT=80          # Lower than local (was 95)
MAX_CPU_DURATION=20         # Shorter than local (was 30)
MAX_MEMORY_PERCENT=20       # Lower than local (was 25)
MAX_CONCURRENT=3            # Fewer than local (was 5)

# Rate limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_PER_HOUR=10
ACTION_COOLDOWN_SECONDS=5

# CORS
ALLOWED_ORIGINS=https://yourapp.railway.app

# Frontend (.env)
VITE_API_URL=https://yourapp.railway.app
VITE_DEMO_MODE=true
```

**Docker Optimization for Cloud:**
```dockerfile
# Smaller images for faster deployment
FROM golang:1.21-alpine AS builder  # Use alpine
# ... build steps ...
FROM alpine:latest                   # Minimal runtime
```

---

### Phase 1.5.5: Testing Public Deployment
**Goal:** Verify everything works online

**Test Checklist:**
- [ ] Application loads via HTTPS URL
- [ ] All 4 actions work
- [ ] Metrics update in real-time
- [ ] Charts render correctly
- [ ] Event log works
- [ ] Demo mode works
- [ ] Rate limiting works (test with multiple requests)
- [ ] Mobile responsive
- [ ] Fast load time (<3s)
- [ ] No console errors
- [ ] SSL certificate valid
- [ ] Links work (GitHub, portfolio, etc.)

**Performance Targets:**
- Initial load: <3 seconds
- Action response: <500ms
- Metric updates: <2 seconds
- Chart rendering: <100ms

**Checkpoint:** ✅ All features work, performance good

---

### Phase 1.5.6: Documentation & Sharing
**Goal:** Make it easy to share and understand

**Create:**
- [ ] Landing page explanation
  ```
  # Interactive System Monitor Dashboard

  Click buttons to trigger system load and watch
  metrics react in real-time!

  Try it:
  1. Click "🔥 CPU Stress"
  2. Watch the CPU graph spike
  3. See the event log update

  [Start Demo Mode]
  ```

- [ ] README badge
  ```markdown
  ## 🌐 Live Demo
  Check out the live demo: [https://yourapp.railway.app](https://yourapp.railway.app)
  ```

- [ ] Share on:
  - [ ] Your portfolio website
  - [ ] LinkedIn
  - [ ] GitHub profile README
  - [ ] Resume (as clickable link)

**Checkpoint:** ✅ Easy to share and impress

---

### Phase 1.5.7: Monitoring & Maintenance
**Goal:** Keep it running smoothly

**Set Up:**
- [ ] Uptime monitoring
  - [ ] UptimeRobot (free tier)
  - [ ] Check every 5 minutes
  - [ ] Email alert if down
- [ ] Simple analytics (optional)
  - [ ] Google Analytics
  - [ ] Track visitors
  - [ ] See popular features
- [ ] Monthly check-in
  - [ ] Verify still working
  - [ ] Check logs for errors
  - [ ] Review costs

**Checkpoint:** ✅ Reliable public demo

---

### Phase 1.5 Complete Checklist

**Functionality:**
- [ ] Accessible via public HTTPS URL
- [ ] All MVP features work online
- [ ] Demo mode impresses visitors
- [ ] Rate limiting prevents abuse
- [ ] Safe for public access

**Quality:**
- [ ] Fast load time (<3s)
- [ ] No errors in console
- [ ] Works on mobile
- [ ] SSL valid
- [ ] Professional appearance

**Documentation:**
- [ ] About page exists
- [ ] Links to GitHub work
- [ ] Portfolio integration done
- [ ] Resume updated with link

**Business Value:**
- [ ] Can share in job applications
- [ ] Works for portfolio
- [ ] Impresses viewers
- [ ] Shows technical skills

**Deliverable:** 🌐 **Live public demo at [your-url]**

**Example URLs:**
- `https://sysmonitor-demo.railway.app`
- `https://monitor-dashboard.onrender.com`
- `https://sysmon.fly.dev`
- Or custom: `https://demo.yourname.com`

---

### Phase 1.5 Hosting Cost Comparison

```
Railway.app:
- Free: $0/month (trial credits)
- Starter: $5/month
- Pro: $20/month
→ RECOMMENDED: Start with trial

Render.com:
- Free: $0/month (with sleep)
- Starter: $7/month (always on)
→ Good for testing

Fly.io:
- Free: $0-5/month (generous free tier)
- Pay as you go
→ Good for minimal cost

+ Domain (optional): ~$12/year
Total: $5-10/month OR free with trial/free tiers
```

---

### When to Choose Phase 1.5 vs Phase 2

**Choose Phase 1.5 (Showcase) if:**
- ✅ Want to show it online quickly
- ✅ Need portfolio link
- ✅ Budget conscious ($10/month OK)
- ✅ Don't need "AWS experience" on resume
- ✅ Want simple deployment

**Skip to Phase 2 (AWS) if:**
- ✅ Want AWS on resume
- ✅ Learning AWS is the goal
- ✅ Budget allows ($100+/month)
- ✅ Want production-grade deployment
- ✅ Targeting DevOps/Cloud roles

**Or do BOTH:**
1. Phase 1.5 first (get it online fast)
2. Phase 2 later (add AWS experience when ready)

---

## 🎯 Phase 2: AWS Dev Environment

**Goal:** Professional AWS deployment for development/testing
**Timeline:** 2-3 weeks
**Prerequisites:** Phase 1 complete (Phase 1.5 optional)
**Cost:** ~$75-100/month

### Phase 2.1: AWS Infrastructure Setup
**Goal:** Terraform infrastructure for dev environment

- [ ] Set up AWS account and credentials
- [ ] Create Terraform configuration
  - [ ] `infrastructure/terraform/dev/`
  - [ ] Provider configuration
  - [ ] State backend (S3 + DynamoDB)
- [ ] Create VPC and networking
  - [ ] VPC with public/private subnets
  - [ ] Internet Gateway
  - [ ] NAT Gateway
  - [ ] Security Groups
- [ ] Create ECR repositories
  - [ ] Backend image repository
  - [ ] Frontend image repository
- [ ] Create ECS cluster (Fargate)
  - [ ] Cluster definition
  - [ ] Task definitions
  - [ ] Service definitions
- [ ] Create Application Load Balancer
  - [ ] Target groups
  - [ ] Listener rules
  - [ ] Health checks
- [ ] Create RDS PostgreSQL instance
  - [ ] Dev-sized instance (t3.micro)
  - [ ] Security group rules
  - [ ] Automated backups
- [ ] Set up CloudWatch
  - [ ] Log groups
  - [ ] Alarms (CPU, Memory, Health)
- [ ] Apply Terraform
  - [ ] `terraform init`
  - [ ] `terraform plan`
  - [ ] `terraform apply`
- [ ] Test: Infrastructure created successfully

**Checkpoint:** ✅ AWS dev infrastructure provisioned

---

### Phase 2.2: Docker Image Build & Push
**Goal:** Images in ECR

- [ ] Build backend Docker image
  - [ ] Multi-stage build for optimization
  - [ ] Tag with git commit SHA
- [ ] Build frontend Docker image
  - [ ] Production build
  - [ ] Nginx configuration
  - [ ] Tag with git commit SHA
- [ ] Push images to ECR
  - [ ] Authenticate with ECR
  - [ ] Push backend image
  - [ ] Push frontend image
- [ ] Test: Images available in ECR

**Checkpoint:** ✅ Docker images in AWS ECR

---

### Phase 2.3: Manual Deployment to Dev
**Goal:** Running application on AWS

- [ ] Update ECS task definitions with image URIs
- [ ] Configure environment variables
  - [ ] Backend env vars
  - [ ] Frontend API URLs
  - [ ] Database connection string
- [ ] Deploy to ECS
  - [ ] Force new deployment
  - [ ] Wait for services to stabilize
- [ ] Configure ALB routing
  - [ ] Frontend → port 80
  - [ ] Backend API → /api/*
  - [ ] WebSocket → /ws
- [ ] Test deployment
  - [ ] Access via ALB DNS
  - [ ] Test all features
  - [ ] Verify metrics collection
  - [ ] Verify actions work
- [ ] Set up DNS (optional)
  - [ ] Route53 hosted zone
  - [ ] A record to ALB
  - [ ] SSL certificate (ACM)

**Checkpoint:** ✅ Application running on AWS dev environment

---

### Phase 2.4: Monitoring & Debugging
**Goal:** Observability in dev environment

- [ ] Set up CloudWatch dashboards
  - [ ] ECS metrics
  - [ ] Application metrics
  - [ ] RDS metrics
- [ ] Configure log aggregation
  - [ ] Container logs to CloudWatch
  - [ ] Application logs
  - [ ] Error logs
- [ ] Set up alerts
  - [ ] High CPU/Memory alerts
  - [ ] Health check failures
  - [ ] Error rate alerts
- [ ] Test monitoring
  - [ ] Trigger actions
  - [ ] Check CloudWatch metrics
  - [ ] Verify logs appear
- [ ] Document troubleshooting steps

**Checkpoint:** ✅ Full observability in dev environment

---

## 🎯 Phase 3: AWS Production Environment

**Goal:** Production-ready deployment
**Timeline:** 1-2 weeks
**Prerequisites:** Phase 2 tested and stable

### Phase 3.1: Production Infrastructure
**Goal:** Production-grade AWS setup

- [ ] Create Terraform for production
  - [ ] `infrastructure/terraform/prod/`
  - [ ] Separate state backend
- [ ] Enhanced networking
  - [ ] Multi-AZ setup
  - [ ] NAT Gateways in multiple AZs
- [ ] Production-sized resources
  - [ ] ECS tasks with appropriate resources
  - [ ] RDS instance (t3.small or larger)
  - [ ] Multi-AZ RDS for HA
- [ ] Security hardening
  - [ ] Restrictive security groups
  - [ ] SSL/TLS enforced
  - [ ] Secrets Manager for credentials
  - [ ] IAM roles with least privilege
- [ ] Backup and recovery
  - [ ] RDS automated backups
  - [ ] Backup retention policy
  - [ ] Disaster recovery plan
- [ ] Apply Terraform for prod

**Checkpoint:** ✅ Production infrastructure provisioned

---

### Phase 3.2: CI/CD Pipeline
**Goal:** Automated deployments

- [ ] Create GitHub Actions workflows
  - [ ] `.github/workflows/backend-ci.yml`
  - [ ] `.github/workflows/frontend-ci.yml`
  - [ ] `.github/workflows/deploy-dev.yml`
  - [ ] `.github/workflows/deploy-prod.yml`
- [ ] Backend CI pipeline
  - [ ] Run tests on PR
  - [ ] Lint code
  - [ ] Build Docker image
  - [ ] Push to ECR
- [ ] Frontend CI pipeline
  - [ ] Run tests on PR
  - [ ] Lint code
  - [ ] Build Docker image
  - [ ] Push to ECR
- [ ] Dev deployment pipeline
  - [ ] Auto-deploy on merge to `develop`
  - [ ] Update ECS services
  - [ ] Health check verification
- [ ] Prod deployment pipeline
  - [ ] Manual approval required
  - [ ] Deploy on merge to `main`
  - [ ] Blue-green deployment
  - [ ] Rollback capability
- [ ] Test pipelines end-to-end

**Checkpoint:** ✅ Automated CI/CD pipelines working

---

### Phase 3.3: Production Deployment
**Goal:** Live production system

- [ ] Pre-deployment checklist
  - [ ] All tests passing
  - [ ] Security review complete
  - [ ] Performance testing done
  - [ ] Monitoring configured
- [ ] Initial production deployment
  - [ ] Deploy via pipeline
  - [ ] Verify health checks
  - [ ] Test all features
- [ ] Configure custom domain
  - [ ] SSL certificate
  - [ ] DNS configuration
  - [ ] HTTPS redirect
- [ ] Final testing
  - [ ] Full user journey test
  - [ ] Load testing
  - [ ] Security scan
- [ ] Documentation update
  - [ ] Architecture diagram
  - [ ] Deployment guide
  - [ ] Runbook for incidents

**Checkpoint:** ✅ Production system live and stable

---

### Phase 3.4: Production Monitoring & Optimization
**Goal:** Maintain and improve production

- [ ] Set up production monitoring
  - [ ] CloudWatch dashboards
  - [ ] Custom metrics
  - [ ] Error tracking
- [ ] Configure alerts
  - [ ] PagerDuty/SNS integration
  - [ ] On-call rotation (if team)
  - [ ] Escalation policies
- [ ] Performance optimization
  - [ ] Identify bottlenecks
  - [ ] Optimize database queries
  - [ ] CDN for static assets
  - [ ] Caching strategy
- [ ] Cost optimization
  - [ ] Right-size resources
  - [ ] Reserved instances
  - [ ] Spot instances for non-critical tasks
- [ ] Regular maintenance
  - [ ] Security patches
  - [ ] Dependency updates
  - [ ] Performance reviews

**Checkpoint:** ✅ Production system optimized and monitored

---

## 📊 Progress Tracking

### Overall Progress

**Phase 1 - Local MVP:** [ ] 0/10 sections complete (0%)
**Phase 2 - AWS Dev:** [ ] 0/4 sections complete (0%)
**Phase 3 - AWS Prod:** [ ] 0/4 sections complete (0%)

### Current Focus
**Active Phase:** Phase 1 - Local MVP
**Current Section:** Not started
**Next Milestone:** Foundation & Setup complete

---

## 🚀 Quick Start Guide

### Getting Started (Phase 1)

1. **Clone repository:**
```bash
git clone <repository-url>
cd monitoring-dashboard
```

2. **Start local development:**
```bash
docker-compose up
```

3. **Access application:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Docs: http://localhost:8080/api/docs

### Migration to Cloud (Phase 2+)

See Phase 2 and Phase 3 sections for detailed cloud deployment steps.

---

## 📁 Repository Structure

```
monitoring-dashboard/
├── backend/                  # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── metrics/         # Metrics collection
│   │   ├── loadgen/         # Load generation engine
│   │   ├── websocket/       # WebSocket hub
│   │   └── api/             # REST API handlers
│   ├── pkg/
│   │   └── models/          # Shared models
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
│
├── frontend/                # React frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── Dashboard/
│   │   │   ├── ControlPanel/
│   │   │   ├── EventLog/
│   │   │   └── Shared/
│   │   ├── services/
│   │   │   ├── websocket.js
│   │   │   └── api.js
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── Dockerfile
│
├── infrastructure/          # AWS infrastructure
│   └── terraform/
│       ├── dev/            # Dev environment
│       └── prod/           # Prod environment
│
├── .github/
│   └── workflows/          # CI/CD pipelines (Phase 3)
│
├── docker-compose.yml      # Local development
├── README.md               # Project documentation
└── RELEASE-PLAN.md         # This file
```

---

## 🎯 Success Criteria

### Phase 1 Complete When:
- [ ] `docker-compose up` starts all services
- [ ] Frontend accessible at localhost:3000
- [ ] Can trigger all 4 actions via UI
- [ ] Metrics update in real-time
- [ ] Event log shows action lifecycle
- [ ] Charts are animated and responsive
- [ ] All tests passing (>80% coverage)
- [ ] README has clear setup instructions

### Phase 2 Complete When:
- [ ] Application running on AWS dev
- [ ] Accessible via public URL
- [ ] All features work identically to local
- [ ] Monitoring dashboards configured
- [ ] Logs aggregated in CloudWatch

### Phase 3 Complete When:
- [ ] Production environment live
- [ ] CI/CD pipeline deploying automatically
- [ ] Custom domain with HTTPS
- [ ] Production monitoring active
- [ ] Disaster recovery plan tested

---

## 💡 Development Tips

### Working with Docker Compose
```bash
# Start services
docker-compose up

# Rebuild after code changes
docker-compose up --build

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Stop services
docker-compose down

# Clean everything
docker-compose down -v
```

### Backend Development
```bash
cd backend

# Run locally (outside Docker)
go run cmd/server/main.go

# Run tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Hot reload (install air)
air
```

### Frontend Development
```bash
cd frontend

# Install dependencies
npm install

# Run locally (outside Docker)
npm run dev

# Run tests
npm test

# Build for production
npm run build
```

---

## 📝 Notes

- **Database:** PostgreSQL included in docker-compose but not required for MVP
- **Persistence:** Local data lost on container restart (dev only)
- **Cloud migration:** Architecture designed for easy AWS transition
- **Costs:** AWS dev ~$50/month, prod ~$100-150/month
- **Timeline:** Flexible, adjust based on progress

---

**Last Updated:** 2025-01-09
**Status:** Planning Phase
**Next Action:** Initialize repository and start Phase 1.1
