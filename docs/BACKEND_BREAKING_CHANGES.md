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

This document should be updated whenever additional intentional compatibility breaks are introduced.
