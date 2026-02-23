# Go Backend Architecture (WIP)

This document defines the **concrete contracts** and architectural boundaries for the Trackstack Go backend. It is the source of truth for the rewrite and should guide implementation decisions.

Roadmap and implementation walkthrough: `docs/WIP/GO_BACKEND_ROADMAP.md`.

## 1) Goals and readiness

- The repo is ready to begin the backend rewrite **as a staged migration**.
- The Astro app remains the UI shell (SSG), and all domain logic moves to Go.
- The Go backend must run unchanged in Lambda (prod) and ECS/EKS (labs).

## 2) Core principles

1. **Hexagonal architecture**
   - `internal/modules/**` contains pure business logic, ports, and adapters.
   - No imports of HTTP, Lambda, or framework types in modules.
2. **Stateless compute**
   - All state lives in Turso (or S3 for static assets). No local disk reliance.
3. **Transport isolation**
   - `cmd/lambda` and `cmd/server` wire the same modules to different transports.
4. **Environment-driven DB mode**
   - `DB_CONNECTION_MODE=HTTP|WS` selects the Turso protocol.

## 3) Repository layout (new backend)

```
apps/trackstack/
├── cmd/
│   ├── lambda/          # AWS Lambda entrypoint (serverless prod)
│   ├── server/          # HTTP server entrypoint (docker monolith)
│   └── service/         # Domain-scoped microservice entrypoints
│       ├── users/
│       ├── calories/
│       ├── expenses/
│       └── heat/
├── internal/
│   ├── core/            # config, logger, db, auth primitives
│   ├── modules/         # domain modules (ports + services + adapters)
│   └── transport/       # shared transport wiring
│       ├── http/        # router, middleware, error mapping
│       └── lambda/      # lambda request adapters
└── pkg/                 # optional shared types (use sparingly)
```

## 4) Domain modules and contracts

Each module has three layers:

1) **Ports (interfaces)**
2) **Service (use-case logic)**
3) **Adapters (DB, time, crypto, etc.)**

### Module directory template

```
internal/modules/<domain>/
├── ports.go           # interface definitions
├── types.go           # request/response DTOs
├── service.go         # use-case implementation
└── adapters/
    └── db/            # Turso implementation for the port
```

### Dependency rules (hard boundaries)

1. **Modules do not import each other.**
2. **Modules do not import transport or AWS packages.**
3. **Cross-domain access uses ports** defined by the consuming module.
4. **Transport code depends on modules, never the reverse.**
5. **Adapters depend on ports, never the reverse.**

### Cross-domain adapter placement (consistency rule)

- **Outbound adapters live with the consumer.**
  - Example: `expenses` needs `users` data, so the `users` client adapter lives in `internal/modules/expenses/adapters/`.
- **Inbound adapters live with the provider.**
  - Example: the `users` gRPC/HTTP server adapter lives under `cmd/service/users/`.
- This rule keeps the module boundaries stable across phases:
  - In-process adapter (monolith) and gRPC/HTTP adapter (microservice) are in the same place.
  - Only the adapter implementation changes when you split services.

Example (consumer vs provider placement):

```go
// internal/modules/expenses/ports.go (consumer owns the port)
package expenses

type UserPlanReader interface {
  GetPlanTier(userID string) (string, error)
}
```

```go
// internal/modules/expenses/adapters/users_grpc.go (outbound adapter lives with consumer)
package expenses

type UsersGrpcClient interface {
  GetPlanTier(userID string) (string, error)
}

type GrpcUserPlanReader struct {
  client UsersGrpcClient
}

func (r *GrpcUserPlanReader) GetPlanTier(userID string) (string, error) {
  return r.client.GetPlanTier(userID)
}
```

