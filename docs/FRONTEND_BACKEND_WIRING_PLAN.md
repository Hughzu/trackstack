# Frontend / Backend Wiring Plan

Date: 2026-03-23

This document explains what must be migrated so `apps/web` can use `apps/server-next` as its backend cleanly and safely.

It is not a backend target document. It is a wiring and migration plan for the browser/runtime boundary.

## Goal

Make `apps/web` work against `apps/server-next` with:

- bearer-auth instead of cookie-auth
- canonical `server-next` route shapes
- canonical `server-next` request payloads
- canonical `server-next` response payloads
- repeatable verification across backend smoke tests and browser tests

## Current Situation

The backend is now far enough along to begin the frontend migration:

- `apps/server-next` has `auth`, `users`, `heat`, `calories`, and `expenses`
- it exposes local HTTP and Lambda runtimes
- it has CORS and `/openapi.yaml`
- it has backend-owned tooling via `cmd/seed-user`
- it has backend-first smoke coverage via `apps/server-next/scripts/e2e.sh`

The main missing work is now on the frontend wiring side.

## Why The Frontend Cannot Just Point At `server-next`

The current Astro frontend still assumes the old auth and calories contracts.

Main mismatch areas:

1. auth transport
2. calories request field names
3. calories response field names
4. browser fetch/form behavior for protected requests
5. test expectations still centered on the old browser contract

Intentional backend contract breaks are tracked in `docs/BACKEND_BREAKING_CHANGES.md` and must be treated as migration inputs, not bugs.

## Backend Facts The Frontend Must Wire To

### Auth

`apps/server-next` now uses bearer auth only.

- `POST /api/auth/login` returns JSON with `accessToken`, `tokenType`, `expiresAt`, and `userId`
- `GET /api/auth/session` requires `Authorization: Bearer <jwt>`
- `POST /api/auth/logout` returns `204` and is stateless
- protected routes under `/api/calories/*`, `/api/expenses/*`, and `/api/heat/*` require `Authorization: Bearer <jwt>`

This means the frontend must own token storage, token attachment, and token clearing.

### Route Surface

The browser should wire to these canonical backend endpoints:

- `/api/auth/login`
- `/api/auth/logout`
- `/api/auth/session`
- `/api/heat/dashboard`
- `/api/heat/refills`
- `/api/heat/refills/{id}`
- `/api/calories/dashboard`
- `/api/calories/target`
- `/api/calories/log`
- `/api/calories/logs/{id}`
- `/api/expenses/settings`
- `/api/expenses/sheet/current`
- `/api/expenses/entries`
- `/api/expenses/entries/{id}`
- `/api/expenses/checklists`
- `/api/expenses/checklists/{id}`
- `/api/expenses/checklists/complete`
- `/api/expenses/recurring`
- `/api/expenses/recurring/{id}`
- `/api/expenses/sheet/close`

Important:

- `/api/dashboard` is intentionally removed in `server-next`
- the home overview must keep loading module cards independently

### Calories Contract Breaks

The frontend must switch from legacy calories names to explicit ones.

Request payload changes:

- `protein` -> `proteinGrams`
- `carbs` -> `carbGrams`
- `fat` -> `fatGrams`
- `targetKcal` -> `targetCalories`
- `targetProtein` -> `targetProteinGrams`
- `targetCarbs` -> `targetCarbGrams`
- `targetFat` -> `targetFatGrams`

Response payload changes:

- `proteinG` -> `proteinGrams`
- `carbsG` -> `carbGrams`
- `fatG` -> `fatGrams`
- `targetKcal` -> `targetCalories`
- `targetProteinG` -> `targetProteinGrams`
- `targetCarbsG` -> `targetCarbGrams`
- `targetFatG` -> `targetFatGrams`
- `consumed` -> `consumedCalories`
- `protein` -> `proteinGrams`
- `carbs` -> `carbGrams`
- `fat` -> `fatGrams`

### Heat And Expenses Contract Breaks

Delete routes are canonical path-based resources only:

- heat: `DELETE /api/heat/refills/{id}`
- expenses entries: `DELETE /api/expenses/entries/{id}`
- expenses checklists: `DELETE /api/expenses/checklists/{id}`
- expenses recurring: `DELETE /api/expenses/recurring/{id}`

The current frontend already looks close on those route shapes.

## Frontend Areas That Must Change

## 1. Auth Runtime

This is the first migration slice and the main blocker.

Files involved now:

