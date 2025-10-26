---
title: 'What is TrackStack'
description: 'Here is a sample of some basic Markdown syntax that can be used when writing Markdown content in Astro.'
pubDate: '2025-10-26'
heroImage: '../../assets/blog-placeholder-2.jpg'
---

## The Problem: Ideas Without Execution

As a software engineer with 8 years of experience transitioning into SRE/DevOps, I have a constant stream of ideas. Small frustrations in daily life that I *know* I could solve with code. But here's the thing: **ideas are worthless without execution**, and execution is expensive.

I needed a framework. A system that lets me:
1. Validate multiple ideas quickly
2. Use them daily to prove their value
3. Evolve the good ones into portfolio-worthy architectures
4. Do all of this without burning out on premature optimization

So I'm building TrackStack: **a personal tracking platform that evolves from scrappy MVP to enterprise-grade architecture**.

## Business First, Technology Second

Let me be clear: I'm not building this to learn a new framework or because microservices are trendy. I'm building it because I have **real, daily frustrations**:

- **Expense tracking apps** are either too complex or too simplistic. I want something that just works for *my* workflow.
- **Heat consumption monitoring** - I burn wood pellets for heating. How many bags per month? What's the correlation with temperature? Am I being efficient? No app does this.
- **Recipe management** - My girlfriend and I cook together. We need a shared recipe book that actually tracks nutritional information properly.
- **Calorie tracking** - Every app wants to sell me a subscription and spam me with notifications. I just want to log my calories and move on.

These aren't revolutionary ideas. They're *boring problems* that I face every single day. And that's exactly why they're perfect for this project.

**No solution exists until you validate the problem.** So I'm starting with proof-of-concept, using it daily, and *then* - only then - investing in architecture.

## The Evolution: Four Phases

### Phase 1: The Scrappy MVP

**Goal:** Validate all four ideas. Use them daily. Figure out what works.

**Stack:**
- **Frontend:** Astro as the shell application (static, fast, SEO-friendly)
- **"Backend":** Initially embedded in Astro (yes, really)
- **Apps:** Separate Angular applications for each domain
  - Expense Tracker
  - Heat Monitor
  - Recipe Book
  - Calorie Counter
- **Database:** Whatever works fastest (probably JSON files or SQLite)
- **Hosting:** AWS S3 + CloudFront (cheapest possible setup)

**Architecture Philosophy:**
- Astro acts as the **shell** - handles authentication, routing, and layout
- Angular apps are **micro-frontends** - each one isolated and independently deployable
- Design for extraction: Each Angular app is architected so it can transition to a standalone application with minimal effort if it proves valuable

**Why this approach?**
- **Speed:** I can ship an idea in days, not weeks
- **Validation:** Real daily usage tells me which ideas have legs
- **Pragmatism:** Minimal infrastructure costs, minimal DevOps overhead, just pure product validation
- **Exit strategy:** If "Heat Monitor" becomes something bigger, I can extract it without rewriting everything

**Success criteria:** After consistent daily use, I'll know which ideas are keepers and which are dead weight.

---

### Phase 2: The Pragmatic Backend

**Goal:** Extract the backend, make it cheap and maintainable.

By now, I know which features I actually use. Time to build something real, but **keep it pragmatic**.

