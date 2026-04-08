# Application Architecture: AI Directives

## 🗺️ Domains
- **Auth:** `src/pages/login.astro`, `src/layouts/AuthBootstrap.astro`
- **Dashboard:** `src/modules/dashboard`
- **Calories:** `src/modules/calories`, `src/pages/calories`
- **Expenses:** `src/modules/expenses`, `src/pages/expenses`
- **Heat:** `src/modules/heat`, `src/pages/heat`

## Solid Migration Scaffold
- **Location:** `apps/web-next/src`
- **Rule:** `src/core/` owns shared config, theme runtime, auth token helpers, refresh/bootstrap auth state, route guards, shared formatters, and the typed `openapi-fetch` client.
- **Rule:** `src/components/ui/` is the only place where Tailwind utility composition should live in the Solid app.
- **Rule:** `src/features/{auth,dashboard,calories,expenses,heat}` owns route entry pages plus domain-local API wrappers.
- **Rule:** Solid feature files should compose UI primitives instead of spraying layout classes around like confetti.
- **Rule:** Repeated screen chrome in the Solid app belongs in semantic UI primitives first - think `Panel` header actions, list meta rows, compact counter pills, and action toggles - so feature files stay focused on domain mapping and state.
- **Rule:** Route entry files under `src/features/**` are orchestration layers only. They load data, hold route-level state, wire mutations, and compose feature components. They do not define page-local cards, ad-hoc dialogs/sheets, or long inline view sections.
- **Rule:** Feature-specific sections and cards belong in `src/features/<domain>/components/`. If a page starts growing subviews, extract them instead of letting `index.tsx` or sibling route files turn into a junk drawer.
- **Rule:** Feature mapping/helpers must not import prop types from `src/components/ui/`. Domain and feature logic should return feature-local view models or primitives, and the final UI adapter lives at the component boundary.
- **Rule:** Shared interaction patterns such as confirmation sheets, destructive action flows, skeleton structures, and reusable form layouts belong in `src/components/ui/` before they are repeated across features.

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
  - **Auth wiring:** On successful `POST /api/auth/login`, it stores the returned bearer token in browser storage. On successful `POST /api/auth/logout`, it clears that client token before redirect or reload.
  - **Fetching & UI Feedback:** Finally, it performs an AJAX `fetch` request using `window.signedFetch` (created by the runtime when needed). Based on the HTTP response status, it either reloads the page, navigates to `data-redirect`, or automatically reveals a hidden error element defined by `data-error-target` containing the JSON error message sent back by the server.
- **Rule [API Mutations Only]:** UI components must never mutate data directly. Browser mutations must target canonical `/api/*` endpoints owned by Go.
- **Deployment note:** In development, browser calls stay same-origin and use the Astro/Vite `/api` proxy. In production/static builds, browser calls target Go through `PUBLIC_API_BASE_URL`.
- **Rule [Auth Header]:** Protected browser mutations and confirm-modal requests attach the bearer token through `X-Trackstack-Authorization` so the deployed CloudFront -> Lambda path stays compatible with SigV4 origin signing.
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

- **How it works (`AuthBootstrap.astro`):** `apps/web/src/layouts/AuthBootstrap.astro` reads the stored bearer token, performs a browser-side session bootstrap against `GET /api/auth/session`, and publishes auth readiness into the client runtime. In the CloudFront -> Lambda Function URL deployment, the browser token travels in `X-Trackstack-Authorization` so CloudFront can keep using `Authorization` for SigV4 origin signing.
- **Current page guard:** protected pages wait for browser-side auth bootstrap and redirect to `/login` on the client when the session is missing. Public-only pages redirect back to `/` once a valid session is confirmed.
- **Current auth flow:** login, refresh, logout, and session verification are all direct browser-to-Go interactions over `/api/auth/*`. Successful login stores the short-lived bearer token client-side and gets a refresh token through an `HttpOnly` cookie; refresh rotates that cookie and returns a fresh access token.
- **Session behavior:** protected routes still use bearer auth only, but session lifecycle is no longer fake-stateless nonsense. Full page reloads and Astro navigations revalidate the stored token against `GET /api/auth/session`, and browser retries can recover through `POST /api/auth/refresh`.
- **Solid auth flow:** `apps/web-next/src/core/auth/store.ts` bootstraps auth once after the SPA mounts, keeps the access JWT in `sessionStorage`, lets route guards in `apps/web-next/src/core/auth/guards.tsx` own the checking/loading state, and logs out from the shared `AppShell` UI.
- **Solid route UX:** `apps/web-next/src/components/ui/AppRoot.tsx` keeps the shared `AppShell` mounted at the router root and swaps only the `main` content with a skeleton fallback while lazy route chunks resolve, so first-visit navigations stop flashing the whole shell.
- **Solid navigation loading:** client-side route changes in `apps/web-next/src/components/ui/AppRoot.tsx` now use the router pending signal to replace stale page content with a shared `RouteStatus` skeleton during tab switches, so navigation loading behaves consistently across dashboard, expenses, calories, and heat.
- **Solid refresh contract:** `apps/web-next/src/core/api/client.ts` retries protected requests once after `401` by calling `POST /api/auth/refresh` with `credentials: include`; if Go does not expose that route yet, the Solid app falls back to a clean guest state instead of pretending everything is fine.
- **Solid logout behavior:** an explicit logout writes a browser-side logout marker so the SPA does not silently rehydrate itself from a still-valid refresh cookie during the same session window.
- **Solid dashboard wiring:** `apps/web-next/src/features/dashboard/index.tsx` now uses three independent `createResource` reads against the monolith contracts exposed through `apps/web-next/src/features/{expenses,calories,heat}/api/client.ts`, so each card keeps its own skeleton/error state instead of faking a single mocked payload.
- **Solid expenses flow:** `/expenses` now owns live dashboard mutations for close-sheet, checklist completion, and entry deletion, `/expenses/new` posts real expense entries through the typed OpenAPI client, and `/expenses/settings` now reads and mutates live income, ratio, checklist-template, and recurring-template data from the canonical expenses endpoints.
- **Solid expenses safeguard:** the `/expenses` close-month action now requires an explicit in-app confirmation before posting `POST /api/expenses/sheet/close`, so the UI stops treating period rollover like a casual click.
- **Solid expenses delete UX:** expense history deletions now reuse the same confirmation-sheet pattern before calling `DELETE /api/expenses/entries/{id}`, so destructive actions behave consistently instead of firing on the first tap.
- **Current split:** Go is the source of truth for login, logout, session verification, page data, and API contracts. The Astro app is now a static frontend shell plus client runtime.
- **Milestone reached:** the home, calories, expenses, and heat dashboards plus the calories and expenses settings pages now load their authenticated read models in the browser after auth bootstrap, so they no longer depend on `getCurrentUserId()` or SSR request-local auth context.
- **Overview behavior:** the home overview no longer waits on a single aggregated `/api/dashboard` response. It loads expenses, calories, and heat cards independently from their canonical module endpoints so one cold module path does not block the whole screen.
- **Shared browser auth transport:** `ApiFormHandler.astro`, `ClientRuntime.astro`, and the client-loaded dashboard/settings modules attach `X-Trackstack-Authorization: Bearer <jwt>` on protected browser calls.

