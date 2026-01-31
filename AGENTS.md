# TrackStack Project Context

## Project Overview

**TrackStack** is a personal "Life Observability Platform" built as a Go Modular Monolith. It solves two problems simultaneously:
1. **Personal Drift Prevention:** Track three life domains with minimal friction:
   - **Calories:** Nutrition and macro tracking (MVP - currently in development)
   - **Money:** Financial expense tracking (upcoming)
   - **Heat:** Home inventory and resource tracking (upcoming)
2. **SRE Portfolio Showcase:** Demonstrates ability to deploy as monolith OR microservices via Terraform changes only

## Technology Stack

- **Language:** Go 1.22+
- **Web Framework:** Standard Library (`net/http`) with `http.ServeMux`
- **Router:** Native Go stdlib router (sufficient for modular monolith, allows clean module route registration)
- **Templating:** Templ (Type-safe HTML templates)
- **Interactivity:** HTMX (server-side rendered, dynamic updates)
- **Styling:** Tailwind CSS (dark mode first)
- **Database:** SQLite 3 (CGO enabled via `mattn/go-sqlite3`)
- **Query Builder:** SQLC (Type-safe SQL)
- **Migrations:** `golang-migrate`
- **Testing:** Go Test + Playwright-Go for E2E
- **Infrastructure:** Terraform + Docker Compose + Caddy
- **Deployment:** AWS ECS (t4g.nano target for <€7/month)

## Architecture Patterns

### The "Schrödinger's Binary" Pattern
The application can behave as a Monolith OR as specific Microservices based solely on environment variables:
- `DEPLOY_MODE=monolith` → All modules run internally
- `DEPLOY_MODE=microservice` + `MODULE=calories` → Only Calories module exposed

### Modular Monolith Structure
```
internal/
├── common/          # Shared infrastructure (db, server, middleware)
├── calories/        # Calorie tracking module (MVP)
├── money/           # Financial tracking (Post-MVP)
└── heat/            # Inventory tracking (Post-MVP)
```

**Critical Rule:** Modules NEVER import each other's stores/models. Communication ONLY via interfaces.

### Interface-Based Dependency Injection
- Define interfaces where they are *used* (Consumer Layer), not where implemented
- Use `NewService(deps...)` constructor pattern
- Runtime injection allows switching implementations (in-memory vs gRPC) without code changes

### Multi-Module Extensibility
Since TrackStack will grow from 1 module (calories) to 3 modules (calories, money, heat), all patterns must be designed for reuse:
- **Service/Store/Handler patterns** established in `calories/` serve as templates for `money/` and `heat/`
- **UI components** must be reusable across modules - shared layouts, consistent card patterns, modular dashboard sections
- **Interface boundaries** must support future gRPC implementations for microservice mode without code changes
- **Shared domain concepts** (daily tracking, "remaining" calculations, recents patterns) use consistent data structures
- **Database patterns** (SQLC queries, migration structure) are identical across all modules

## Critical Coding Rules

### Go Language Rules
- **JSON Tags:** ALWAYS use `json:"snake_case"` to align with SQLite columns
- **Error Handling:** Wrap all errors with context: `fmt.Errorf("failed to X: %w", err)`
- **Never return concrete types** from constructors if they implement an interface
- **NO Global State:** Do not use `init()` functions for side effects
- **NO Cross-Module Imports:** `internal/calories` cannot import `internal/money`

### HTMX/Templ Rules
- **OOB Swaps:** For non-fatal errors, return OOB Swap targeting `#toast-container`
- **Component Location:** Feature-specific `.templ` files live in `internal/{feature}/components/`
- **State Changes:** Prefer `hx-post` over `hx-get` for any state-changing action
- **Loading States:** Always include `hx-indicator` for async operations
- **NO Custom JavaScript:** Avoid `script` tags. Use HTMX or Alpine.js (`x-data`) only

### Database Rules
- **Strict Module Isolation:** A module must NEVER import Store/Models of another module
- **SQLC Types:** Use generated `db.CreateMealParams` structs. Don't manually map if avoidable
- **Transactions:** Use `RunTx` pattern in `internal/common/db` for atomic operations

### Testing Rules
- **Co-location:** Test files (`_test.go`) must live next to the code they test
- **Database Tests:** Use `testutil.NewTestDB()` (in-memory SQLite). Don't mock the driver
- **E2E Tests:** Use Playwright for user flows. Avoid fragile HTML regex parsing in unit tests

