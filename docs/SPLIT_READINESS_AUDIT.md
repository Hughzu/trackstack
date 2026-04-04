# Split Readiness Audit

Date: 2026-04-04

## Verdict

The rebuilt Go backend in `apps/server` is a solid modular monolith with a mostly clean hexagonal shape inside each context.

It is not fully ready to split into microservices yet.

The main issue is not the domain packages. The main issue is the shared runtime assembly: one bootstrap, one router, one config surface, and one process that still knows about every database and every module.

## What Was Audited

- Docs and intended backend shape in `docs/ARCHITECTURE.md`, `docs/APPLICATION.md`, and `docs/TESTING.md`
- Runtime assembly in `apps/server/internal/app/monolithapi`
- Current entrypoints in `apps/server/cmd`
- Auth/session middleware and handler boundaries
- Local validation with:
  - `go test ./...`
  - `go run ./cmd/seed-user`
  - `./scripts/e2e.sh`

## Current Strengths

### 1. The context shape is mostly good

The following contexts already follow a sensible hexagonal structure:

- `auth`
- `calories`
- `expenses`
- `heat`
- `users`

Each context mostly separates:

- domain
- application ports
- application services
- inbound adapters where needed
- outbound adapters

That is exactly what you want before a split.

### 2. Cross-context coupling is limited

`calories`, `expenses`, and `heat` look reasonably isolated and are the strongest candidates for extraction into separate services.

### 3. The current monolith is runtime-portable

The backend already supports:

- local HTTP via `cmd/server`
- Lambda via `cmd/lambda`

That shared-runtime discipline is useful. It proves the business code is not welded to one deployment target.

### 4. Local validation works

The current backend passed:

- `go test ./...` in `apps/server`
- `go run ./cmd/seed-user`
- backend smoke e2e via `apps/server/scripts/e2e.sh`

Important caveat: `go test ./...` currently acts as a compile check more than a true regression suite because there are no dedicated Go test files yet.

## Current Weaknesses

### 1. The composition root is still monolithic

`apps/server/internal/app/monolithapi/runtime.go` builds the whole system in one shot.

Today the runtime:

- loads one global config
- opens all domain databases
- wires all modules together
- mounts all routes in one router

This is the biggest blocker to a clean microservice split.

### 2. Config is global instead of service-local

`apps/server/internal/platform/config/config.go` currently requires credentials for all databases plus `JWT_SECRET`.

That means every process must know about everything.

For real services, each binary should load only what it actually needs.

### 3. Database wiring is all-or-nothing

`apps/server/internal/app/monolithapi/database.go` opens:

- calories DB
- expenses DB
- heat DB
- users DB

That is fine for the modular monolith. It is wrong for extracted services.

### 4. Auth and users are not separate service seams today

`apps/server/internal/app/monolithapi/auth_module.go` builds auth by directly constructing the users service.

That means `auth` and `users` are not really two service boundaries yet.

## Auth/User Decision

Decision for the next architecture step:

- keep `auth` and `users` in one service
- keep two internal hexagons if that helps maintain clarity
- do not treat the user service as an API gateway

This is the right call.

The current codebase already behaves more like one identity boundary than two independent services.

## Recommended Auth Model For The Split

### Short version

Each future microservice should validate JWTs locally.

Do not call the auth/user service on every request just to verify a token.

### Why

If every service has to round-trip to auth before doing anything, then:

- auth becomes a latency bottleneck
- auth becomes a single point of failure
- the whole system becomes tightly coupled at runtime

That is how people accidentally build a distributed monolith.

### Recommended model

- one identity service owns login and token issuance
- each domain service validates bearer tokens locally using the shared signing secret or public key
- handlers extract claims and pass explicit `userID` values into application commands
- if a service needs richer user data than token claims provide, it calls the identity service explicitly for that use case only

## Split Readiness By Domain

### Ready sooner

- `calories`
- `expenses`
- `heat`

These look closest to real service boundaries.

### Needs a boundary decision first

- `auth`
- `users`

These should now be treated as one future service boundary: identity.

## What Is Missing Before The Split

### 1. Service-specific composition roots

You need separate bootstraps instead of one global `NewRuntime()`.

Expected direction:

- identity bootstrap
- calories bootstrap
- expenses bootstrap
- heat bootstrap

### 2. Service-specific entrypoints under `cmd/`

The current `cmd/` layout is correct for the monolith:

- `cmd/server`
- `cmd/lambda`
- `cmd/seed-user`

For a split simulation, the backend should gain dedicated service entrypoints such as:

- `cmd/identity-api`
- `cmd/calories-api`
- `cmd/expenses-api`
- `cmd/heat-api`

These should stay thin and delegate to service-local bootstrap packages.

### 3. Service-local config validation

Each service should validate only its own environment variables.

Examples:

- identity service: users DB + JWT secret
- calories service: calories DB + JWT verification config
- expenses service: expenses DB + JWT verification config
- heat service: heat DB + JWT verification config

### 4. Service-local router assembly

Each service should mount only its own routes plus health/openapi endpoints as needed.

### 5. Real Go tests

Before splitting, add backend tests for:

- auth middleware behavior
- auth handler contracts
- handler parsing and status mapping for each domain
- a few focused application service invariants per context

Right now the backend has smoke coverage and compilation coverage, but not enough contract-level regression protection for a split.

## Next Steps

### Phase 1: Refactor for service-local assembly inside the monolith repo

1. Create service-level bootstrap packages for:
   - identity
   - calories
   - expenses
   - heat
2. Keep `auth` and `users` as one service boundary named `identity`
3. Move config loading toward service-local config structs
4. Move DB wiring toward service-local open/close logic
5. Keep the existing all-in-one runtime temporarily if it still helps local development

### Phase 2: Add dedicated service entrypoints

Add thin binaries under `apps/server/cmd/` for each future service:

- `cmd/identity-api`
- `cmd/calories-api`
- `cmd/expenses-api`
- `cmd/heat-api`

Each binary should:

- load only its own config
- build only its own dependencies
- expose only its own routes

### Phase 3: Prove local split readiness

Create a separate compose file for split simulation, for example:

- `docker-compose.microservices.yml`

It should run:

- identity service
- calories service
- expenses service
- heat service
- optional local gateway or proxy only if needed

This should validate that the current architecture can be deployed as separate processes without changing domain logic.

### Phase 4: Lock down tests before moving further

Add Go tests before treating the split as complete.

Recommended first targets:

- `ResolveSession` middleware
- `auth` login/session HTTP handlers
- one happy path and one failure path per domain handler
- one or two application service tests per context

## Final Take

The backend rewrite is working.

The hexagonal structure inside the contexts is good enough to support a future service split.

What still needs work is the runtime boundary around those contexts.

So the answer is:

- yes, the architecture is heading in the right direction
- no, it is not yet truly split-ready
- the next move is to extract service-local bootstrap, config, and `cmd/` entrypoints
- the identity boundary should be one service containing the `auth` and `users` hexagons
