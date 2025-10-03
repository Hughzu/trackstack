# Building TrackStack: A Tale of Two Architectures

## Why This Project?

As a software engineer with 7 years of experience transitioning into SRE/DevOps, I wanted to build something that demonstrates not just what I *can* do, but how I *think* about software architecture. Too many portfolios show a single implementation - but the real skill in our field is understanding **trade-offs**.

So I'm building the same application twice: once for reality, once for resume.

## Meet TrackStack

TrackStack is a personal tracking application I actually use daily to manage:
- **Expenses** - where my money goes
- **Gym routines** - tracking reps and progress
- **Heating consumption** - monitoring energy usage

It's nothing revolutionary, but it's *real*. And that authenticity matters when demonstrating architectural decisions.

## The Two Versions

### The Pragmatic Version: $3-5/month

This is the version I actually run. Every day. With my real data.

**Stack:**
- **Frontend:** Astro (static site, blazing fast, SEO-friendly) as a PWA
- **Backend:** C# monolith in a Docker container on AWS Fargate Spot
- **Database:** SQLite on EFS
- **Messaging:** AWS SQS for event-driven patterns
- **Hosting:** S3 + CloudFront for frontend, single Fargate task for backend

**Why these choices?**
- **Cost first:** Running 24/7 for less than a coffee
- **Simple operations:** One container to deploy, one database file to back up
- **Battle-tested:** C# and SQLite are boring technology (in the best way)
- **Event-driven foundation:** Even the monolith publishes domain events, preparing for future evolution

**Monthly cost breakdown:**
- Fargate Spot: $1-2
- EFS (SQLite): $0.30
- S3 + CloudFront: $1
- SQS: Essentially free (1M requests/month free tier)

### The Enterprise Version: $30-32/month (weekend-only)

This is the version I show clients. It demonstrates what I can architect when scale, resilience, and observability matter.

**Stack:**
- **Frontend:** Same Astro PWA (code reuse!)
- **Backend:** Microservices architecture
  - Expenses Service
  - Gym Service  
  - Heating Service
  - Each service: C# REST API + SQS event consumer
- **Database:** SQLite per service on EFS (will migrate to RDS Postgres in v2)
- **Messaging:** 
  - AWS SNS for domain events (pub/sub)
  - AWS SQS queues per service (guaranteed delivery)
- **API Layer:** AWS API Gateway (managed)
- **Infrastructure:**
  - VPC with public/private subnets across 2 AZs
  - Application Load Balancer
  - ECS Fargate for container orchestration
  - NAT Gateway for private subnet internet access
- **Observability:** Prometheus + Grafana (self-hosted on Lightsail)
- **IaC:** 100% Terraform
- **CI/CD:** GitHub Actions with OIDC (no long-lived credentials)

**Why these choices?**
- **Domain-Driven Design:** Each service owns its bounded context
- **Event-driven architecture:** Services communicate via events, fully decoupled
- **Eventual consistency:** Services maintain denormalized data, listening to events from other domains
- **Multi-AZ resilience:** Survives availability zone failures
- **Modern DevOps practices:** IaC, OIDC authentication, immutable infrastructure

**Cost management:**
This version only runs on **weekends** (Saturdays and Sundays, 8 AM - 10 PM). Terraform automatically:
- Spins up Saturday/Sunday mornings
- Tears down each evening at 10 PM
- **~112 hours/month runtime**

This keeps monthly costs around $30-32 instead of $150+ if running 24/7.

## The Data Bridge: Pragmatic → Enterprise

Here's where it gets interesting: **the enterprise version uses real data from the pragmatic version**.

Every night, the pragmatic version backs up its SQLite database to S3. When the enterprise version spins up on weekends:

1. A Lambda function triggers on Terraform apply
2. Downloads the latest pragmatic backup from S3  
3. Splits the monolithic database into domain-specific databases
4. Each microservice starts with fresh, real-world data

This means I'm not demoing with fake data. I'm showing how the same information flows through two completely different architectures.

