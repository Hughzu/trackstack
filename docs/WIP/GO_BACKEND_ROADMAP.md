# Go Backend Roadmap (WIP)

This roadmap is designed to help you learn Go while building the Trackstack backend without locking into a rigid, waterfall process. It focuses on small, repeatable building blocks that scale from monolith to microservices.

## 1) Roadmap (phased, minimal lock-in)

1. **Bootstrap the runtime**
   - Build `cmd/server` with Chi and a `/health` endpoint.
   - Add `internal/core/config` and `internal/core/logging`.
   - Goal: one running HTTP binary with JSON responses.

2. **First domain module (end-to-end)**
   - Pick the smallest domain (e.g., `heat` or `calories`).
   - Implement ports → service → db adapter → HTTP handler.
   - Goal: one real feature wired end-to-end.

3. **Make a repeatable module template**
   - Convert the first module into a checklist (see below).
   - Goal: every new module follows the same steps.

4. **Introduce auth/session**
   - Implement auth middleware and `auth.FromContext`.
   - Goal: handlers stay clean; user identity is always from context.

5. **Add Lambda transport**
   - Wrap the same router in `cmd/lambda`.
   - Goal: transport swap without touching domain code.

6. **Prove microservice extraction**
   - Create `cmd/service/<domain>` for one domain.
   - Swap in-process adapter for HTTP/gRPC adapter.
   - Goal: same domain logic runs in a dedicated service.

## 2) Module implementation checklist

Use this checklist for every new domain module.

1. **Define ports (domain boundary)**
   - `internal/modules/<domain>/ports.go`
   - Keep interfaces minimal and domain-focused.

2. **Define DTOs (API + persistence)**
   - `internal/modules/<domain>/types.go`
   - Use RFC3339 for API time values; keep nullable fields explicit.

3. **Implement service (pure logic)**
   - `internal/modules/<domain>/service.go`
   - Validate input, normalize dates, map domain errors.

4. **DB adapter (Turso)**
   - `internal/modules/<domain>/adapters/db/`
   - Use `database/sql` with `QueryContext` / `ExecContext` and explicit scan mapping.

5. **Transport handlers (JSON‑only)**
   - `internal/transport/http/<domain>.go`
   - Parse JSON → call service → map errors → JSON response.
   - **No redirects, no HTML** (UI handles navigation).

6. **Routes**
   - `internal/transport/http/router.go`
   - Register under `/api/<domain>` with plural noun paths.
   - Add transition aliases only if needed (remove later).

7. **Wiring (composition root)**
   - `internal/wiring/<domain>/` builds DSN → DB → adapter → service.
   - `cmd/server/main.go` only wires handlers/router and lifecycle.

8. **OpenAPI**
   - Update `internal/transport/http/openapi.yaml` to match routes + response shapes.

9. **Compatibility check (Astro parity)**
   - Ensure status codes + JSON shape align with existing Astro API behavior.
   - Keep API behavior clean; update Astro UI instead of adding redirects.

10. **Test the boundary**
   - Handler test or service unit test for a core use-case.

## 3) First module walkthrough (minimal example)

### Step 1: Ports

```go
// internal/modules/heat/ports.go
package heat

type RefillStore interface {
  ListByRange(userID string, from string, to string) ([]Refill, error)
  Create(userID string, in CreateRefillInput) (Refill, error)
}
```

### Step 2: Types

```go
// internal/modules/heat/types.go
package heat

type Refill struct {
  ID       string
  Date     string
  WeightKg float64
  Bags     int
  Season   *string
}

type ListRefillsRequest struct {
  UserID string
  From   string
  To     string
}

type CreateRefillInput struct {
  Date     string
  WeightKg float64
  Bags     int
}
```

### Step 3: Service

```go
// internal/modules/heat/service.go
package heat

type Service struct {
  store RefillStore
}

func NewService(store RefillStore) *Service {
  return &Service{store: store}
}

func (s *Service) ListRefills(req ListRefillsRequest) ([]Refill, error) {
  return s.store.ListByRange(req.UserID, req.From, req.To)
}
```

### Step 4: Transport handler

```go
// internal/transport/http/heat.go
package httptransport

type HeatHandler struct {
  svc *heat.Service
}

func (h *HeatHandler) ListRefills(w http.ResponseWriter, r *http.Request) {
  userID := auth.FromContext(r.Context())
  from := r.URL.Query().Get("from")
  to := r.URL.Query().Get("to")

  resp, err := h.svc.ListRefills(heat.ListRefillsRequest{
    UserID: userID,
    From:   from,
    To:     to,
  })

  writeJSON(w, resp, err)
}
```

### Step 5: Routes

```go
// internal/transport/http/router.go
r.Route("/api/heat", func(r chi.Router) {
  r.Get("/refills", handlers.Heat.ListRefills)
})
```

### Step 6: Wiring

```go
// cmd/server/main.go
heatStore := heatdb.NewRefillStore(db)
heatSvc := heat.NewService(heatStore)
handlers := httptransport.NewHandlers(httptransport.Deps{
  HeatService: heatSvc,
})
```

## 4) Learning goals per phase

- Phase 1: core Go HTTP, middleware, error handling
- Phase 2: domain modeling + service patterns
- Phase 3: adapters + config + runtime wiring
- Phase 4: transport swapping (server → lambda)
- Phase 5: service extraction patterns
