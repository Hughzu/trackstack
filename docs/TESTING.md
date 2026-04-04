# Testing Guide for Trackstack

This document defines the regression workflow for the Astro frontend plus Go backend split.

## Overview

Use three test layers:

1. `Go tests` to catch backend and transport regressions quickly.
2. `Vitest` to validate frontend contracts and form wiring.
3. `Playwright e2e` to verify real browser flows across the frontend -> Go boundary.

## Test Commands

### Backend (Go)

From repo root:

Or run the pieces directly:

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

`pnpm test:e2e` runs the Playwright browser flows against the currently running frontend/backend pair. If you want to refresh the backend test user first, run `pnpm seed:e2e-user` from `apps/web` on the host or `go run ./cmd/seed-user` from `apps/server`.

### Backend-First Smoke Checks

Run the rebuild backend directly:

```bash
cd apps/server
go run ./cmd/server
```

In another shell:

```bash
cd apps/server
go run ./cmd/seed-user
./scripts/e2e.sh
```

Notes:

- `./scripts/e2e.sh` defaults to `BASE_URL=http://localhost:8080`.
- It loads `apps/server/.env` automatically and requires `E2E_TEST_EMAIL` and `E2E_TEST_PASSWORD` there or in the shell.
- `go run ./cmd/seed-user` loads `apps/web/.env` first and `apps/server/.env` second.
- Override `BASE_URL` when validating a non-default local port or remote environment.

### Split Validation: Heat Service

Run the standalone heat service directly:

```bash
cd apps/server
PORT=18080 go run ./cmd/heat-api
```

In another shell, validate the split runtime is alive:

```bash
curl http://localhost:18080/health
curl -i http://localhost:18080/api/heat/dashboard
```

Notes:

- `GET /health` should return `200` with `{"status":"ok"}`.
- `GET /api/heat/dashboard` without a bearer token should return `401`.
- `cmd/heat-api` loads only heat-service config: `TURSO_HEAT_*`, `JWT_SECRET`, shared DB pool envs, and standard runtime envs.

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

- `apps/web/tests/forms.test.ts`
  - Required `data-api-form` attributes across mutation forms
  - Redirect contract checks
  - Guard against inline fetch logic bypassing the shared form runtime
  - Guard the client auth bootstrap wiring for login/public-only vs protected pages

- `apps/web/tests/read-paths.test.ts`
  - Guard that the home dashboard no longer imports SSR auth context
  - Guard that the dashboard overview reads expenses, calories, and heat module endpoints independently in the browser after auth bootstrap
  - Guard that the calories dashboard no longer imports SSR auth context
  - Guard that the calories dashboard reads `/api/calories/dashboard` in the browser after auth bootstrap
  - Guard that the expenses and heat dashboards no longer import SSR auth context
  - Guard that the expenses and heat dashboards read their Go endpoints in the browser after auth bootstrap
  - Guard that the calories settings page no longer imports SSR auth context
  - Guard that the calories settings page reads `/api/calories/target` in the browser after auth bootstrap
  - Guard that the expenses settings page no longer imports SSR auth context
  - Guard that the expenses settings page reads `/api/expenses/settings` in the browser after auth bootstrap
  - Guard that Astro middleware is removed and `AuthBootstrap` is the active page guard
  - Guard that Astro auth adapter route files are removed
  - Guard that legacy SSR auth helper files are removed
  - Guard that legacy SSR service wrappers are removed

- `apps/web/tests/api-config.test.ts`
  - Guard API path normalization for browser calls
  - Guard that direct browser calls can target Go-owned auth and domain endpoints consistently

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
  - Assert canonical `DELETE /api/heat/refills/{id}` succeeds and the list shrinks
  - Assert heat mutations target canonical Go-owned `/api/heat/refills` resources

- `apps/server/scripts/e2e.sh`
  - Backend-first curl smoke coverage for the rebuild runtime
  - Logs in with bearer auth and validates `GET /api/auth/session`
  - Exercises heat create/list/delete
  - Exercises calories target update, log create, dashboard read, and delete
  - Exercises expenses settings read/update, entry create/delete, checklist create/delete, recurring create/delete, and dashboard read
  - Validates `POST /api/auth/logout` stays stateless and returns `204`

- `apps/server/**`
  - No dedicated Go regression test files yet in the rebuilt backend workspace
  - The current backend safety net is `go test ./...` plus the seed-user command and `./scripts/e2e.sh`
  - Add focused Go tests next for handler parsing, auth behavior, and service invariants as contracts settle

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
  - verify Playwright browsers are installed in the environment where the command is running, and verify `apps/web/.env` contains valid `E2E_TEST_EMAIL` and `E2E_TEST_PASSWORD` values.
- If login e2e fails unexpectedly:
  - rerun `pnpm seed:e2e-user` on the host or `go run ./cmd/seed-user` from `apps/server`, then rerun `pnpm test:e2e`.
- If Go tools are missing in container runs:
  - execute the command through the `go-backend` service.
- If `apps/server/scripts/e2e.sh` fails early:
  - verify `apps/server` is running on the expected `BASE_URL`, rerun `go run ./cmd/seed-user`, and confirm `JWT_SECRET` plus `TURSO_*` env vars are set for the rebuild runtime.
## Auth Boundary

- Browser auth now talks directly to Go under `/api/auth/*`.
- `ApiFormHandler.astro` stores the login bearer token in browser storage and clears it on successful logout.
- Protected-page gating happens through `AuthBootstrap.astro` calling `GET /api/auth/session` with `X-Trackstack-Authorization` and redirecting client-side.
- Playwright login helpers now explicitly assert that `/api/auth/login` succeeds before continuing into module tests.
- When auth transport behavior changes, update both Go auth transport tests and at least one browser flow that covers login plus post-login page bootstrap.

For `apps/server`, bearer-auth changes should still be validated backend-first with the rebuild runtime and `curl`:

- run `go run ./cmd/server` from `apps/server`
- run `go run ./cmd/seed-user` from `apps/server`
- run `./scripts/e2e.sh` from `apps/server`
- verify `POST /api/auth/login` returns a JSON bearer token payload and no auth cookie
- verify `GET /api/auth/session` returns `401` without `Authorization: Bearer <jwt>` and `200` with it
- verify protected routes under `/api/heat/*`, `/api/calories/*`, and `/api/expenses/*` reject missing bearer headers and succeed with a valid bearer token
- verify `POST /api/auth/logout` returns `204` and does not invalidate an already-issued token because the current rebuild stays stateless
- verify browser login stores client auth state, a protected page reload reboots auth through `GET /api/auth/session`, and logout clears the client token before redirect
