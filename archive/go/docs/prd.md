---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]
inputDocuments: ['_bmad-output/planning-artifacts/product-brief-trackstack-2026-01-23.md', '_bmad-output/analysis/brainstorming-session-2026-01-23.md']
workflowType: 'prd'
classification:
  projectType: web_app
  domain: general
  complexity: medium
  projectContext: greenfield
---

# Product Requirements Document - trackstack

**Author:** Hsi
**Date:** 2026-01-23

## Executive Summary

**TrackStack** is a personal observability platform designed to solve two distinct problems simultaneously: the functional need to track personal "Drift" (Financial, Physical, Inventory) and the professional need to overcome Imposter Syndrome by demonstrating SRE mastery.

Unlike typical tracking apps which focus on historical logging, TrackStack focuses on **Drift Prevention** (managing "remainders" and "capacities") via a frictionless "Lazy UX" (AI Ingest, Recents Carousel). Simultaneously, the project serves as a "Full-Stack Thesis," proving the author's ability to engineer a **Modular Monolith** in Go that can be refactored into microservices via purely infrastructure changes (Terraform), offering a "Zero-Code-Change" deployment flexibility that serves as a powerful portfolio asset.

## Core Vision

### Problem Statement
The user, a Site Reliability Engineer (SRE), faces a dual challenge:
1.  **Personal Drift:** Without low-friction tracking, they "drift" into failure states (overspending, weight gain, running out of heating pellets).
2.  **Professional Imposter Syndrome:** Despite being employed as an SRE, they feel a need to prove deep competence across the entire engineering lifecycle (Dev -> Arch -> Infra -> Obs), particularly in Go and Cloud-Native patterns.

### Problem Impact
*   **Financial/Health:** "Stupid decisions" lead to overspending and weight gain.
*   **Inventory:** Heating pellet shortages risk "Winter Bollocks" (running out mid-season).
*   **Career:** Lack of a tangible, complex portfolio piece reinforces feelings of technical inadequacy.

### Proposed Solution
**TrackStack** is a "Life Observability Platform" built as a **Go Modular Monolith**.
*   **UX:** "Pain-Driven" interfaces including AI Voice/Text Ingest and "Recents Carousels" to minimize data entry friction.
*   **Architecture:** A cellular monolith where modules (Money, Calories, Heat, Body) are isolated by strict interfaces.
*   **Infrastructure:** A Terraform-managed AWS deployment that can deploy the app as a single binary (Monolith) or distinct services (Microservices) without changing application code.

### Key Differentiators
*   **The "Infrastructure-Driven" Refactor:** The unique ability to split the monolith into microservices via Terraform configuration alone, proving elite SRE skills.
*   **"Lazy" Observability:** Features designed for the absolute minimum effective dose of interaction (e.g., "3-Week Pulse" for heating, "Meal Totals" instead of ingredients).
*   **Dual-Value Prop:** It is both a daily utility for life management and a living resume for career advancement.

## Target Users

### Primary Users

#### **Hsi (The Dual-Role Creator)**
*   **Role 1: The "Lazy Human" (End User)**
    *   **Motivation:** Prevent "Drift" (weight gain, overspending) with minimal effort.
    *   **Success Vision:** "I tap one button or say one sentence, and I know exactly how many calories/Euros I have left for the month."
*   **Role 2: The "SRE Builder" (Engineer)**
    *   **Motivation:** Overcome Imposter Syndrome; prove mastery of Go, Terraform, and Architecture.
    *   **Success Vision:** "I can show a recruiter how I split a monolith into microservices just by changing a Terraform variable."

### Secondary Users

#### **The Partner (Independent User)**
*   **Context:** Managing her own personal expenses and data within the same system.
*   **Engineering Implication:** Standard Authentication/Authorization (User ID scoping) is required.

#### **The Technical Recruiter / Interviewer**
*   **Context:** Reviewing Hsi's portfolio for < 5 minutes.
*   **Success Vision:** "This candidate understands the trade-offs between Monoliths and Microservices and has the code to prove it."

## Success Criteria

### User Success
*   **Low Friction:** Time-to-Log for recurring items < 5 seconds.
*   **Immediate Insight:** User sees "Remaining Calories" immediately upon opening the app (No clicks required).
*   **Habit Formation:** Active usage > 5 days/week for the first month.

### Business Success (Portfolio & Cost)
*   **Cost Efficiency:** Total monthly AWS bill < €7 (Goal: Serverless Containers or T4g.nano).
*   **Engineering Thesis:** Successful deployment of the "Pilot Slice" (Calories) to AWS via Terraform.
*   **Documentation:** 1 comprehensive blog post explaining the "Zero-Code-Change" architecture.

### Technical Success
*   **Zero-Code-Change Deployment:** The application binary must support both Monolithic and Microservice modes purely via configuration/dependency injection.
*   **Deployment Success:** 100% success rate for deploying the MVP architecture to AWS.
*   **Performance:** Performance/Latency tuning is explicitly deferred to Phase 2.

