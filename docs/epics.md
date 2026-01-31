---
stepsCompleted: [1, 2, 3, 4]
inputDocuments: ["_bmad-output/planning-artifacts/prd.md", "_bmad-output/planning-artifacts/architecture.md", "_bmad-output/planning-artifacts/ux-design-specification.md"]
---

# trackstack - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for trackstack, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

FR1: The User can input meal data via natural language text (e.g., "I had a pizza").
FR2: The System can parse unstructured text into structured nutritional data (Name, Calories, Protein, Carbs, Fat) using an LLM provider.
FR3: The User can select a meal from a "Recent Meals" list to log it with a single action.
FR4: The User can manually input numeric values for calories/macros if the AI is unavailable.
FR5: The System can maintain a running total of calories and macros consumed for the current day.
FR6: The System can reset the daily total at a configured time (e.g., midnight).
FR7: The User can view the remaining budget (Target - Consumed) for the current day.
FR8: The User can edit or delete a previously logged entry for the current day to correct errors.
FR9: The User can view a "Dashboard" summary immediately upon application load.
FR10: The System can display visual indicators (e.g., progress bars) for Calorie and Macro limits.
FR11: The System can display alerts if the User exceeds their daily limit.
FR12: The User can configure their daily Calorie and Macro targets.
FR13: The System can authenticate the User via a persistent session (Cookie).
FR14: The System can export all user data in a standard format (JSON/CSV) [Post-MVP, but good to list].
FR15: The System can be deployed as a single monolithic binary containing all modules.
FR16: The System can be deployed with specific modules disabled/enabled via environment variables.
FR17: The System can emit structured health checks (/health) for orchestration.

### NonFunctional Requirements

NFR1: API P99 latency < 100ms for read operations.
NFR2: Cold start < 500ms to prove Go efficiency and readiness for serverless/Lambda execution.
NFR3: CI pipeline (Commit to Artifact) < 2 minutes.
NFR4: No PII collected beyond email/username.
NFR5: Session-based authentication for PWA; API Key for external integrations.
NFR6: Automated vulnerability scanning (Dependabot/Trivy) integrated into CI.
NFR7: Must run on 128MB RAM (t4g.nano) in Monolith mode.
NFR8: Stateless application tier allowing N replicas behind a load balancer.
NFR9: Strict golangci-lint passing with zero warnings.
NFR10: 100% Terraform coverage; no manual AWS console changes.
NFR11: Architecture Decision Records (ADRs) committed to the repository for every major architectural decision.

### Additional Requirements

- Schrödinger's Binary Pattern: Single binary, topology decoupled from artifact, runtime injection for monolith vs microservices.
- Interface-based Dependency Injection for strict modular boundaries.
- SQLite CGO with t4g.nano + Caddy for <€7 budget.
- OpenTelemetry Context Propagation and Prometheus-compatible endpoints.
- Structured JSON logging (Slog) with trace_id.
- Mobile PWA with manifest.json and Service Worker cache for "Recents".
- "5-Second Rule" UX: Instant Glance Dashboard (remaining calories/protein).
- "Recents Carousel" and integrated AI chat for logging.
- Auto-scroll to dashboard after logging (Brunch Flow).

### FR Coverage Map

FR1: Epic 2 - Natural language input
FR2: Epic 2 - AI parsing to nutritional data
FR3: Epic 1 - Recents Carousel logging
FR4: Epic 1 - Manual numeric input (essential fallback)
FR5: Epic 1 - Daily total maintenance
FR6: Epic 3 - Daily reset
FR7: Epic 1 - "What Remains" calculation
FR8: Epic 3 - Edit/Delete entries
FR9: Epic 1 - Instant Glance Dashboard
FR10: Epic 1 - Visual progress bars
FR11: Epic 3 - Over-budget alerts
FR12: Epic 3 - Configure targets
FR13: Epic 1 - Persistent Session
FR15: Epic 1 - Monolithic deployment
FR16: Epic 4 - Module toggling / Microservice mode
FR17: Epic 1 - Health checks

## Epic List

