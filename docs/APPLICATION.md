# Application Architecture: AI Directives

## 🗺️ Domains
- **Auth:** `src/pages/login.astro`, `src/server/auth`
- **Dashboard:** `src/modules/dashboard`
- **Calories:** `src/modules/calories`, `src/pages/calories`
- **Expenses:** `src/modules/expenses`, `src/pages/expenses`
- **Heat:** `src/modules/heat`, `src/pages/heat`

---

## 📂 Architecture Rules

### 1. `src/pages/` (Routing)
- **Role:** Astro pages plus thin frontend-facing API adapters.
- **Rule:** Map URLs to views or small adapter endpoints. Handle query params, cookies, headers, redirects, and Go API proxying.
- **Rule:** NO complex business logic here. Defer to `src/modules/`, `src/server/`, or the Go backend.

### 2. `src/components/` & `src/layouts/` (UI Elements)
- **Role:** "Dumb" presentational UI.
- **Rule:** NO data fetching. NO DB connections. NO direct mutations.
- **Rule:** Receive data exclusively via `Astro.props`.

### 3. `src/modules/` (Business Logic)
- **Role:** Domain-specific logic and queries.
- **Rule:** Keep domains isolated (e.g., Calories cannot query Expenses directly).
- **Rule:** New mutation logic should prefer the Go backend unless there is a documented reason to keep it in Astro.
- **Rule:** As domains migrate, Astro modules should become view-model helpers and page data loaders, not the source of mutation rules.

### 4. `src/server/` (Core Services)
- **Role:** Global services (`auth`, `db`).

---

## 🧠 Core Directives

### UI & Styling
- **Rule [Tailwind]:** ALWAYS use Tailwind CSS. NEVER use inline `style="..."` or external `.css` files unless Tailwind absolutely cannot handle dynamic calculations.
- **Rule [Content Separation]:** Data fetching occurs ONLY in the Astro Frontmatter (`---`). The template renders validated data without DB queries.

### Client Interactivity & Global State
*The problem: Astro is incredible for fast, zero-JS page loads. But interactive UI features (like modals or dropdowns) need global state management without pulling in a heavy framework like React.*

- **How it works (`ClientRuntime.astro`):** `apps/web/src/layouts/ClientRuntime.astro` acts as the single orchestrator for all vanilla JavaScript UI state across the app. It binds event listeners globally based on `data-*` attributes (e.g., `data-menu-trigger`) avoiding inline scripts spread across dozens of components. It manages accessibility (escape to close), click-outside behavior, and ensures UI state is correctly managed across view transitions.
- **Rule [No Frameworks]:** Do NOT add React, Vue, or heavy client-side frameworks. Any generic interactive UI behaviors you build must be written in Vanilla JS and integrated into the global `ClientRuntime.astro` using data attribute bindings.

### Mutations & Serverless Forms
*The problem: The production application runs behind an AWS Lambda Function URL secured by AWS IAM. A standard HTML form submitted by the browser (`<form method="POST">`) doesn't know how to sign the payload with AWS SigV4 credentials, resulting in a `403 Forbidden` error.*

- **How it works (`ApiFormHandler.astro`):** `apps/web/src/layouts/ApiFormHandler.astro` is a global interceptor script injected into the main layouts. 
  - **Interception:** It listens for submisssion events on any form carrying the `data-api-form` attribute and prevents the browser's default navigation behavior.
  - **Serialization:** It extracts the form inputs via the `FormData` API and serializes them into a clean JSON object.
  - **SigV4 Signing Prep:** AWS SigV4 requires the request payload to be hashed. The script uses the browser's native `window.crypto.subtle` API to calculate the `SHA-256` hash of the JSON body and injects it into the `x-amz-content-sha256` header.
  - **Fetching & UI Feedback:** Finally, it performs an AJAX `fetch` request using `window.signedFetch` (which is configured elsewhere in `AppShell` to handle the final AWS signing). Based on the HTTP response status, it either reloads the page, navigates to `data-redirect`, or automatically reveals a hidden error element defined by `data-error-target` containing the JSON error message sent back by the server.