### Go Backend Boundary

- **Role:** Go backend business logic now centers on context-local boundaries under `apps/server/internal/contexts/{auth,users,heat,calories,expenses}/**`.
- **Role:** Runtime assembly in `apps/server` lives under `apps/server/internal/app/monolithapi` and is shared by both `apps/server/cmd/monolith-api` and `apps/server/cmd/lambda`.
- **Rule:** New domain behavior should be added in Go first.
- **Rule:** Transport code in Go stays thin: parse request, call service, map status/error, serialize JSON.
- **Rule:** Mutating Go endpoints should prefer a single JSON request contract; the frontend runtime may still submit JSON from forms, but Astro pages do not adapt requests server-side.
- **Rule:** Expenses mutations and command-like actions (`close sheet`, `complete checklist`, template upserts/deletes) should call Go over HTTP rather than reimplementing logic in Astro routes.
- **Rule:** Delete-style Go endpoints should use an explicit URL identifier contract rather than accepting ids from multiple locations.
- **Rule:** Go route aliases should be temporary migration shims only; once the browser uses the canonical backend paths, remove the aliases from Go and OpenAPI.
- **Rule:** Astro forms and UI triggers should use the canonical migrated API paths directly once those paths are stable.
- **Rule:** When changing a Go endpoint contract, update the frontend client runtime, transport tests, and e2e coverage together.
- **Heat transition note:** heat now owns its inbound HTTP adapter under `apps/server/internal/contexts/heat/adapters/inbound/http`, and the browser delete flow now targets canonical `DELETE /api/heat/refills/{id}` while the query-param delete route remains as a temporary compatibility alias.
- **Heat facade note:** runtime assembly and Go transport depend on the context-local heat application facade in `apps/server/internal/contexts/heat/application/service.go`; the legacy backend `internal/modules/heat` package has been removed.
- **Heat rebuild note:** `apps/server/internal/contexts/heat/**` now exposes `GET /api/heat/refills`, `POST /api/heat/refills`, and canonical `DELETE /api/heat/refills/{id}` with bearer-auth enforced by shared middleware at the runtime boundary.
- **Calories rebuild note:** `apps/server/internal/contexts/calories/**` now exposes `GET /api/calories/dashboard`, `GET /api/calories/target`, `POST /api/calories/target`, `POST /api/calories/log`, and canonical `DELETE /api/calories/logs/{id}`. Unlike the legacy calories contract, the rebuild intentionally uses explicit nutrient field names such as `proteinGrams`, `carbGrams`, `fatGrams`, `targetCalories`, and `targetProteinGrams`. Any accepted contract break must be recorded in `docs/BACKEND_BREAKING_CHANGES.md`.
- **Expenses rebuild note:** `apps/server/internal/contexts/expenses/**` now exposes `GET /api/expenses/settings`, `POST /api/expenses/settings`, `GET /api/expenses/sheet/current`, `POST /api/expenses/entries`, `POST /api/expenses/checklists`, `POST /api/expenses/checklists/complete`, `POST /api/expenses/recurring`, and `POST /api/expenses/sheet/close` with the existing frontend JSON contract for settings, dashboard, entries, and templates.
- **Expenses delete contract note:** the rebuild intentionally removes legacy query-param deletes and uses canonical resource routes instead: `DELETE /api/expenses/entries/{id}`, `DELETE /api/expenses/checklists/{id}`, and `DELETE /api/expenses/recurring/{id}`. Accepted contract breaks must be recorded in `docs/BACKEND_BREAKING_CHANGES.md`.
- **Expenses dashboard note:** the expenses dashboard use case now depends on a snapshot-style read port in Go so the application layer does not orchestrate several separate read calls for one page against Turso.
- **Backend-first rebuild verification note:** `apps/server/scripts/e2e.sh` is the current backend smoke harness and should be kept aligned with auth and route contract changes.