### Epic 1: The Core "Lazy" Experience
Enable users to see their "Calorie Drift" immediately and log repeat meals with a single tap. This establishes the "Instant Glance" dashboard and the "Muscle Memory" carousel.
**FRs covered:** FR3, FR4, FR5, FR7, FR9, FR10, FR13, FR15, FR17.

### Epic 2: Conversational Ingest
Integrate AI to allow logging of new meals using natural language, eliminating the need for manual macro searching.
**FRs covered:** FR1, FR2.

### Epic 3: Goal Personalization & Correction
Allow users to tune the system to their specific targets and fix logging errors to ensure "Drift" detection is accurate.
**FRs covered:** FR6, FR8, FR11, FR12.

#### Epic 4: Cloud-Native Portfolio Flex
Implement the "Schrödinger's Binary" pattern and full SRE observability stack to prove technical mastery.
**FRs covered:** FR16.

## Epic 1: The Core "Lazy" Experience

Enable users to see their "Calorie Drift" immediately and log repeat meals with a single tap. This establishes the "Instant Glance" dashboard and the "Muscle Memory" carousel.

### Story 1.1: Project Initialization & Monolith Skeleton

As a Developer,
I want to initialize the repository using the Go-HTMX starter template and set up the modular monolith directory structure,
So that I have a clean foundation that supports future module isolation.

**Acceptance Criteria:**

**Given** a clean repository
**When** I run the initialization script and set up `internal/common`, `internal/calories`, and `cmd/monolith`
**Then** the project should compile and run a basic health check at `/health`
**And** the `golangci-lint` check should pass with zero warnings.

### Story 1.2: Persistent User Session & SQLite Foundation

As a User,
I want the system to remember who I am across sessions using a persistent cookie,
So that I don't have to log in every time I open the app.

**Acceptance Criteria:**

**Given** a running server and a defined SQLite schema for users
**When** I access the application
**Then** the system should issue/recognize a session cookie
**And** my user context should be available to the backend for data scoping.

### Story 1.3: "Instant Glance" Dashboard (Passive Check)

As a Lazy Human,
I want to see my remaining calories and protein in large, color-coded numbers on the home screen,
So that I can make a meal decision in < 1 second.

**Acceptance Criteria:**

**Given** I have a defined calorie/protein target and current logs in the database
**When** I open the app
**Then** the dashboard should display "Remaining" values (Target - Consumed)
**And** the numbers should be Green (>20% left), Yellow (5-20%), or Red (<5% or over).

### Story 1.4: Manual Numeric Quick-Log

As a User,
I want to manually enter a numeric calorie/protein value when I'm in a rush,
So that I can track my intake even if I don't want to describe the meal.

**Acceptance Criteria:**

**Given** I am on the dashboard
**When** I enter numbers into the "Quick Add" fields and submit
**Then** the dashboard numbers should update instantly via HTMX without a page reload.

### Story 1.5: "Muscle Memory" Recents Carousel

As a User,
I want to see a list of my 5-10 most recent meals and log one with a single tap,
So that I can track recurring meals in < 3 seconds.

**Acceptance Criteria:**

**Given** I have previous meal logs in the database
**When** I tap a meal in the Recents list
**Then** the system should create a new log entry based on that meal's data
**And** the dashboard should auto-scroll to the top to show the updated "Remaining" numbers.

## Epic 2: Conversational Ingest

Integrate AI to allow logging of new meals using natural language, eliminating the need for manual macro searching.

### Story 2.1: Gemini API Integration & Structured Output

As a Developer,
I want to integrate the Gemini API with a system prompt that enforces a JSON response (Name, Calories, Macros),
So that the application can reliably convert unstructured text into actionable nutritional data.

**Acceptance Criteria:**

**Given** a text string like "I had a chicken sandwich"
**When** the system sends the request to Gemini
**Then** the response should be a valid JSON object containing `name`, `calories`, `protein`, `carbs`, and `fat`
**And** the system should handle API timeouts or errors gracefully without crashing the UI.

### Story 2.2: Fixed AI Chat Interface (Frontend)

As a User,
I want a persistent chat input bar at the bottom of my screen,
So that I can type a meal description at any time without navigating away from the dashboard.

