# Testing Guide for Trackstack

This document defines the regression workflow for the Astro frontend plus Go backend split.

## Overview

Use three test layers:

1. `Go tests` to catch backend and transport regressions quickly.
2. `Vitest` to validate frontend contracts and form wiring.
3. `Playwright e2e` to verify real browser flows across the frontend -> Go boundary.

## Test Commands

### Backend (Go)

```bash
cd apps/server
go test ./...
```

### Frontend Unit Tests (Vitest)

```bash
cd apps/web
pnpm test
pnpm test:watch
```

### Frontend Integration Tests (Playwright)

```bash
cd apps/web
pnpm test:e2e
```

`pnpm test:e2e` seeds the configured users database through the backend-owned seed command before running browser flows. The test workflow remains compatible with `apps/web/.env` for `E2E_TEST_EMAIL` and `E2E_TEST_PASSWORD`; the backend command also supports `apps/server/.env` or exported environment variables for DB config.

## Compose Regression Workflow

From repo root:

```bash
docker compose up --build -d
docker compose exec -T go-backend sh -lc 'go test ./...'
docker compose exec -T astro-frontend sh -lc 'pnpm test'
docker compose exec -T astro-frontend sh -lc 'pnpm test:e2e'
```

This is the default local loop after changing Astro auth routes, Go handlers, or request/response contracts between them.

## Current Regression Coverage

- `apps/server/internal/transport/http/auth_test.go`
  - Login JSON request
  - Logout JSON request
  - Session verification request
  - Assert cookie issuance and auth transport behavior

- `apps/server/internal/transport/http/calories_test.go`
  - Authenticated `POST /api/calories/log` JSON request
  - Authenticated non-JSON `POST /api/calories/log` rejection
  - Authenticated `POST /api/calories/target` JSON request
  - Authenticated `DELETE /api/calories/log?id=...` request
  - Assert JSON-only transport behavior at the HTTP boundary

- `apps/server/internal/transport/http/expenses_test.go`
  - Authenticated `POST /api/expenses/settings` JSON request
  - Authenticated `POST /api/expenses/checklists` JSON request
  - Authenticated `POST /api/expenses/recurring` JSON request
  - Authenticated non-JSON `POST /api/expenses/entries` rejection
  - Assert removed legacy expense aliases return `404`

- `apps/server/internal/modules/heat/service_test.go`
  - Basic service validation and season behavior regression checks

- `apps/web/tests/forms.test.ts`
  - Required `data-api-form` attributes across mutation forms
  - Redirect contract checks
  - Guard against inline fetch logic bypassing the shared form runtime
  - Guard the client auth bootstrap wiring for login/public-only vs protected pages

- `apps/web/tests/read-paths.test.ts`
  - Guard that the home dashboard no longer imports SSR auth context
  - Guard that the dashboard overview reads `/api/dashboard` in the browser after auth bootstrap

- `apps/web/tests/e2e/calories.spec.ts`
  - Login through Astro
  - Submit calorie form
  - Update calorie targets
  - Delete a calorie log
  - Assert `POST /api/calories/log` succeeds
  - Assert target updates and delete flow work against canonical Go-owned `/api/calories/*` endpoints

- `apps/web/tests/e2e/expenses.spec.ts`
  - Login through Astro
  - Submit expense form
  - Save expense settings
  - Add monthly checklist template
  - Add recurring template
  - Assert canonical Go-owned expense endpoints (`/api/expenses/entries`, `/api/expenses/checklists`) succeed
  - Assert settings and template mutations work against Go-owned `/api/expenses/*`

- `apps/web/tests/e2e/heat.spec.ts`
  - Login through Astro
  - Create a refill
  - Delete a refill from the heat history list
  - Assert `DELETE /api/heat/refills` succeeds and the list shrinks
  - Assert heat mutations target canonical Go-owned `/api/heat/refills`

## Adding Regression Tests

When a regression is found in the frontend/backend boundary:

1. Add or update a Go transport test if the issue is in request parsing, auth, or HTTP response behavior.
2. Add or update a Playwright test if the issue is only reproducible through the browser flow.
3. Add or update `apps/web/tests/forms.test.ts` when a form contract changes.
4. Keep assertions at the boundary: status code, redirect target, and minimal visible success behavior.

## Troubleshooting

- If frontend container fails after changing dependencies:
  - recreate the `web_node_modules` volume.
- If the homepage or any server-rendered page throws `fetch failed` in Docker:
  - rebuild `go-backend` after any `apps/server/go.mod` or Go runtime image change so the running container is not stuck on an older toolchain image.
- If `pnpm test:e2e` fails before the browser starts:
  - verify `apps/web/.env` contains valid `E2E_TEST_EMAIL` and `E2E_TEST_PASSWORD` values, and verify `apps/server/.env` or your shell environment provides valid Turso users DB credentials if they are not already shared.
- If login e2e fails unexpectedly:
  - rerun `pnpm test:e2e`; it reseeds the test credentials each run.
- If Go tools are missing in container runs:
  - execute the command through the `go-backend` service.
## Auth Boundary

- `apps/web/src/pages/api/auth/login.ts` and `apps/web/src/pages/api/auth/logout.ts` are thin Astro adapters that proxy directly to the Go auth endpoints and forward `Set-Cookie` / redirect headers back to the browser.
- `apps/web/src/middleware.ts` now verifies page-request sessions through `GET /api/auth/session`, so session validation and rotation semantics come from Go as well.
- Cookie names and session timing defaults are aligned between `apps/web` and `apps/server`; if you change one, update both configs together.
- Playwright login helpers now explicitly assert that `/api/auth/login` succeeds before continuing into module tests.
- When auth transport behavior changes, update both Go auth transport tests and at least one browser flow that authenticates through `/api/auth/login`.