### Multi-Module Development Rules
Code the `calories` module as if `money` and `heat` already exist:
- **Pattern reusability:** Every service, store, and handler pattern must work for all three modules without modification
- **Shared UI library:** Extract reusable components to `internal/common/components/` (cards, forms, lists, charts)
- **Constructor extensibility:** Service constructors must accept interface dependencies only - never concrete types
- **Dashboard modularity:** Layout must accommodate multiple module cards side-by-side (calories + money + heat)
- **Route isolation:** Each module's HTTP routes are isolated under `/calories/`, `/money/`, `/heat/` prefixes
- **Consistent UX:** All three modules follow identical "Zero-Click Dashboard" and "Recents Carousel" patterns

### Git & Version Control Rules
- **NO COMMITS:** Do not create git commits, ever. The user is the final approver for all code changes.
- **NO PUSH:** Never push changes to remote repositories.

## Directory Structure

```
/home/hsi/src/trackstack/
├── AGENTS.md                 # This file
├── docs/                     # Planning & architecture docs
│   ├── prd.md
│   ├── architecture.md
│   ├── epics.md             # For reference only - not agent tasks
│   ├── ux-design.md
│   └── project-context.md
├── app/                      # Application code (Go Modular Monolith)
│   ├── cmd/
│   │   └── monolith/        # Single binary entry point
│   ├── internal/
│   │   ├── common/          # DB, server, middleware, shared components
│   │   │   ├── components/  # Reusable UI components (cards, forms, layouts)
│   │   │   ├── db/          # Database connection and transactions
│   │   │   └── server/      # HTTP server, middleware, routing
│   │   ├── calories/        # MVP module - nutritional tracking
│   │   │   ├── components/  # Calories-specific UI components
│   │   │   ├── handlers/    # HTTP handlers
│   │   │   ├── models/      # Domain models
│   │   │   └── store/       # Database access layer
│   │   ├── money/           # Financial tracking module (upcoming)
│   │   └── heat/            # Inventory tracking module (upcoming)
│   ├── migrations/
│   └── go.mod
├── infrastructure/          # Terraform configs
└── .opencode/              # OpenCode configuration (vanilla)
```

## UX Design Principles

### The "5-Second Rule"
Any interaction taking > 5 seconds will kill adoption. Prioritize speed over features.

### Core UX Patterns
1. **"Zero-Click Dashboard":** Open app → See "Calories Remaining" immediately (no navigation)
2. **"Recents Carousel":** 5-10 most recent meals as one-tap buttons (exploits habit formation)
3. **Traffic Light Colors:**
   - Green (>20% budget remaining) = Safe
   - Yellow (5-20% remaining) = Caution
   - Red (<5% or over) = Stop

### Key Interactions
- **Quick Log:** Tap Recent meal → Confirm → Dashboard updates (< 3 seconds)
- **AI Ingest:** Type "chicken sandwich" → AI returns macros → Confirm & Log (< 5 seconds)
- **The "Brunch Flow":** After logging, auto-scroll to Dashboard (see updated numbers), ready to log again

## Design System

- **Dark Mode First:** Background `#0a0a0a`, Text `#ffffff`
- **Typography:** System fonts only (no web font loading)
- **Dashboard Numbers:** 80px bold for calories, 60px for protein
- **Touch Targets:** Minimum 44px × 44px
- **Progress Bars:** Visual reinforcement showing consumed portion

## Module: Calories (MVP - First of Three)

**Note:** The `calories` module is currently in active development and serves as the blueprint for the `money` and `heat` modules. All patterns, UI components, and architectural decisions must be designed with multi-module reuse in mind.

### Domain Concepts
- **Target:** Daily calorie/protein goals (user-configurable)
- **Consumed:** Sum of today's logged meals
- **Remaining:** Target - Consumed (the key metric)
- **Meal:** Name + Calories + Protein + Carbs + Fat

### Key Features
- Instant Glance Dashboard (Remaining numbers)
- Recents Carousel (one-tap re-logging)
- AI Chat Input (Gemini API for meal parsing)
- Manual numeric quick-log (fallback when AI unavailable)
- Daily auto-reset at midnight

## Development Workflow

### Typical Session
1. You tell opencode what you want to implement (e.g., "Add the Recents Carousel to the dashboard")
2. Agent has full context from this file (tech stack, patterns, constraints)
3. You stay engaged, review code, guide decisions
4. Agent follows all documented patterns automatically

### Build Commands
Always use the Makefile for build operations:
```bash
cd app/
make build          # Build the application
```

**DO NOT** run raw `go run`, `go test`, or `go build` commands directly - always use `make` targets to ensure consistency.

### Running the App (via Makefile)
```bash
cd app/
make run            # Production build and run
```

## Additional Context

For detailed reference:
- **Architecture decisions:** @docs/architecture.md
- **Full PRD:** @docs/prd.md
- **UX specifications:** @docs/ux-design.md
- **Epics & stories:** @docs/epics.md (for your planning only)

---

**Last Updated:** 2026-01-31
**Project:** TrackStack - Personal Life Observability Platform
**Author:** Hsi