**Acceptance Criteria:**

**Given** any view in the app
**When** I look at the bottom of the screen
**Then** I should see a rounded text input with a "Send" icon
**And** the input should remain fixed even when I scroll through the dashboard or recents.

### Story 2.3: Refinement Loop & "Confirm & Log" Action

As a User,
I want to see the AI's estimate and have the option to refine it (e.g., "add bacon") before logging it,
So that I can ensure my calorie tracking is accurate.

**Acceptance Criteria:**

**Given** I have sent a meal description
**When** the AI returns the nutritional data
**Then** the system should display a "Confirm & Log" button and allow me to send a follow-up message to adjust the values
**And** clicking "Confirm" should save the meal and trigger the "Brunch Flow" (auto-scroll to dashboard + update numbers).

## Epic 3: Goal Personalization & Correction

Allow users to tune the system to their specific targets and fix logging errors to ensure "Drift" detection is accurate.

### Story 3.1: Daily Target Configuration

As a User,
I want to set my daily calorie and protein targets in a settings view,
So that the dashboard correctly reflects my personal health goals.

**Acceptance Criteria:**

**Given** I am in the Settings view
**When** I enter new numeric targets and save
**Then** the system should update my user profile in the database
**And** the dashboard should immediately reflect these new targets in its "Remaining" calculations.

### Story 3.2: Log Management (Edit/Delete)

As a User,
I want to view my logs for the current day and edit or delete them if I made a mistake,
So that my daily totals remain accurate.

**Acceptance Criteria:**

**Given** I have logged meals today
**When** I view the daily log list and select "Edit" or "Delete"
**Then** the system should update the database record
**And** the dashboard numbers should update to reflect the change.

### Story 3.3: Automated Daily Reset

As a Lazy Human,
I want my consumed totals to reset to zero at midnight every day,
So that I don't have to manually clear my data every morning.

**Acceptance Criteria:**

**Given** the current time passes midnight
**When** I open the app the next morning
**Then** the dashboard should show the full daily budget as "Remaining"
**And** the "Consumables" list should be empty for the new day.

### Story 3.4: Budget Limit Alerts

As a User,
I want to see a visual alert when I am very close to or over my budget,
So that I am consciously aware that I've reached my limit.

**Acceptance Criteria:**

**Given** my "Remaining" calories fall below 5% of my target
**When** the dashboard renders
**Then** the numbers should turn Red
**And** a subtle warning message should appear on the dashboard.

## Epic 4: Cloud-Native Portfolio Flex

Implement the "Schrödinger's Binary" pattern and full SRE observability stack to prove technical mastery.

### Story 4.1: "Schrödinger's Binary" Runtime Injection

As a Developer,
I want the application to choose its mode (Monolith vs. Microservice) based on environment variables and interface-based dependency injection,
So that I can change the deployment topology without recompiling the code.

**Acceptance Criteria:**

**Given** a single binary
**When** I start it with `DEPLOY_MODE=monolith`, it starts all modules internally
**When** I start it with `DEPLOY_MODE=microservice` and `MODULE=calories`, it only exposes the Calories gRPC/HTTP endpoints
**Then** the core business logic remains identical in both modes.

### Story 4.2: OpenTelemetry & Distributed Tracing

As a Recruiter/SRE,
I want to see complete request traces across the system, even when split into microservices,
So that I can verify Hsi's mastery of distributed observability.

**Acceptance Criteria:**

**Given** a request entering the system
**When** the request travels through modules (or across network boundaries in microservice mode)
**Then** the system should propagate the W3C traceparent context
**And** structured logs should include the `trace_id`.

### Story 4.3: Terraform-Driven Topology Refactor

As an SRE,
I want to use a single Terraform variable to toggle the entire AWS infrastructure between Monolith and Microservices,
So that I can demonstrate infrastructure-driven architectural agility.

**Acceptance Criteria:**

**Given** the `infrastructure/terraform` directory
**When** I change `is_microservice = true` and run `terraform apply`
**Then** AWS should spin up separate ECS tasks for each module and update the Load Balancer routing
**And** the PWA should continue to function seamlessly for the user.




