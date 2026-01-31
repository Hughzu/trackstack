# Story 1.4: Manual Numeric Quick-Log

Status: ready-for-dev

<!-- Note: Validation is optional. Run validate-create-story for quality check before dev-story. -->

## Story

As a User,
I want to manually enter a numeric calorie/protein value when I'm in a rush,
so that I can track my intake even if I don't want to describe the meal.

## Acceptance Criteria

1. **Given** I am on the dashboard
2. **When** I enter numbers into the "Quick Add" fields and submit
3. **Then** the system should create a new meal log entry with those values
4. **And** the dashboard numbers should update instantly via HTMX without a page reload
5. **And** the meal name should be auto-generated as "Quick Entry" (user can edit via Story 3.2)

## Tasks / Subtasks

- [ ] Add Quick-Add Form to Dashboard (AC: 1, 4)
  - [ ] Create compact inline form below metric cards
  - [ ] Input fields: Calories (number), Protein (number, optional)
  - [ ] Submit button with visual indicator
- [ ] Implement HTMX Endpoint for Quick-Log (AC: 2, 3)
  - [ ] Create `POST /api/meals/quick-log` handler in `internal/calories/handlers`
  - [ ] Parse form values, create `meal_logs` entry with name="Quick Entry"
  - [ ] Return updated dashboard HTML fragment for HTMX swap
- [ ] Update Dashboard Service Layer (AC: 3)
  - [ ] Add `QuickLogMeal(ctx, userID, calories, protein)` method to service
  - [ ] Auto-set carbs=0, fat=0 for quick entries (to be refined via AI in Story 2.x)
- [ ] HTMX Integration (AC: 4)
  - [ ] Add `hx-post` to form with target="#dashboard-container"
  - [ ] Use `hx-swap="innerHTML"` for seamless update
  - [ ] Add loading state indicator on submit

## Dev Notes

- **Relevant architecture patterns:** Modular Monolith, HTMX for interactivity, Repository Pattern.
- **Source tree components to touch:** `internal/calories/handlers/`, `internal/calories/service.go`, `internal/calories/components/dashboard.templ`.
- **UI/UX alignment:** Follow "Progress-Enhanced Cards" design from UX spec. Quick-add should be compact and unobtrusive (below metric cards, above Recents).
- **Previous Story Context:** Story 1.3 established the dashboard, database schema (`meal_logs`, `user_targets`), and service layer. Build on these patterns.

### Project Structure Notes

- Handlers: `internal/calories/handlers/calories_handler.go` (extend existing)
- Service: `internal/calories/service.go` (extend existing)
- Components: `internal/calories/components/dashboard.templ` (modify to add quick-log form)
- All meal logs go to the same `meal_logs` table used in Story 1.3

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.4]
- [Source: _bmad-output/planning-artifacts/ux-design-specification.md#Design System Foundation]
- [Source: _bmad-output/implementation-artifacts/1-3-instant-glance-dashboard-passive-check.md]

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

