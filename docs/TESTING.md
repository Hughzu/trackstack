# Testing Guide for Trackstack

This document defines the regression workflow for the Astro frontend plus Go backend split.

## Overview

Use three test layers:

1. `Go tests` to catch backend and transport regressions quickly.
2. `Vitest` to validate frontend contracts and form wiring.
3. `Playwright e2e` to verify real browser flows across the Astro -> Go boundary.

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

`pnpm test:e2e` seeds the configured users database before running browser flows. It expects `apps/web/.env` to provide `TURSO_USERS_URL`, `E2E_TEST_EMAIL`, and `E2E_TEST_PASSWORD`, plus `TURSO_USERS_TOKEN` for remote Turso.

## Compose Regression Workflow

From repo root:

```bash
docker compose up --build -d
docker compose exec -T go-backend sh -lc 'go test ./...'
docker compose exec -T astro-frontend sh -lc 'pnpm test'
docker compose exec -T astro-frontend sh -lc 'pnpm test:e2e'
```

This is the default local loop after changing Astro API routes, Go handlers, or request/response contracts between them.

## Current Regression Coverage

- `apps/server/internal/transport/http/calories_test.go`
  - Authenticated `POST /api/calories/log` JSON request
  - Authenticated `POST /api/calories/log` form request
  - Assert created vs redirect behavior at the HTTP transport boundary

- `apps/server/internal/modules/heat/service_test.go`
  - Basic service validation and season behavior regression checks

- `apps/web/tests/forms.test.ts`
  - Required `data-api-form` attributes across mutation forms
  - Redirect contract checks
  - Guard against inline fetch logic bypassing the shared form runtime

- `apps/web/tests/e2e/calories.spec.ts`
  - Login through Astro
  - Submit calorie form
  - Assert `POST /api/calories/log` succeeds
  - Assert redirect to `/calories`

- `apps/web/tests/e2e/expenses.spec.ts`
  - Login through Astro
  - Submit expense form
  - Assert `POST /api/expenses/expense` succeeds
  - Assert redirect to `/expenses`

## Adding Regression Tests

When a regression is found in the frontend/backend boundary:

1. Add or update a Go transport test if the issue is in request parsing, auth, or HTTP response behavior.
2. Add or update a Playwright test if the issue is only reproducible through the browser flow.
3. Add or update `apps/web/tests/forms.test.ts` when a form contract changes.
4. Keep assertions at the boundary: status code, redirect target, and minimal visible success behavior.

## Troubleshooting

- If frontend container fails after changing dependencies:
  - recreate the `web_node_modules` volume.
- If `pnpm test:e2e` fails before the browser starts:
  - verify `apps/web/.env` contains valid `TURSO_USERS_URL`, `E2E_TEST_EMAIL`, and `E2E_TEST_PASSWORD` values.
- If login e2e fails unexpectedly:
  - rerun `pnpm test:e2e`; it reseeds the test credentials each run.
- If Go tools are missing in container runs:
  - execute the command through the `go-backend` service.
