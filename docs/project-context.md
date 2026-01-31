---
project_name: 'trackstack'
user_name: 'Hsi'
date: '2026-01-25'
sections_completed: ['technology_stack', 'critical_rules']
status: 'complete'
optimized_for_llm: true
---

# Project Context for AI Agents

_This file contains critical rules and patterns that AI agents must follow when implementing code in this project. Focus on unobvious details that agents might otherwise miss._

---

## Technology Stack & Versions

-   **Language:** Go 1.22+
-   **Web Framework:** Standard Library (`net/http`) + `chi` or `echo` (defined in Architecture as StdLib + Middleware)
-   **Templating:** Templ (Type-safe HTML)
-   **Interactivity:** HTMX
-   **Styling:** Tailwind CSS
-   **Database:** SQLite 3 (CGO enabled via `mattn/go-sqlite3`)
-   **ORM/Query:** SQLC (Type-safe SQL)
-   **Migrations:** `golang-migrate`
-   **Testing:** Go Test + Playwright-Go
-   **Infrastructure:** Terraform + Docker Compose + Caddy

## Critical Implementation Rules

### Language-Specific Rules (Go)
-   **Interfaces:** Define interfaces where they are *used* (Consumer Layer), not where they are implemented.
-   **JSON Tags:** ALWAYS use `json:"snake_case"` to align with SQLite column naming.
-   **Error Handling:** Wrap all errors with context: `fmt.Errorf("failed to process meal: %w", err)`.
-   **Constructors:** Use `NewService(deps...)` pattern. Never return concrete types from constructors if they implement an interface.

### Framework-Specific Rules (Templ/HTMX)
-   **OOB Swaps:** For non-fatal errors (validation), return an OOB Swap targeting `#toast-container` rather than a full page error.
-   **Component Co-location:** Feature-specific `.templ` files should live in `internal/{feature}/components/`.
-   **Type Safety:** Use `templ.Component` as the return type for all UI render functions.
-   **HTMX Attributes:** Prefer `hx-post` over `hx-get` for any state-changing action. Always include `hx-indicator`.

### Database Rules (SQLite/SQLC)
-   **Strict Isolation:** A module (e.g., `calories`) must NEVER import the Store/Models of another module (`money`).
-   **SQLC Types:** Use generated `db.CreateMealParams` structs. Do not manually map if avoidable.
-   **Transactions:** Use the `RunTx` pattern in `internal/common/db` for atomic operations.

### Testing Rules
-   **Co-location:** Test files (`_test.go`) must live next to the code they test.
-   **Pattern:** Use `testutil.NewTestDB()` (in-memory SQLite) for storage tests. Do not mock the database driver.
-   **UI Testing:** Use Playwright for E2E flows; avoid fragile Regex parsing of HTML in unit tests.

### Critical Don't-Miss Rules
-   ❌ **NO Global State:** Do not use `init()` functions to register side effects.
-   ❌ **NO Cross-Module Imports:** `internal/calories` cannot import `internal/money`. Use Interfaces.
-   ❌ **NO Javascript:** Avoid `script` tags. Use HTMX or Alpine.js (`x-data`) only.

---

## Usage Guidelines

**For AI Agents:**
-   Read this file before implementing any code.
-   Follow ALL rules exactly as documented.
-   When in doubt, prefer the more restrictive option.

**For Humans:**
-   Keep this file lean and focused on agent needs.
-   Update when technology stack changes.
