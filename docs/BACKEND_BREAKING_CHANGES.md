# Backend Breaking Changes

This file tracks intentional contract breaks between frontend and backend during the backend rewrite.

## Calories (`apps/server-next`)

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

## Expenses (`apps/server-next`)

Date: 2026-03-22

### HTTP contract breaks

- `DELETE /api/expenses/entries?id=<id>` becomes `DELETE /api/expenses/entries/{id}`.
- `DELETE /api/expenses/checklists?id=<id>` becomes `DELETE /api/expenses/checklists/{id}`.
- `DELETE /api/expenses/recurring?id=<id>` becomes `DELETE /api/expenses/recurring/{id}`.

### Compatibility notes

- Legacy query-parameter delete routes are intentionally removed and are not kept as compatibility aliases.
- Existing expenses payload and response shapes otherwise stay aligned with the current frontend contract during the `apps/server-next` migration slice.

## Authentication (`apps/server-next` target)

Date: 2026-03-22

### Upcoming Auth Contract Breaks

- The backend migration will adopt a **Stateless JWT** approach to achieve $0 infrastructure cost and 0 Turso DB reads on API requests.
- The frontend must be updated to handle this change. It should expect an `HttpOnly` JWT cookie instead of the existing Turso session string. 
- If the Astro frontend currently validates sessions directly against Turso during Server-Side Rendering (SSR), this logic must be rethought to either cryptographically verify the JWT using the shared secret natively, or delegate authentication checks entirely to the Go backend proxy.

This document should be updated whenever additional intentional compatibility breaks are introduced.
