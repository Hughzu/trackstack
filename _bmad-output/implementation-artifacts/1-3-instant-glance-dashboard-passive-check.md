# Story 1.3: "Instant Glance" Dashboard (Passive Check)

Status: ready-for-dev

## Story

As a Lazy Human,
I want to see my remaining calories and protein in large, color-coded numbers on the home screen,
So that I can make a meal decision in < 1 second.

## Acceptance Criteria

1. **Given** I have a defined calorie/protein target and current logs in the database
2. **When** I open the app
3. **Then** the dashboard should display "Remaining" values (Target - Consumed) for the current day
4. **And** the numbers should be Green (>20% left), Yellow (5-20%), or Red (<5% or over)
5. **And** the UI should follow the "Progress-Enhanced Cards" design:
   - Bold hero numbers (80px for calories, 60px for protein)
   - Muted labels ("kcal remaining")
   - Horizontal progress bars showing consumption percentage

## Tasks / Subtasks

- [ ] Initialize `calories` module schema (AC: 1)
  - [ ] Create `meal_logs` table: `id`, `user_id`, `name`, `calories`, `protein`, `carbs`, `fat`, `created_at`
  - [ ] Create `user_targets` table: `user_id`, `calorie_target`, `protein_target` (default: 2000 kcal, 150g protein)
- [ ] Implement Dashboard Backend Logic (AC: 3)
  - [ ] Create `internal/calories` repository and service
  - [ ] Implement `GetDailySummary(ctx, userID)` to calculate totals for today
- [ ] Implement UI Components with Templ + Tailwind (AC: 5)
  - [ ] Set up Tailwind CSS configuration with status colors
  - [ ] Create `internal/calories/components/dashboard.templ`
  - [ ] Implement `MetricWithProgress` component with color-coding logic (AC: 4)
- [ ] Update Routing (AC: 2)
  - [ ] Update root route `/` in `internal/common/server` to render the Dashboard

## Dev Notes

- **Relevant architecture patterns:** Modular Monolith, Component-based UI (Templ), Server-Side Rendering.
- **Source tree components to touch:** `internal/calories/`, `internal/common/server/`, `tailwind.config.js`.
- **Timezone handling:** Use the server's local time or UTC for "today" calculation for now (Story 3.1 will handle personalization).
- **Default Targets:** If no record exists in `user_targets`, use 2000 kcal / 150g protein as defaults.

### References

- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#The Defining Experience: "The Instant Glance"]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Design System Foundation]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.3]

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
