# Application Architecture: AI Directives

## 🗺️ Domains
- **Auth:** `src/pages/login.astro`, `src/layouts/AuthBootstrap.astro`
- **Dashboard:** `src/modules/dashboard`
- **Calories:** `src/modules/calories`, `src/pages/calories`
- **Expenses:** `src/modules/expenses`, `src/pages/expenses`
- **Heat:** `src/modules/heat`, `src/pages/heat`

---

## 📂 Architecture Rules

### 1. `src/pages/` (Routing)
- **Role:** Astro pages and static route shells.
- **Rule:** Map URLs to views and frontend navigation concerns. Astro pages must not act as backend adapters.
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

### 4. `src/server/`
- **Role:** Reserved for future frontend-only helpers if needed.
- **Rule:** Do not rebuild backend-like auth or data access layers here.

---

## 🧠 Core Directives

### UI & Styling
- **Rule [Tailwind]:** ALWAYS use Tailwind CSS. NEVER use inline `style="..."` or external `.css` files unless Tailwind absolutely cannot handle dynamic calculations.
- **Rule [Content Separation]:** Protected app data loads in the browser after auth readiness and should come from Go endpoints, not Astro frontmatter.

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
- **Rule [API Mutations Only]:** UI components must never mutate data directly. Browser mutations must target canonical `/api/*` endpoints owned by Go.
- **Deployment note:** In development, browser calls stay same-origin and use the Astro/Vite `/api` proxy. In production/static builds, browser calls target Go through `PUBLIC_API_BASE_URL`.
- **Rule [SigV4 Forms - CRITICAL]:** ANY form triggering a mutation must use the `data-api-form` attribute to be intercepted by `ApiFormHandler.astro`.
- **Rule [Go Boundary]:** The browser should call Go-owned `/api/*` endpoints directly. Astro does not own API adapter routes.
- **Rule [Timestamp Ownership]:** Calorie quick-add and manual meal logging should let the Go backend stamp the current UTC timestamp unless the product explicitly adds a backdating feature.
- **Usage:**
  - `data-redirect="/path"` controls the success destination. 
  - `data-error-target="element-id"` injects the API error message softly on failure avoiding a full page crash.
  - *Example:* `<form method="POST" action="/api/calories" data-api-form data-redirect="/calories" data-error-target="error-id">`

### Database (Turso Distributed SQLite)
*The problem: We need to connect to 4 completely separate Turso databases simultaneously (Users, Calories, Expenses, Heat). Furthermore, local development uses hardcoded `.env` secrets, but AWS Lambda must securely fetch connection strings from AWS SSM Parameter Store at runtime to prevent leaks.*

- **Rule [Frontend DB Boundary]:** Astro application runtime must not read or write domain databases directly. Page data, auth verification, and mutations go through Go endpoints.
- **Rule [Backend-Owned Tooling]:** Direct Turso access for seeding or maintenance belongs in backend-owned commands under `apps/server/cmd`, not inside the frontend app.

### Authentication
*The problem: We need a secure, custom session system that protects both API routes and server-rendered pages without constantly prop-drilling a `userId` through every UI component and business logic function.*

- **How it works (`AuthBootstrap.astro`):** `apps/web/src/layouts/AuthBootstrap.astro` performs a browser-side session bootstrap against `GET /api/auth/session` and publishes auth readiness into the client runtime.
- **Current page guard:** protected pages now rely on browser-side auth bootstrap against `GET /api/auth/session` and redirect to `/login` on the client when the session is missing.
- **Current auth flow:** login, logout, and session verification are all direct browser-to-Go interactions over `/api/auth/*`.
- **Session behavior:** Go owns a sliding session window. When authenticated traffic extends the server-side idle expiry, the API also refreshes the auth cookie expiry so the browser keeps sending the DB-backed session token.
- **Current split:** Go is the source of truth for login, logout, session verification, page data, and API contracts. The Astro app is now a static frontend shell plus client runtime.
- **Milestone reached:** the home, calories, expenses, and heat dashboards plus the calories and expenses settings pages now load their authenticated read models in the browser after auth bootstrap, so they no longer depend on `getCurrentUserId()` or SSR request-local auth context.
- **Overview behavior:** the home overview no longer waits on a single aggregated `/api/dashboard` response. It loads expenses, calories, and heat cards independently from their canonical module endpoints so one cold module path does not block the whole screen.

### Go Backend Boundary

- **Role:** `apps/server/internal/modules/**` owns domain rules, DTOs, and persistence contracts.
- **Rule:** New domain behavior should be added in Go first.
- **Rule:** Transport code in Go stays thin: parse request, call service, map status/error, serialize JSON.
- **Rule:** Mutating Go endpoints should prefer a single JSON request contract; the frontend runtime may still submit JSON from forms, but Astro pages do not adapt requests server-side.
- **Rule:** Expenses mutations and command-like actions (`close sheet`, `complete checklist`, template upserts/deletes) should call Go over HTTP rather than reimplementing logic in Astro routes.
- **Rule:** Delete-style Go endpoints should use an explicit URL identifier contract rather than accepting ids from multiple locations.
- **Rule:** Go route aliases should be temporary migration shims only; once the browser uses the canonical backend paths, remove the aliases from Go and OpenAPI.
- **Rule:** Astro forms and UI triggers should use the canonical migrated API paths directly once those paths are stable.
- **Rule:** When changing a Go endpoint contract, update the frontend client runtime, transport tests, and e2e coverage together.