- `apps/web/src/layouts/AuthBootstrap.astro`
- `apps/web/src/layouts/ApiFormHandler.astro`
- `apps/web/src/pages/login.astro`
- `apps/web/src/components/TopBar.astro`

What must change:

- store the token returned by `POST /api/auth/login`
- attach `Authorization: Bearer <jwt>` to protected fetches
- attach `Authorization: Bearer <jwt>` to protected form submissions handled by `ApiFormHandler`
- remove reliance on `credentials: "include"` as the auth mechanism
- clear the stored token on logout
- have auth bootstrap read session state using bearer auth instead of cookie auth

What the frontend needs to decide explicitly:

- where the token is stored in the browser
- how auth readiness is restored after page loads/navigation
- how logout behaves when the backend is stateless

Recommended direction:

- keep token handling in the shared browser runtime/layout layer, not scattered across domain modules
- centralize token read/write/clear in one client helper
- let `AuthBootstrap` remain the source of truth for browser auth readiness

## 2. Shared Browser Fetch Layer

The shared browser request path must be made auth-aware.

Files involved now:

- `apps/web/src/layouts/ApiFormHandler.astro`
- `apps/web/src/config/api.ts`
- any module using direct `fetch(resolveBrowserApiUrl(...))`

What must change:

- add a shared way to include bearer headers in browser requests
- keep SigV4/body hashing behavior intact for deployments that still need signed requests
- keep same-origin local dev behavior working
- preserve redirect/error handling already built around `data-api-form`

Target rule:

- the frontend should have one shared way to make authenticated browser API calls
- forms and ad hoc module reads should both reuse it

## 3. Auth Bootstrap And Page Guards

Files involved now:

- `apps/web/src/layouts/AuthBootstrap.astro`
- `apps/web/src/layouts/Layout.astro`
- `apps/web/src/layouts/AppShell.astro`

What must change:

- `GET /api/auth/session` must use bearer auth
- missing/expired token should publish unauthenticated state cleanly
- protected pages should still redirect to `/login`
- public-only pages should still redirect away after a valid login/session

Target behavior:

- login stores token
- bootstrap validates token on page load
- protected pages render only after auth readiness
- logout clears token and returns to `/login`

## 4. Login / Logout Browser Flow

Files involved now:

- `apps/web/src/pages/login.astro`
- `apps/web/src/components/TopBar.astro`
- `apps/web/src/layouts/ApiFormHandler.astro`

What must change:

- login success must capture the JSON token response before redirect
- login failure handling should still show a friendly invalid-credentials message
- logout should clear the browser token even though backend logout is stateless
- do not wait for cookie clearing because there is no auth cookie in `server-next`

## 5. Calories Forms And Read Models

This is the largest contract rename slice.

Files already showing legacy assumptions:

- `apps/web/src/pages/calories/new.astro`
- `apps/web/src/modules/calories/components/QuickAdd.astro`
- `apps/web/src/modules/calories/components/CaloriesDashboardClient.astro`
- `apps/web/src/modules/calories/components/CaloriesSettingsClient.astro`
- `apps/web/src/modules/calories/components/LogList.astro`
- `apps/web/src/modules/calories/components/DailySummary.astro`
- `apps/web/src/modules/dashboard/components/DashboardOverviewClient.astro`
- `apps/web/src/modules/dashboard/components/CaloriesSummaryCard.astro`

What must change:

- rename outgoing form field names to explicit backend names
- rename incoming UI reads to explicit backend response names
- rename quick-add hidden fields
- rename dashboard summary shape usage
- rename settings form field names and settings fetch shape usage
- rename log list renderers and recent meal/quick-add model assumptions

This slice should be done consistently in one pass to avoid mixed-contract bugs.

## 6. Expenses And Heat Under Bearer Auth

These domains look much closer structurally.

Files to revisit mainly for auth transport, not route shape:

- `apps/web/src/modules/expenses/components/ExpensesSettingsClient.astro`
- `apps/web/src/modules/expenses/components/ExpensesDashboardClient.astro`
- `apps/web/src/modules/heat/components/HeatDashboardClient.astro`
- any mutation forms already pointing to canonical expense and heat routes

What must change:

- replace cookie-dependent protected fetch behavior with bearer auth
- ensure modal/form helpers still work once auth is header-based
- verify no hidden legacy delete-by-query assumptions remain

## 7. Dashboard Overview Wiring

The home overview is already on the right backend direction because it does not depend on `/api/dashboard`.

Files involved now:

- `apps/web/src/modules/dashboard/components/DashboardOverviewClient.astro`