```go
// cmd/service/users/grpc_server.go (inbound adapter lives with provider)
package main

type UsersServer struct {
  svc *users.Service
}

func (s *UsersServer) GetPlanTier(ctx context.Context, req *pb.GetPlanTierRequest) (*pb.GetPlanTierResponse, error) {
  tier, err := s.svc.GetPlanTier(req.UserId)
  if err != nil {
    return nil, err
  }
  return &pb.GetPlanTierResponse{Tier: tier}, nil
}
```

### Common DTO conventions

- All request structs are named `XRequest` and response structs `XResponse`.
- All time values are RFC3339 UTC strings in API responses.
- Domain errors are typed and mapped at transport boundaries.

## 5) Core services (internal/core)

### Config contract

```
type Config struct {
  Env                string // local|serverless|ecs|eks
  DBConnectionMode   string // HTTP|WS
  OriginVerifyHeader string
  OriginVerifyValue  string
  AuthCookieName     string
  AuthCookieSecure   bool
  AuthCookieSameSite string // lax|strict|none
  AuthSessionIdleSeconds      int
  AuthSessionAbsoluteSeconds  int
  AuthSessionRotateAfterSeconds int
  AuthSessionRotationGraceSeconds int
  AuthSessionTouchSeconds int
}
```

### Logger contract

```
type Logger interface {
  Debug(msg string, fields ...any)
  Info(msg string, fields ...any)
  Warn(msg string, fields ...any)
  Error(msg string, fields ...any)
}
```

### DB client contract

```
type DB interface {
  Query(ctx context.Context, sql string, args ...any) (Rows, error)
  Exec(ctx context.Context, sql string, args ...any) (Result, error)
  Close() error
}
```

The DB constructor must accept `DBConnectionMode` and select:

- HTTP: Turso HTTP endpoint (serverless)
- WS: Turso WebSocket endpoint (ecs/eks)

## 6) Transport layer contracts

Transport layers are thin adapters that:

1. Parse request input
2. Call module service
3. Map errors to HTTP
4. Serialize response

Minimal handler shape (auth + error mapping):

```go
// internal/transport/http/handlers.go
func (h *ExpensesHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
  userID := auth.FromContext(r.Context())
  resp, err := h.svc.GetSettings(r.Context(), expenses.GetSettingsRequest{UserID: userID})
  writeJSON(w, resp, err)
}
```

### HTTP transport (cmd/server)

- Base path: `/api` for JSON endpoints
- JSON only, no HTML

Chi routing (transport only):

```go
// internal/transport/http/router.go
package httptransport

func NewRouter(handlers Handlers) http.Handler {
  r := chi.NewRouter()
  r.Use(requestID, logging, recoverer)

  r.Route("/api/expenses", func(r chi.Router) {
    r.Get("/settings", handlers.Expenses.GetSettings)
    r.Post("/settings", handlers.Expenses.SetSettings)
  })

  return r
}
```

```go
// cmd/server/main.go
func main() {
  handlers := buildHandlers()
  router := httptransport.NewRouter(handlers)

  srv := &http.Server{
    Addr:    ":8080",
    Handler: router,
  }

  log.Fatal(srv.ListenAndServe())
}
```

### Lambda transport (cmd/lambda)

- API Gateway or Function URL
- Same JSON contract as HTTP server

```go
// cmd/lambda/main.go
func main() {
  handlers := buildHandlers()
  router := httptransport.NewRouter(handlers)
  lambda.Start(lambdaadapter.New(router))
}
```

### Domain microservice transport (cmd/service/<domain>)

- Exposes only one domain's API surface
- Uses the same module contracts as monolith
- Cross-domain dependencies swap from in-process adapters to HTTP client adapters

### Error mapping contract

```
type ErrorCode string

const (
  ErrUnauthorized ErrorCode = "unauthorized"
  ErrForbidden    ErrorCode = "forbidden"
  ErrNotFound     ErrorCode = "not_found"
  ErrInvalidInput ErrorCode = "invalid_input"
  ErrConflict     ErrorCode = "conflict"
  ErrInternal     ErrorCode = "internal"
)

type APIError struct {
  Code    ErrorCode `json:"code"`
  Message string    `json:"message"`
}
```

