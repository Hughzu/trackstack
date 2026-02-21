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
- **Role:** Astro SSR & API endpoints.
- **Rule:** Map URLs to views or APIs. Handle query params, cookies, and headers.
- **Rule:** NO complex business logic here. Defer to `src/modules/` or `src/server/`.

### 2. `src/components/` & `src/layouts/` (UI Elements)
- **Role:** "Dumb" presentational UI.
- **Rule:** NO data fetching. NO DB connections. NO direct mutations.
- **Rule:** Receive data exclusively via `Astro.props`.

### 3. `src/modules/` (Business Logic)
- **Role:** Domain-specific logic and queries.
- **Rule:** Keep domains isolated (e.g., Calories cannot query Expenses directly).

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
- **Rule [API Mutations Only]:** Database mutations MUST happen in `src/pages/api/`. UI components must NEVER mutate data directly.
- **Rule [SigV4 Forms - CRITICAL]:** ANY form triggering a mutation must use the `data-api-form` attribute to be intercepted by `ApiFormHandler.astro`.
- **Usage:**
  - `data-redirect="/path"` controls the success destination. 
  - `data-error-target="element-id"` injects the API error message softly on failure avoiding a full page crash.
  - *Example:* `<form method="POST" action="/api/calories" data-api-form data-redirect="/calories" data-error-target="error-id">`

### Database (Turso Distributed SQLite)
*The problem: We need to connect to 4 completely separate Turso databases simultaneously (Users, Calories, Expenses, Heat). Furthermore, local development uses hardcoded `.env` secrets, but AWS Lambda must securely fetch connection strings from AWS SSM Parameter Store at runtime to prevent leaks.*

- **How it works (`sqlite.ts`):** `apps/web/src/server/db/sqlite.ts` manages client instantiation and caching. If the environment variable starts with `/trackstack/` (which happens in production Terraform), the script uses the AWS SDK to pull the *real* values from SSM seamlessly.
- **Rule [No Raw Clients]:** NEVER instantiate a LibSQL client directly. Rely on `sqlite.ts` to orchestrate credential resolution.
- **Usage:** `import { getDb } from "@/server/db/sqlite"; const db = getDb("expenses");`

### Authentication & Middleware
*The problem: We need a secure, custom session system that protects both API routes and server-rendered pages without constantly prop-drilling a `userId` through every UI component and business logic function.*

- **How it works (`middleware.ts` & `session.ts`):** 
  - `apps/web/src/middleware.ts` is the bouncer. It intercepts every request. If you hit an API route without a session, it returns a 401. If you hit a Page, it redirects to `/login`.
  - It handles session security transparently: validating the cookie, updating the "last seen" time (touch), or rotating the cookie if it's getting old, using the logic defined in `session.ts` and the `users` Turso database.
- **How it works (`currentUser.ts`):** To solve the prop-drilling problem, the middleware wraps the entire request lifecycle inside Node.js `AsyncLocalStorage`.
- **Rule [AsyncLocalStorage]:** Do NOT pass `userId` as arguments to functions if the action is being performed by the currently logged-in user. 
- **Usage:** Deep inside any backend logic, securely grab the context: `import { getCurrentUserId } from "@/server/auth/currentUser"; const userId = getCurrentUserId();`