**Stack:**
- **Frontend:** Same Astro + Angular apps (if it works, don't break it)
- **Backend:** C# monolith with hexagonal architecture
  - **Why hexagonal?** Domain logic isolated from infrastructure
  - **Why monolith?** I don't have Google's problems. I have *my* problems.
  - **Why C#?** Fast, type-safe, mature ecosystem
- **Database:** SQLite (or similar lightweight option)
  - My focus is on cost/performance, and for my load, there's no need for expensive database infrastructure
  - Using an ORM means I can pivot the database implementation if needed
  - I want to demonstrate broad architectural principles - the specific database technology isn't the point
- **Messaging:** Event-driven patterns with a lightweight message broker
- **Hosting:** 
  - Frontend: S3 + CloudFront (pennies per month)
  - Backend: AWS Fargate Spot (staying under $5/month total)

**Architecture highlights:**
- **Hexagonal architecture** means I can swap databases/queues/APIs without touching business logic
- **Event-driven patterns** even in the monolith - domain events flow through the system
- **Domain boundaries** clearly defined, making future service extraction straightforward if needed

**Monthly cost target: $3-5**

**Why SQLite?** Because I'm not Netflix. My "scale" is me, my girlfriend, and maybe a few friends. The operational simplicity is unbeatable, and having multiple expensive databases just to demonstrate domain splitting doesn't make sense.

---

### Phase 3: Observability on a Budget

**Goal:** Instrument everything. Know when things break before users do.

**What I'll add:**
- **Logging:** Structured logging infrastructure
- **Metrics:** Prometheus + Grafana (or similar cost-effective solution after research)
- **Alerts:** Critical paths only (don't wake me up at 3 AM unless it matters)
- **Dashboards:**
  - Request latency (p50, p95, p99)
  - Error rates by endpoint
  - Database query performance
  - Cost tracking (AWS bills integrated)

**Target cost increase: $5-8/month**

**Philosophy:** Observability isn't optional. Even personal projects need monitoring. This phase teaches me to build sustainable systems, not just "works on my machine" code.

---

### Phase 4: The Enterprise-Grade Architecture

**Goal:** Demonstrate what I can architect at scale for potential employers and clients.

Now we get serious. This version runs **on-demand** - when someone wants to see it - but I'll document it extensively and explain every architectural decision.

**Stack:**
- **Frontend:** Same Astro + Angular (code reuse!)
- **Backend:** Microservices architecture
  - Expense Service
  - Heat Monitor Service
  - Recipe Service
  - Calorie Tracker Service
  - Each service: C# REST API with dedicated event consumers
- **Database:** Lightweight per-service databases (leveraging ORM flexibility)
- **Messaging:** Event-driven pub/sub architecture
- **API Layer:** Managed API gateway
- **Infrastructure:**
  - Multi-AZ deployment for resilience
  - Container orchestration
  - Proper network segmentation (public/private subnets)
  - Load balancing
  - NAT gateways for secure egress
- **Observability:** Enhanced monitoring with service-level dashboards
- **IaC:** 100% Terraform
- **CI/CD:** GitHub Actions with OIDC (no long-lived credentials)

**Architecture highlights:**
- **Domain-Driven Design:** Each service owns its bounded context
- **Event-driven architecture:** Services communicate via events, fully decoupled
- **Eventual consistency:** Services maintain their own read models
- **Multi-AZ resilience:** Survives availability zone failures as per SRE best practices
- **Immutable infrastructure:** Every deploy is a fresh container

**Cost management - FinOps best practices:**
- Automatic shutdown during non-business hours (no point paying while I sleep)
- Resource right-sizing based on actual usage patterns
- Spot instances where appropriate
- Careful monitoring of spend vs. value

I'll provide detailed infrastructure schemas and explain multi-AZ deployment strategies, service mesh considerations, and all architectural trade-offs in dedicated blog posts.

---

### Phase 5: Multi-Cloud Architecture

Once AWS is dialed in, rebuild the enterprise-grade architecture on **Azure** to demonstrate:
- Cloud-agnostic architectural thinking
- Platform-specific services and their AWS equivalents
- Multi-cloud cost comparisons
- Infrastructure-as-Code portability

This isn't about showing off - it's about proving I can think beyond a single cloud provider's ecosystem.

---

## The Data Bridge: Reality → Demo

Here's the brilliant part: **The enterprise version uses real data from the pragmatic version.**

Every night:
1. Pragmatic version backs up its database to S3
2. Enterprise startup (on-demand):
   - Automated process triggers on infrastructure deployment
   - Downloads latest pragmatic backup
   - Splits monolithic database into domain-specific databases
   - Each microservice starts with real-world data

**I'm not demoing with fake data. I'm showing how the same real-world usage flows through two architectures.**

---

## What This Demonstrates

If you're a potential client or employer, here's what you're seeing:

✅ **Business-first thinking** - Technology serves problems, not egos  
✅ **Pragmatic engineering** - Build what you need, not what's trendy  
✅ **Cost consciousness** - Every dollar matters  
✅ **Architectural evolution** - Monolith → Microservices done right  
✅ **Modern DevOps** - IaC, OIDC, event-driven patterns, observability  
✅ **SRE mindset** - Multi-AZ resilience, proper monitoring, FinOps practices  
✅ **Real-world validation** - Production software I use daily  
✅ **Communication skills** - I can explain decisions clearly

---

## Follow Along

The code will be open source (coming soon). I'll document successes, failures, and everything in between. Because the best learning happens in the messy middle, not the polished finish.

**Next post:** Phase 1 - Building the Astro + Angular MVP architecture. How to design micro-frontends that can evolve into standalone apps with minimal effort.

Stay tuned!

---

*GitHub: https://github.com/Hughzu/trackstack/*

---

**Note to self:** As a young dad, I work when I have time. This isn't a sprint - it's a marathon with nap-time interruptions. And that's perfectly fine.