HTTP status mapping:

- `unauthorized` -> 401
- `forbidden` -> 403
- `not_found` -> 404
- `invalid_input` -> 422
- `conflict` -> 409
- `internal` -> 500

## 7) Auth/session contracts

The Go backend replaces the Astro middleware + session system.

### Auth ports

```
type SessionStore interface {
  Create(ctx context.Context, userID string, uaHash string, ipPrefix string) (Session, error)
  Validate(ctx context.Context, sessionID string) (Session, error)
  Touch(ctx context.Context, sessionID string) (Session, error)
  Rotate(ctx context.Context, sessionID string) (Session, error)
  Revoke(ctx context.Context, sessionID string) error
}

type PasswordHasher interface {
  Hash(password string) (string, error)
  Verify(password, hash string) (bool, error)
}
```

### Session DTO

```
type Session struct {
  ID                string
  UserID            string
  CreatedAt         time.Time
  ExpiresAt         time.Time
  RotatedAt         time.Time
  LastSeenAt        time.Time
  AbsoluteExpiresAt time.Time
  ParentID          *string
  RevokedAt         *time.Time
}
```

## 8) Domain API contracts (initial set)

These are the first JSON endpoints the SSG frontend will consume.

### Auth

- `POST /api/auth/login`
  - Request: `{ "email": "...", "password": "..." }`
  - Response: `{ "ok": true }`
- `POST /api/auth/logout`
  - Response: `{ "ok": true }`
- `GET /api/auth/me`
  - Response: `{ "id": "...", "email": "..." }`

### Calories

- `GET /api/calories/logs?from=YYYY-MM-DD&to=YYYY-MM-DD`
- `POST /api/calories/logs`
- `GET /api/calories/targets`
- `POST /api/calories/targets`

### Expenses

- `GET /api/expenses/settings`
- `POST /api/expenses/settings`
- `GET /api/expenses/sheet/current`
- `POST /api/expenses/entries`

### Heat

- `GET /api/heat/refills?from=YYYY-MM-DD&to=YYYY-MM-DD`
- `POST /api/heat/refills`

Each endpoint maps to a module service method with a corresponding DTO.

## 9) Module example: Calories contracts

```
type CaloriesService interface {
  ListLogs(ctx context.Context, req ListLogsRequest) (ListLogsResponse, error)
  CreateLog(ctx context.Context, req CreateLogRequest) (CreateLogResponse, error)
  GetTargets(ctx context.Context) (GetTargetsResponse, error)
  SetTargets(ctx context.Context, req SetTargetsRequest) (SetTargetsResponse, error)
}

type ListLogsRequest struct {
  UserID string
  From   string // YYYY-MM-DD
  To     string // YYYY-MM-DD
}

type CalorieLog struct {
  ID       string
  DateTime string // RFC3339
  Calories int
  ProteinG int
  CarbsG   *int
  FatG     *int
  Title    *string
}
```

## 10) Migration rules

1. Astro becomes **SSG only** for pages and layout.
2. All data fetching uses Go API endpoints.
3. All mutations go through Go API endpoints (no SigV4 browser proxy once Go is live).
4. Update infra to route `/api/*` to Go Lambda and `/assets/*` to S3.
5. Update docs (`ARCHITECTURE.md`, `APPLICATION.md`, `TESTING.md`) as boundaries move.

## 11) Monolith to microservices path

1. **Monolith phase:** `cmd/server` or `cmd/lambda` wires all domains in-process.
2. **Hybrid phase:** extract one domain into `cmd/service/<domain>`.
3. **Microservices phase:** one domain per service; cross-domain calls use HTTP adapters.

## 12) Non-goals (for Phase 1)

- WebSockets for real-time updates
- Full E2E tests
- GraphQL
