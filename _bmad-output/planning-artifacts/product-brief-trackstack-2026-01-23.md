---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments: ['_bmad-output/analysis/brainstorming-session-2026-01-23.md', '_bmad-output/analysis/brainstorming-session-2026-01-20.md']
date: 2026-01-23
author: Hsi
---

# Product Brief: trackstack

## Executive Summary

**TrackStack** is a personal observability platform designed to solve two distinct problems simultaneously: the functional need to track personal "Drift" (Financial, Physical, Inventory) and the professional need to overcome Imposter Syndrome by demonstrating SRE mastery.

Unlike typical tracking apps which focus on historical logging, TrackStack focuses on **Drift Prevention** (managing "remainders" and "capacities") via a frictionless "Lazy UX" (AI Ingest, Recents Carousel). Simultaneously, the project serves as a "Full-Stack Thesis," proving the author's ability to engineer a **Modular Monolith** in Go that can be refactored into microservices via purely infrastructure changes (Terraform), offering a "Zero-Code-Change" deployment flexibility that serves as a powerful portfolio asset.

---

## Core Vision

### Problem Statement

The user, a Site Reliability Engineer (SRE), faces a dual challenge:
1.  **Personal Drift:** Without low-friction tracking, they "drift" into failure states (overspending, weight gain, running out of heating pellets).
2.  **Professional Imposter Syndrome:** Despite being employed as an SRE, they feel a need to prove deep competence across the entire engineering lifecycle (Dev -> Arch -> Infra -> Obs), particularly in Go and Cloud-Native patterns.

### Problem Impact

*   **Financial/Health:** "Stupid decisions" lead to overspending and weight gain.
*   **Inventory:** Heating pellet shortages risk "Winter Bollocks" (running out mid-season).
*   **Career:** Lack of a tangible, complex portfolio piece reinforces feelings of technical inadequacy.

### Why Existing Solutions Fall Short

*   **Commercial Apps (MyFitnessPal, Mint):** High friction for data entry, generic feature sets that don't track specific needs (e.g., pellet inventory), and offer zero professional value.
*   **Simple "Todo" Projects:** Build skills in isolation but fail to demonstrate complex architectural reasoning (Evolutionary Architecture, Infrastructure-as-Code boundaries).

### Proposed Solution

**TrackStack** is a "Life Observability Platform" built as a **Go Modular Monolith**.
*   **UX:** "Pain-Driven" interfaces including AI Voice/Text Ingest and "Recents Carousels" to minimize data entry friction.
*   **Architecture:** A cellular monolith where modules (Money, Calories, Heat, Body) are isolated by strict interfaces.
*   **Infrastructure:** A Terraform-managed AWS deployment that can deploy the app as a single binary (Monolith) or distinct services (Microservices) without changing application code.

### Key Differentiators

*   **The "Infrastructure-Driven" Refactor:** The unique ability to split the monolith into microservices via Terraform configuration alone, proving elite SRE skills.
*   **"Lazy" Observability:** Features designed for the absolute minimum effective dose of interaction (e.g., "3-Week Pulse" for heating, "Meal Totals" instead of ingredients).
*   **Dual-Value Prop:** It is both a daily utility for life management and a living resume for career advancement.

---

## Target Users

### Primary Users

#### **Hsi (The Dual-Role Creator)**
*   **Role 1: The "Lazy Human" (End User)**
    *   **Context:** 8 PM, tired after work, needs to log dinner/expenses quickly.
    *   **Motivation:** Prevent "Drift" (weight gain, overspending) with minimal effort.
    *   **Frustration:** High-friction inputs (manual forms) lead to abandonment.
    *   **Success Vision:** "I tap one button or say one sentence, and I know exactly how many calories/Euros I have left for the month."
*   **Role 2: The "SRE Builder" (Engineer)**
    *   **Context:** Weekends, learning sessions, interview prep.
    *   **Motivation:** Overcome Imposter Syndrome; prove mastery of Go, Terraform, and Architecture.
    *   **Frustration:** Projects that are "too simple" (Todo apps) and don't teach real-world complexity.
    *   **Success Vision:** "I can show a recruiter how I split a monolith into microservices just by changing a Terraform variable."

### Secondary Users

