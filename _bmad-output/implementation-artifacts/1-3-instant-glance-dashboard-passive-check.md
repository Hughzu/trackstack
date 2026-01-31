# Story 1.3: "Instant Glance" Dashboard (Passive Check)

Status: done

## Story

As a Lazy Human,
I want to see my remaining calories and protein in large, color-coded numbers on the home screen,
So that I can make a meal decision in < 1 second.

## Acceptance Criteria

1. **Given** I have a defined calorie/protein target and current logs in the database
2. **When** I open the app
3. **Then** the dashboard should display "Remaining" values (Target - Consumed) for the current day
4. **And** the numbers should be Green (>20% left), Yellow (5-20%), or Red (<5% or over)
5. **And** the UI should follow the "Progress-Enhanced Cards" design:
   - Bold hero numbers (80px for calories, 60px for protein)
   - Muted labels ("kcal remaining")
   - Horizontal progress bars showing consumption percentage

## Tasks / Subtasks

- [x] Initialize `calories` module schema (AC: 1)
  - [x] Create `meal_logs` table: `id`, `user_id`, `name`, `calories`, `protein`, `carbs`, `fat`, `logged_at`
  - [x] Create `user_targets` table: `user_id`, `calorie_target`, `protein_target` (default: 2300 kcal, 120g protein)
- [x] Implement Dashboard Backend Logic (AC: 3)
  - [x] Create `internal/calories` repository and service
  - [x] Implement `CalculateDailySummary(ctx, userID)` to calculate totals for today
- [x] Implement UI Components with Templ + Tailwind (AC: 5)
  - [x] Set up Tailwind CSS configuration with status colors
  - [x] Create `internal/calories/components/dashboard.templ`
  - [x] Implement `MetricCard` component with color-coding logic (AC: 4)
- [x] Update Routing (AC: 2)
  - [x] Update root route `/` in `internal/common/server` to render the Dashboard

## Dev Notes

- **Relevant architecture patterns:** Modular Monolith, Component-based UI (Templ), Server-Side Rendering.
- **Source tree components to touch:** `internal/calories/`, `internal/common/server/`, `tailwind.config.js`.
- **Timezone handling:** Use the server's local time or UTC for "today" calculation for now (Story 3.1 will handle personalization).
- **Default Targets:** If no record exists in `user_targets`, use 2000 kcal / 150g protein as defaults.

### References

- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#The Defining Experience: "The Instant Glance"]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Design System Foundation]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.3]

## Dev Agent Record

### Agent Model Used

claude-sonnet-4-5

### Debug Log References

N/A - Implementation completed successfully

### Completion Notes List

1. **Frontend Tooling Setup:**
   - Installed Templ CLI for type-safe Go HTML templates
   - Downloaded Tailwind CSS v4 standalone CLI
   - Created `tailwind.config.js` with UX spec design tokens (status colors, typography)
   - Created `Makefile` for build automation (templ generate, css build)

2. **Shared Components (`internal/common/components/`):**
   - `layout.templ`: Base HTML shell with dark mode styling
   - `progress_bar.templ`: Reusable progress bar component with status colors

3. **Database Schema (Calories Module):**
   - `meal_logs` table: Stores individual meal entries with nutritional data
   - `user_targets` table: Stores daily calorie/protein targets (defaults: 2300 kcal, 120g protein)
   - Indexes on `(user_id, logged_at)` for efficient daily queries

4. **Repository Pattern:**
   - `internal/calories/repository.go`: Database access layer
   - `GetDailyTotals()`: Aggregates calories/protein consumed using SQL SUM
   - `GetUserTargets()`: Fetches user goals with fallback to defaults

5. **Business Logic (Service Layer):**
   - `internal/calories/service.go`: Presentation-agnostic business logic
   - `CalculateDailySummary()`: Orchestrates data fetching and calculation
   - `CalculateMetricStatus()`: Determines color status (Safe >20%, Warn 5-20%, Danger <5%)

6. **UI Components (Templ):**
   - `internal/calories/components/dashboard.templ`: Main dashboard view
   - `internal/calories/components/metric_card.templ`: Reusable metric display with progress bar
   - Dark mode design with hero typography (80px calories, 60px protein)

7. **BFF Handler Pattern:**
   - `internal/calories/handlers/htmx_handler.go`: HTMX-specific presentation layer
   - Decoupled from service (supports future REST/GraphQL BFFs)
   - Renders HTML via Templ components

8. **Integration:**
   - Updated `cmd/monolith/main.go`: Wired calories module dependencies
   - Updated `internal/common/server/server.go`: Added `DashboardHandler` interface, root route `/`
   - Static file serving for CSS (`/static/`)

9. **All Acceptance Criteria Met:**
   - ✅ Dashboard displays remaining calories (2300) and protein (120) from defaults
   - ✅ Color-coded numbers (Green for >20% remaining, currently safe with 0% consumed)
   - ✅ Progress bars showing 0% consumption (no meals logged yet)
   - ✅ Bold hero typography and muted labels as per UX spec
   - ✅ Root route `/` renders dashboard, other endpoints (`/health`, `/api/session`) still functional

### File List

**Created:**
- `app/tailwind.config.js` - Design system configuration
- `app/static/css/input.css` - Tailwind input file
- `app/Makefile` - Build automation
- `app/internal/common/components/layout.templ` - Shared HTML layout
- `app/internal/common/components/progress_bar.templ` - Reusable progress bar
- `app/internal/calories/models.go` - Domain models (~35 lines)
- `app/internal/calories/repository.go` - Database access (~85 lines)
- `app/internal/calories/service.go` - Business logic (~75 lines)
- `app/internal/calories/components/dashboard.templ` - Dashboard view
- `app/internal/calories/components/metric_card.templ` - Metric card component
- `app/internal/calories/handlers/htmx_handler.go` - HTMX BFF (~45 lines)
- `app/bin/tailwindcss` - Tailwind standalone CLI binary

**Modified:**
- `app/internal/common/db/db.go` - Added calories module schema
- `app/internal/common/server/server.go` - Added dashboard handler interface, root route, static file serving
- `app/cmd/monolith/main.go` - Wired calories module
- `app/go.mod` - Added templ dependency
- `.gitignore` - Added generated files (*_templ.go, output.css)