- **Rule [API Mutations Only]:** UI components must never mutate data directly. Mutations must go through `src/pages/api/`, which may proxy to Go.
- **Rule [SigV4 Forms - CRITICAL]:** ANY form triggering a mutation must use the `data-api-form` attribute to be intercepted by `ApiFormHandler.astro`.
- **Rule [Go Boundary]:** For migrated domains, Astro API routes are adapters only. They normalize form payloads, forward cookies to Go, and translate Go responses into browser-friendly redirects or JSON.
- **Usage:**
  - `data-redirect="/path"` controls the success destination. 
  - `data-error-target="element-id"` injects the API error message softly on failure avoiding a full page crash.
  - *Example:* `<form method="POST" action="/api/calories" data-api-form data-redirect="/calories" data-error-target="error-id">`

### Database (Turso Distributed SQLite)
*The problem: We need to connect to 4 completely separate Turso databases simultaneously (Users, Calories, Expenses, Heat). Furthermore, local development uses hardcoded `.env` secrets, but AWS Lambda must securely fetch connection strings from AWS SSM Parameter Store at runtime to prevent leaks.*

- **How it works (`sqlite.ts`):** `apps/web/src/server/db/sqlite.ts` manages client instantiation and caching. If the environment variable starts with `/trackstack/` (which happens in production Terraform), the script uses the AWS SDK to pull the *real* values from SSM seamlessly.
- **Rule [No Raw Clients]:** NEVER instantiate a LibSQL client directly. Rely on `sqlite.ts` to orchestrate credential resolution.
- **Usage:** `import { getDb } from "@/server/db/sqlite"; const db = getDb("expenses");`
- **Rule [Frontend DB Boundary]:** Astro application runtime should not read or write domain databases directly. Page data, auth verification, and mutations should go through Go endpoints.
- **Allowed exceptions:** local tooling and migration support scripts may still talk to Turso directly (for example seeding or one-off maintenance), but that is not part of the app request path.

### Authentication & Middleware
*The problem: We need a secure, custom session system that protects both API routes and server-rendered pages without constantly prop-drilling a `userId` through every UI component and business logic function.*

- **How it works (`middleware.ts`):**
  - `apps/web/src/middleware.ts` is the bouncer. It intercepts every request. If you hit an API route without a session, it returns a 401. If you hit a Page, it redirects to `/login`.
  - Session validation for page requests is now delegated to the Go backend through `GET /api/auth/session`. If Go rotates the session, Astro forwards the returned `Set-Cookie` header back to the browser.
- **How it works (`currentUser.ts`):** To solve the prop-drilling problem, the middleware wraps the entire request lifecycle inside Node.js `AsyncLocalStorage`.
- **How it works (`login.ts` / `logout.ts`):** `apps/web/src/pages/api/auth/login.ts` and `apps/web/src/pages/api/auth/logout.ts` adapt browser form posts into JSON calls to Go and keep browser redirects in the frontend layer.
- **Current split:** Go is the source of truth for login, logout, and session verification. Astro still owns request-local auth context for SSR through `AsyncLocalStorage`.
- **Rule [AsyncLocalStorage]:** Do NOT pass `userId` as arguments to functions if the action is being performed by the currently logged-in user. 
- **Usage:** Deep inside any backend logic, securely grab the context: `import { getCurrentUserId } from "@/server/auth/currentUser"; const userId = getCurrentUserId();`

### Go Backend Boundary

- **Role:** `apps/server/internal/modules/**` owns domain rules, DTOs, and persistence contracts.
- **Rule:** New domain behavior should be added in Go first; Astro adapters should only normalize browser input, forward auth cookies, and map redirects/errors.
- **Rule:** Transport code in Go stays thin: parse request, call service, map status/error, serialize JSON. Redirect and browser-form behavior belong in frontend adapters, not Go handlers.
- **Rule:** Mutating Go endpoints should prefer a single JSON request contract; browser form posts should be normalized by frontend adapters before they reach Go.
- **Rule:** Expenses mutations and command-like actions (`close sheet`, `complete checklist`, template upserts/deletes) should call Go over HTTP rather than reimplementing logic in Astro routes.
- **Rule:** Delete-style Go endpoints should use an explicit URL identifier contract rather than accepting ids from multiple locations.
- **Rule:** Go route aliases should be temporary migration shims only; once Astro adapters call the canonical backend paths, remove the aliases from Go and OpenAPI.
- **Rule:** Astro forms and UI triggers should use the canonical migrated API paths directly once those paths are stable.
- **Rule:** When changing a Go endpoint contract, update the Astro adapter, transport tests, and e2e coverage together.