What must change:

- use bearer-auth fetches
- update calories summary shape from legacy names to explicit names
- keep partial-card loading behavior so one failing module request does not block the whole screen

## Recommended Migration Order

### Phase 1: Shared auth/browser plumbing

Do this first.

1. Introduce browser token storage helpers.
2. Update `ApiFormHandler.astro` to attach bearer auth to protected API calls.
3. Update `AuthBootstrap.astro` to validate bearer session.
4. Update login/logout flows to set and clear the token.

Deliverable:

- browser can authenticate against `apps/server-next`
- protected pages can bootstrap session state correctly

### Phase 2: Calories contract migration

Do this second because it has the highest payload/response mismatch.

1. update calorie create payload names
2. update quick-add payload names
3. update target settings payload names
4. update dashboard/settings/read-model response names
5. update any delete assumptions to `/api/calories/logs/{id}`

Deliverable:

- calories UI reads and writes work against `server-next`

### Phase 3: Expenses and heat bearer migration

1. switch protected fetches to bearer auth
2. verify mutation forms still work through shared runtime helpers
3. confirm delete/resource routes still align

Deliverable:

- expenses and heat work against `server-next` without cookie auth

### Phase 4: Regression/test migration

1. update frontend Vitest expectations that still assume old calories names or old auth transport
2. update Playwright flows to work with the new login/session behavior
3. keep `apps/server-next/scripts/e2e.sh` passing as backend smoke coverage
4. run browser tests against the migrated frontend/backend pair

Deliverable:

- regression suite matches the new boundary

## Verification Plan

## Backend-first checks

Keep these green throughout the frontend migration:

```bash
cd apps/server-next
go run ./cmd/seed-user
./scripts/e2e.sh
```

This proves backend integrity independently of frontend work.

## Frontend integration checks

After each frontend slice:

- manually log in through the browser
- refresh a protected page and confirm session bootstrap still works
- load `/`, `/calories`, `/expenses`, and `/heat`
- perform one mutation per migrated domain

## Regression commands

Use these once the frontend slice is updated:

```bash
cd apps/web
pnpm test
pnpm test:e2e
```

And keep backend checks alongside them:

```bash
cd apps/server-next
go test ./...
./scripts/e2e.sh
```

## Definition Of Done

The frontend/backend wiring is considered good when:

- login stores and uses bearer tokens successfully
- session bootstrap works after full page reloads and Astro navigations
- logout clears client auth state reliably
- all protected browser fetches and mutation forms use bearer auth
- calories request/response shapes fully match `server-next`
- expenses and heat work through canonical `server-next` routes
- the overview page loads module cards without `/api/dashboard`
- backend smoke checks pass
- frontend unit/e2e regression checks pass
- no frontend code still relies on cookie auth for `server-next`

## Non-Goals

This migration plan does not require:

- bringing back `/api/dashboard`
- keeping cookie auth compatibility in `server-next`
- preserving old calories field names forever
- hiding intentional contract breaks instead of documenting them

## File Checklist

High-priority frontend files to touch:

- `apps/web/src/layouts/AuthBootstrap.astro`
- `apps/web/src/layouts/ApiFormHandler.astro`
- `apps/web/src/pages/login.astro`
- `apps/web/src/components/TopBar.astro`
- `apps/web/src/modules/dashboard/components/DashboardOverviewClient.astro`
- `apps/web/src/modules/calories/components/CaloriesDashboardClient.astro`
- `apps/web/src/modules/calories/components/CaloriesSettingsClient.astro`
- `apps/web/src/modules/calories/components/QuickAdd.astro`
- `apps/web/src/pages/calories/new.astro`
- `apps/web/src/modules/expenses/components/ExpensesSettingsClient.astro`
- `apps/web/src/modules/expenses/components/ExpensesDashboardClient.astro`
- `apps/web/src/modules/heat/components/HeatDashboardClient.astro`

Supporting docs/tests likely to change:

- `docs/APPLICATION.md`
- `docs/TESTING.md`
- `apps/web/tests/forms.test.ts`
- `apps/web/tests/read-paths.test.ts`
- `apps/web/tests/e2e/*.spec.ts`

## Recommended Working Rule

Do not migrate everything blindly at once.

Use this loop:

1. keep backend smoke green
2. migrate one frontend slice
3. verify manually
4. update tests
5. move to the next slice

That keeps auth transport bugs, calories contract bugs, and general browser/runtime bugs isolated enough to debug.