#### **The Partner (Independent User)**
*   **Context:** Managing her own personal expenses and data within the same system.
*   **Motivation:** Tracking her own financial/health goals independently.
*   **Interaction:** Logging personal expenses/meals using the same low-friction tools.
*   **Success Vision:** "I can use this tool for my own needs without my data mixing with Hsi's."
*   **Engineering Implication:** Standard Authentication/Authorization (User ID scoping) is required.

#### **The Technical Recruiter / Interviewer**
*   **Context:** Reviewing Hsi's portfolio for < 5 minutes.
*   **Motivation:** Quickly assess technical depth and architectural maturity.
*   **Interaction:** Browsing the GitHub ReadMe, Architecture Diagrams, and Infrastructure Code.
*   **Success Vision:** "This candidate understands the trade-offs between Monoliths and Microservices and has the code to prove it."

### User Journey (The "Feedback Loop")

1.  **Capture (Lazy Human):** Hsi eats a sandwich -> Opens PWA -> Taps "Chicken Sandwich" (Recents) OR Dictates "Chicken Sandwich" -> Done.
2.  **Process (System):** System logs data -> Updates "Remainders" -> Checks for Drift.
3.  **Review (Lazy Human):** Hsi sees "100 Kcal remaining" -> Decides not to eat a cookie.
4.  **Refactor (SRE Builder):** Hsi notices the "Ingest" module is slow -> Decides to refactor it into a separate microservice -> Updates Terraform -> Deploys -> Observes improvement via Grafana.

---

## Success Metrics

### User Success (The "Lazy Human")
*   **Friction Index:** Time-to-Log for recurring items < 5 seconds.
*   **Drift Prevention:** Zero "Surprise" events (e.g., unexpectedly running out of pellets or budget) per quarter.
*   **Engagement:** Active usage > 5 days/week (proving the "Lazy" UX works).

### Business Objectives (The "SRE Builder")
*   **Cost Efficiency:** Total AWS bill < €7/month (Goal: Ultra-lean infrastructure, potentially Serverless Containers vs. EC2 Spot).
*   **Portfolio Impact:** Creation of a "Technical Journey" blog series explaining the architecture, specifically the "Zero-Code-Change" deployment pattern.
*   **Engineering Velocity:** 100% success rate for "Infrastructure-Driven" deployment switches (Monolith <-> Microservice) without code modification.

### Key Performance Indicators
*   **Bill of Materials (BOM):** Monthly AWS cost tracking against the €7 limit.
*   **Code Coverage:** Core Logic test coverage > 80% (to justify the "Quality" claim).
*   **Latency:** API response time < 100ms for P99 interactions (proving Go efficiency).

---

## MVP Scope (The "Pilot Slice")

### Core Features
*   **Domain Focus:** Calories & Macros ONLY (Single Module).
*   **Architecture:** Go Modular Monolith structure (initialized with one module) + SQLite.
*   **Infrastructure:** Terraform-managed AWS environment (t4g.nano or equivalent cheap compute).
*   **UX (PWA):**
    *   **Recents Carousel:** 1-tap entry for frequent meals.
    *   **AI Ingest:** Text-to-JSON parsing via Gemini API ("Chicken sandwich and a coke").
    *   **Drift Dashboard:** Simple view of "Calories Remaining Today."

### Out of Scope for MVP
*   **Modules:** Money, Heating, and Body measurements (deferred to Phase 2).
*   **Microservices Deployment:** The code will support it, but the MVP deployment is Monolithic.
*   **Multi-User Support:** Hardcoded/Single-user auth initially to speed up the Pilot.
*   **Complex Observability:** Basic logging only; LGTM stack deferred.

### MVP Success Criteria
*   **Deployment:** API and PWA are live on AWS and accessible via public URL.
*   **Usage:** User successfully logs 3 days of meals without falling back to other tools.
*   **Cost:** Deployed infrastructure is projected to be < €7/month.

### Future Vision
*   **Phase 2 (Expansion):** Add Money, Heating, and Body modules using the established "Vertical Slice" pattern.
*   **Phase 3 (The Flex):** Demonstrate the "Infrastructure-Driven Refactor" by deploying the Money module as a separate microservice in a staging environment.
*   **Phase 4 (Observability):** Add full LGTM stack for deep tracing and drift alerts.
