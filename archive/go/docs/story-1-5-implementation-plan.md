# Implementation Plan: 1-5 Muscle Memory Recents Carousel

## Design Concept: "Quick Reply" Chips Bar
Instead of a bulky gallery-style carousel, we will implement a high-efficiency **Quick Reply Bar** positioned directly above the fixed log form.

### Key Features
- **Placement:** Fixed at the bottom, sitting between the dashboard content and the input form.
- **UI Pattern:** Horizontal scrolling container of pill-shaped chips.
- **Ranking Logic:** Frequent + Recent. Shows the top 10 unique meals ordered by `COUNT(*)` DESC and `MAX(logged_at)` DESC.
- **One-Tap Logging:** Tapping a chip immediately posts the meal's stored macros to `/calories/log`.

## Technical Tasks

### 1. Backend (Repository & Service)
- [ ] **Repository:** Add `GetFrequentMeals(ctx, userID, limit) ([]Meal, error)` to `Repository` interface.
  - SQL: `SELECT name, calories, protein, carbs, fat FROM meals WHERE user_id = ? GROUP BY name ORDER BY COUNT(*) DESC, MAX(logged_at) DESC LIMIT ?`
- [ ] **Models:** Add `FrequentMeals []Meal` to the `DailySummary` struct in `models.go`.
- [ ] **Service:** Update `CalculateDailySummary` to fetch and include frequent meals.

### 2. UI Components (Templ)
- [ ] **Component:** Create `internal/calories/components/frequent_meals_bar.templ`.
  - Styling: `flex overflow-x-auto no-scrollbar gap-2 px-4 py-2`
  - Item: `button` with `whitespace-nowrap rounded-full bg-white/5 border border-white/10 px-3 py-1 text-xs`
- [ ] **Integration:** Update `internal/calories/components/quick_log_form.templ` to include the bar.
- [ ] **Interactivity:** Add HTMX attributes to chips:
  ```html
  hx-post="/calories/log"
  hx-vals='{"name": "Oatmeal", "calories": 350, ...}'
  hx-target="#dashboard-metrics"
  hx-swap="outerHTML"
  ```

### 3. UX Polish
- [ ] **Auto-Scroll:** Ensure the "Brunch Flow" works (Dashboard updates and user can immediately see the result).
- [ ] **Empty State:** Hide the bar if no historical data exists.

## Success Criteria
- [ ] User can log a repeat meal with exactly **one tap**.
- [ ] The dashboard updates instantly via HTMX OOB swap or standard target swap.
- [ ] The interface remains "thumb-friendly" and responsive on mobile.