## Product Scope

### MVP - Minimum Viable Product (The "Pilot Slice")
*   **Domain:** Calories & Macros ONLY.
*   **View:** Daily "What Remains" Dashboard (No historical charts).
*   **Input:** Recents Carousel + AI Text Ingest.
*   **Infrastructure:** Modular Monolith on AWS (Terraform-managed).

### Growth Features (Post-MVP)
*   **Domains:** Money, Heating, Body measurements.
*   **Views:** Historical trends, charts, and analysis.
*   **Technical:** Microservices deployment performance tuning, Full LGTM Observability stack.

### Vision (Future)
*   **Full "Life Observability Platform":** A single dashboard for health, wealth, and home inventory, running on a sophisticated, self-healing cloud-native architecture that serves as a permanent portfolio showcase.

## User Journeys

### Journey 1: The "Lazy" Feed (Daily Usage - Primary)
*   **Persona:** Hsi (The Lazy Human)
*   **Scene:** 8:00 PM on a Tuesday. Hsi is tired after work and just finished a chicken sandwich. He wants to track it but has zero mental energy for forms.
*   **Action:** He opens the TrackStack PWA on his phone. He taps the microphone icon and says, "I had a chicken sandwich and a coke."
*   **System Response:** The system uses the Gemini API to parse the intent. It identifies "Chicken Sandwich (approx 450kcal)" and "Coke (140kcal)". It updates the daily log.
*   **Outcome:** The dashboard immediately updates to show "400 Kcal Remaining."
*   **Key Requirement:** AI Text/Voice Ingest with low latency (< 2s).

### Journey 2: The "SRE" Refactor (Weekend Engineering - Primary)
*   **Persona:** Hsi (The SRE Builder)
*   **Scene:** Saturday morning. Hsi is exploring architectural patterns and decides to test if the "Calories" module can run independently.
*   **Action:** He opens the `infrastructure/terraform` directory. He changes a single variable `deploy_mode = "microservice"` for the Calories module. He runs `terraform apply`.
*   **System Response:** Terraform spins up a separate container for the Calories service and configures the networking/ingress to route traffic there instead of the monolith.
*   **Outcome:** The application continues to function identical to the user, but the architecture has fundamentally shifted.
*   **Key Requirement:** Strict Modular Monolith boundaries + Terraform-driven deployment logic.

### Journey 3: The Recruiter Scan (Portfolio Review - Secondary)
*   **Persona:** Technical Recruiter / Hiring Manager
*   **Scene:** Reviewing Hsi's CV for a Senior SRE role. They click the GitHub link for "TrackStack."
*   **Action:** They read the README, which clearly diagrams the "Modular Monolith" architecture. They navigate to `/cmd` and see the clean separation between `monolith/main.go` and `microservices/calories/main.go`.
*   **System Response:** The code structure visually confirms the architectural claims.
*   **Outcome:** The recruiter validates Hsi's claim of understanding complex system design.
*   **Key Requirement:** Clear directory structure and high-quality README/Documentation.

### Journey 4: The Data Fix (Admin/Ops - Edge Case)
*   **Persona:** Hsi (The Admin)
*   **Scene:** Hsi realizes the AI incorrectly logged "Coke" as "Diet Coke" yesterday.
*   **Action:** He opens the "Log" view, finds the entry, and taps "Edit." He corrects the entry.
*   **System Response:** The system updates the record and recalculates the "Yesterday" totals (though the UI focuses on Today).
*   **Outcome:** Data integrity is restored.
*   **Key Requirement:** CRUD capabilities for logged data.

## Domain-Specific Requirements (SRE & Infrastructure)

### Compliance & Standards (Self-Imposed)
*   **Observability:** Must support OpenTelemetry Context Propagation (W3C traceparent) to prove distributed tracing capability.
*   **Metrics:** Must expose Prometheus-compatible endpoints (`/metrics`) for all services.
*   **Logging:** Must use structured JSON logging (Slog) with injected `trace_id` for correlation.

### Technical Constraints
*   **Infrastructure-as-Code (IaC):** All infrastructure must be defined in Terraform. "ClickOps" in the AWS Console is forbidden.
*   **Secrets Management:** Secrets must be injected at runtime (via Env Vars or Secrets Manager) and never committed to Git.
*   **State Management:** Terraform state must be remote (S3 + DynamoDB locking) to simulate professional team environments.

### Deployment & Operations
*   **Immutable Artifacts:** The build process must produce a single Docker image artifact that is promoted across environments/deployments.
*   **CI/CD:** Deployments must be triggered via GitHub Actions pipelines, not local CLI commands.
*   **GitOps-Lite:** The `main` branch state should reflect the deployed state.

## Innovation & Novel Patterns

