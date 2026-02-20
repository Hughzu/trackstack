# Frontend Regression: Delete Modals Failing in Production

This document captures the production-only regression we hit after refactoring client runtime logic, plus the test coverage we can add to prevent repeats.

## Issue Summary

- Symptoms
  - In production only: delete actions stopped working across calories, expenses, and heat.
  - Clicking the trash icon did nothing (confirm dialog never opened).
  - Local dev and preview appeared fine.

- Impact
  - Users could not delete entries from list views (major UX regression).
  - The problem was not caught before deploy.

## Root Cause

We moved UI runtime logic (confirm modal and dropdown bindings) from inline scripts into an external module loaded via `ClientRuntime`. In the deployed build (Lambda SSR + CloudFront static assets), the runtime module did not reliably execute or re-bind after transitions. This made the delete triggers inert in production while still working locally.

Key contributing factors:

- Client runtime was bundled as a separate asset and loaded after swap. In prod, asset loading behavior diverged from local dev.
- Confirm modal handlers relied on document-level event binding that was not re-established correctly in deployed navigation flow.

Resolution:

- Moved the runtime back inline (`ClientRuntime.astro`) so it is embedded in the HTML payload and executes consistently in production.
- Kept `astro:page-load` and `astro:after-swap` re-binding to support client router transitions.

## Tests to Add (to prevent this regression)

### 1. API Route Contract Tests (fast)

Goal: Validate the server behavior for POST/DELETE routes independent of UI.

Suggested coverage:

- `POST /api/expenses/expense`
- `DELETE /api/expenses/expense`
- `POST /api/calories/log`
- `DELETE /api/calories/log`
- `POST /api/heat/refill` and delete endpoint

Assertions:

- 2xx on valid payload
- 4xx on invalid payload
- Expected JSON response shape

Why this helps:

- Ensures the backend behavior is correct even if UI binding changes.

### 2. Preview Build UI Tests (Playwright)

Goal: Catch production-only regressions caused by build output differences.

Run tests against `astro build` + `astro preview`, not `astro dev`.

Scenarios:

- Expenses list: click trash icon, confirm modal opens, confirm deletion, list updates.
- Calories: click trash icon, confirm modal opens, delete, list updates.
- Heat: click trash icon, confirm modal opens, delete, list updates.

Assertions:

- Confirm dialog appears when clicking the trash icon.
- Network request is sent (DELETE to expected endpoint).
- UI updates (row removed or page reloads).

Why this helps:

- Mirrors the production asset split and deployment stack, where the regression occurred.

### 3. Client Runtime Smoke Test

Goal: Verify that runtime bindings are attached after navigation.

Minimal checks:

- `window.signedFetch` exists.
- After `astro:after-swap`, a click on a trash icon opens a dialog.

Why this helps:

- Specifically validates the part that broke in production.

## How These Tests Prevent the Regression

- The Playwright tests would fail in the exact production build if the runtime code does not bind, catching the issue before deploy.
- Contract tests ensure the API stays stable when refactors happen on the UI side.
- The runtime smoke test validates that client router transitions still activate the modal bindings.

## Suggested CI Placement

- Run contract tests on every PR for fast feedback.
- Run Playwright tests on `main` or for changes under `apps/web/**`.
- Execute Playwright against `astro preview` (dist/client + dist/server), not dev server.

## Follow-up Implementation Notes

- Keep runtime bindings in `ClientRuntime.astro` inline unless we add explicit guarantees for external client assets in the Lambda + CloudFront stack.
- If we re-extract runtime into a module, add a dedicated test that asserts the runtime script loads and binds in the preview build.