## The Architectural Philosophy

**Pragmatic version** answers: *"What would I build if I were a solo developer with limited budget?"*
- Optimize for cost and operational simplicity
- Leverage managed services where possible
- Accept single points of failure for non-critical personal use
- Event-driven foundation allows future evolution

**Enterprise version** answers: *"What would I build if I were architecting for a scaling startup or enterprise?"*
- Optimize for resilience, observability, and team autonomy  
- Enable independent service deployment and scaling
- Design for failure (multi-AZ, graceful degradation)
- Demonstrate modern cloud-native patterns

Neither is "better" - they're **appropriate for different contexts**. That's the whole point.

## What's Next?

Over the coming weeks and months, I'll be documenting:

1. **Infrastructure deep-dives:** Terraform modules, networking decisions, security postures
2. **Event-driven patterns:** How domain events flow through the system
3. **Cost optimizations:** Every dollar saved in the pragmatic version
4. **Observability:** Prometheus metrics, Grafana dashboards, alerting strategies
5. **CI/CD pipelines:** GitHub Actions workflows, OIDC setup, deployment strategies
6. **Database evolution:** Migrating from SQLite to RDS Postgres in the enterprise version
7. **Azure parallel:** Building the same enterprise architecture on Azure for comparison

## Why You Should Care

If you're a potential client or employer reading this, here's what this project demonstrates:

✅ **Cost consciousness** - I understand business constraints  
✅ **Architectural judgment** - I know when complexity is warranted  
✅ **Modern DevOps practices** - IaC, OIDC, event-driven architecture  
✅ **Communication skills** - I can explain technical decisions clearly  
✅ **Real-world experience** - This isn't a tutorial project, it's production software I use daily

## The Tech Stack at a Glance

**Shared:**
- Frontend: Astro (PWA)
- Backend Language: C#
- Containerization: Docker
- Cloud Provider: AWS (Azure version coming)
- IaC: Terraform  
- CI/CD: GitHub Actions with OIDC
- Messaging: SQS/SNS

**Pragmatic-specific:**
- Compute: Fargate Spot (single task)
- Database: SQLite on EFS
- Architecture: Monolith with event-driven patterns

**Enterprise-specific:**
- Compute: Fargate (microservices)
- Database: SQLite per service (RDS Postgres in v2)
- Architecture: Domain-driven microservices
- API Gateway: AWS API Gateway (managed)
- Load Balancer: Application Load Balancer
- Observability: Prometheus + Grafana
- Multi-AZ: Yes

## Follow Along

I'll be publishing detailed articles as I build this. Topics will include technical deep-dives, cost breakdowns, and lessons learned. Whether you're interested in cloud architecture, DevOps practices, or just want to see how someone thinks through real-world trade-offs, there should be something here for you.

The code will be open source (coming soon), and I'll be documenting both successes and failures. Because let's be honest - the mistakes are often more educational than the wins.

---

**Next post:** Setting up the foundational AWS infrastructure with Terraform - VPC, subnets, security groups, and the networking decisions that matter.

Stay tuned, and feel free to reach out with questions or suggestions!

---

*Last updated: [Your Date]*  
*GitHub: [Your GitHub]*  
*LinkedIn: [Your LinkedIn]*



Suggested Edits:

Add your actual dates, links, and personal info at the bottom
Adjust the tone to match your voice (I kept it professional but approachable)
If you want it more casual or more formal, let me know
Add a hero image/architecture diagram when you have one

Follow-Up Blog Ideas:

"Terraform Modules for TrackStack: Building Reusable Infrastructure"
"OIDC with GitHub Actions: Eliminating AWS Access Keys"
"Event-Driven Monoliths: The Missing Middle Ground"
"SQLite in Production: When Simple is Smart"
"Cost Breakdown: How I Run TrackStack for $3/month"
"Weekend-Only Infrastructure: Terraform Strategies for Demo Environments"

Want me to adjust anything in the blog post, or ready to move on to planning the actual implementation?