### The "Schrödinger's Binary" Pattern
TrackStack introduces a novel architectural pattern where the deployment topology is decoupled from the build artifact.
*   **Concept:** A single application binary can behave as a Monolith OR as a specific Microservice based solely on environmental configuration (Injection).
*   **Innovation:** Most systems require code changes or separate repos to split. TrackStack uses "Cellular Architecture" interfaces to make the boundary (In-Memory vs. gRPC) a runtime decision.
*   **Value:** Allows "Start Cheap" (Monolith/T4g.nano) and "Scale Infinite" (Microservices/Lambda) without rewriting business logic.

### "Pain-Driven" Observability
Moving beyond "Quantified Self" vanity metrics to "Drift Detection."
*   **Concept:** The system ignores historical analysis (which is often ignored by users) and focuses 100% on "Capacity Management" (what remains).
*   **Innovation:** Applying SRE "Error Budget" principles to personal diet and finance. You have a "Calorie Budget" and "Money Budget"; the system alerts only on burn-rate anomalies.

## Web App & API Specific Requirements

### Project-Type Overview
TrackStack is a hybrid **Progressive Web App (PWA)** backed by a **REST API**. It prioritizes mobile-first usage (add-to-homescreen) and offline resilience for data entry. The backend serves both the UI (via templates or API) and future microservices.

### Technical Architecture Considerations

#### Web App (PWA)
*   **Architecture:** Server-Side Rendering (SSR) with Go Templates (e.g., Templ) + HTMX for interactivity. avoiding complex JS frameworks (React/Vue) to keep the "Monolith" cohesive.
*   **Offline Strategy:** Service Worker cache for "Recents" data; Background Sync for submitting offline logs when connectivity returns.
*   **Responsive Design:** Mobile-first viewport optimization. Desktop view is secondary.
*   **Installation:** `manifest.json` for "Add to Homescreen" capability on iOS/Android.

#### API Backend
*   **Protocol:** REST (JSON) for external clients; gRPC allowed for internal module-to-module communication in the future.
*   **Authentication:** Session-based (Cookie) for PWA; API Key for "Gemini Ingest" or external scripts.
*   **Rate Limiting:** Per-user rate limiting to prevent API abuse/cost spikes.

### Implementation Considerations
*   **State:** Minimal client-side state; server is source of truth.
*   **Assets:** Static assets (CSS/Images) embedded in the Go binary for single-file deployment.

## Functional Requirements

### Ingest & Acquisition
*   **FR1:** The User can input meal data via natural language text (e.g., "I had a pizza").
*   **FR2:** The System can parse unstructured text into structured nutritional data (Name, Calories, Protein, Carbs, Fat) using an LLM provider.
*   **FR3:** The User can select a meal from a "Recent Meals" list to log it with a single action.
*   **FR4:** The User can manually input numeric values for calories/macros if the AI is unavailable.

### Tracking & State
*   **FR5:** The System can maintain a running total of calories and macros consumed for the current day.
*   **FR6:** The System can reset the daily total at a configured time (e.g., midnight).
*   **FR7:** The User can view the remaining budget (Target - Consumed) for the current day.
*   **FR8:** The User can edit or delete a previously logged entry for the current day to correct errors.

### Insights & Feedback
*   **FR9:** The User can view a "Dashboard" summary immediately upon application load.
*   **FR10:** The System can display visual indicators (e.g., progress bars) for Calorie and Macro limits.
*   **FR11:** The System can display alerts if the User exceeds their daily limit.

### System Management
*   **FR12:** The User can configure their daily Calorie and Macro targets.
*   **FR13:** The System can authenticate the User via a persistent session (Cookie).
*   **FR14:** The System can export all user data in a standard format (JSON/CSV) [Post-MVP, but good to list].

### Infrastructure Control
*   **FR15:** The System can be deployed as a single monolithic binary containing all modules.
*   **FR16:** The System can be deployed with specific modules disabled/enabled via environment variables.
*   **FR17:** The System can emit structured health checks (`/health`) for orchestration.

## Non-Functional Requirements

### Performance
*   **Response Time:** API P99 latency < 100ms for read operations.
*   **Startup Time:** Cold start < 500ms to prove Go efficiency and readiness for serverless/Lambda execution.
*   **Build Time:** CI pipeline (Commit to Artifact) < 2 minutes.

### Security
*   **Data Minimization:** No PII collected beyond email/username.
*   **Authentication:** Session-based for PWA; API Key for external integrations.
*   **Dependency Management:** Automated vulnerability scanning (Dependabot/Trivy) integrated into CI.

### Scalability
*   **Vertical Efficiency:** Must run on 128MB RAM (t4g.nano) in Monolith mode.
*   **Horizontal Scaling:** Stateless application tier allowing N replicas behind a load balancer.

### Maintainability (SRE Core)
*   **Code Quality:** Strict `golangci-lint` passing with zero warnings.
*   **Infrastructure-as-Code:** 100% Terraform coverage; no manual AWS console changes.
*   **Documentation:** Architecture Decision Records (ADRs) committed to the repository for every major architectural decision.
