# Backend Breaking Changes

This file tracks intentional contract breaks between frontend and backend during the backend rewrite.

## Calories (`apps/server`)

Date: 2026-03-22

### HTTP contract breaks

- `DELETE /api/calories/log?id=<id>` becomes `DELETE /api/calories/logs/{id}`.
- Calories response fields move from gram suffix abbreviations to explicit names:
  - `proteinG` -> `proteinGrams`
  - `carbsG` -> `carbGrams`
  - `fatG` -> `fatGrams`
- Target response fields move from legacy names to explicit names:
  - `targetKcal` -> `targetCalories`
  - `targetProteinG` -> `targetProteinGrams`
  - `targetCarbsG` -> `targetCarbGrams`
  - `targetFatG` -> `targetFatGrams`
- Dashboard summary fields move to explicit names:
  - `consumed` -> `consumedCalories`
  - `protein` -> `proteinGrams`
  - `carbs` -> `carbGrams`
  - `fat` -> `fatGrams`

### Request payload breaks

- Log creation payload moves from legacy request fields to explicit names:
  - `protein` -> `proteinGrams`
  - `carbs` -> `carbGrams`
  - `fat` -> `fatGrams`
- Target update payload moves from legacy request fields to explicit names:
  - `targetKcal` -> `targetCalories`
  - `targetProtein` -> `targetProteinGrams`
  - `targetCarbs` -> `targetCarbGrams`
  - `targetFat` -> `targetFatGrams`

## Expenses (`apps/server`)

Date: 2026-03-22

### HTTP contract breaks

- `DELETE /api/expenses/entries?id=<id>` becomes `DELETE /api/expenses/entries/{id}`.
- `DELETE /api/expenses/checklists?id=<id>` becomes `DELETE /api/expenses/checklists/{id}`.
- `DELETE /api/expenses/recurring?id=<id>` becomes `DELETE /api/expenses/recurring/{id}`.

### Compatibility notes

- Legacy query-parameter delete routes are intentionally removed and are not kept as compatibility aliases.
- Existing expenses payload and response shapes otherwise stay aligned with the current frontend contract during the `apps/server` migration slice.

## Heat (`apps/server`)

Date: 2026-03-23

### HTTP contract breaks

- `DELETE /api/heat/refills?id=<id>` becomes `DELETE /api/heat/refills/{id}`.
- The legacy query-parameter delete route is intentionally removed and is not kept as a compatibility alias.

## General API Surface (`apps/server`)

Date: 2026-03-23

### HTTP contract breaks

- `GET /api/dashboard` is intentionally removed from `apps/server`.
- The overview/read-model direction is now explicit module-level endpoints instead of one aggregate backend dashboard endpoint.
- Health payload shape changes from `{"ok":true}` to `{"status":"ok"}` for both `/health` and `/api/health`.

## Authentication (`apps/server` target)

Date: 2026-03-23

### Auth Contract Breaks

- The backend migration will adopt a **Stateless JWT** approach to achieve $0 infrastructure cost and 0 Turso DB reads on API requests.
- `apps/server` now uses `Authorization: Bearer <jwt>` as the only auth transport for protected endpoints.
- `POST /api/auth/login` now returns `200 OK` JSON with `accessToken`, `tokenType`, `expiresAt`, and `userId` instead of issuing an auth cookie.
- `GET /api/auth/session` now requires an `Authorization: Bearer <jwt>` header.
- Protected `apps/server` routes under `/api/calories/*`, `/api/expenses/*`, and `/api/heat/*` no longer accept cookie auth.
- `POST /api/auth/logout` is now a stateless `204 No Content` endpoint; it does not revoke tokens server-side and does not clear a browser cookie.
- The current Astro frontend is intentionally incompatible with `apps/server` until it is migrated from cookie auth to bearer auth.

### Config Contract Breaks

- `apps/server` now requires `JWT_SECRET` for bearer token signing and verification.
- `AUTH_COOKIE_NAME` and `AUTH_COOKIE_SECURE` are intentionally removed from the rebuild config surface.
- Legacy `SESSION_COOKIE_NAME` and `SESSION_COOKIE_SECURE` remain unsupported.
- Existing `serverless-next` runtime and deploy wiring must be updated to provide `JWT_SECRET` before `apps/server` can fully replace the old backend in that environment.

This document should be updated whenever additional intentional compatibility breaks are introduced.
