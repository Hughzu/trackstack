# Backend Target Architecture

Date: 2026-03-17

This document describes the target shape for the Go backend rewrite. It is not the description of the current implementation. It is the direction we want to move toward, slice by slice.

## Why This Exists

The current backend works, but the structure is hard to reason about:

- `transport`, `modules`, and `wiring` do not make adapter direction obvious.
- HTTP handlers currently depend on concrete services too directly.
- Some module ports are too large.
- Some reads hide write-side behavior.
- The folder shape does not make future Lambda and service extraction concerns easy to see.

The goal of this target architecture is to make boundaries explicit so the codebase is easier to navigate, safer to evolve with an LLM, and ready for both Lambda deployment and future service extraction.

## Core Direction

We keep a hexagonal style, but make the vocabulary clearer.

Data flow should read like this:

`inbound adapters -> application -> domain -> outbound ports -> outbound adapters`

At the module level, we want each bounded context to be easy to understand in isolation.

## Target Repository Shape

```text
apps/server/
├── cmd/
│   ├── lambda/                 # Lambda entrypoint
│   └── server/                 # HTTP server entrypoint
└── internal/
    ├── platform/               # shared technical building blocks
    │   ├── config/
    │   ├── db/
    │   └── logging/
    ├── app/
    │   └── bootstrap/          # composition root and runtime assembly
    └── contexts/
        ├── heat/
        │   ├── adapters/
        │   │   ├── inbound/
        │   │   │   └── http/
        │   │   └── outbound/
        │   │       └── db/
        │   ├── application/
        │   │   ├── ports/
        │   │   └── services/
        │   └── domain/
        ├── auth/
        ├── users/
        ├── calories/
        └── expenses/
```

## Vocabulary

### `platform/`

Shared technical concerns used by multiple contexts.

- configuration
- database opening and pooling
- logging
- other reusable runtime helpers that are not domain behavior

`platform` must not become a hidden business logic bucket.

### `app/bootstrap/`

The composition root.

This is where the runtime is assembled:

- load config
- open databases
- construct adapters
- wire application services
- mount HTTP router
- expose the final handler for `cmd/server` and `cmd/lambda`

This is the only place allowed to know the full concrete graph.

### `contexts/<name>/domain/`

Pure business concepts and invariants.

Examples:

- entities
- value objects
- domain rules
- domain errors
- logic that should survive any transport or persistence change

The domain must not know about HTTP, Lambda, SQL, or framework concerns.

### `contexts/<name>/application/ports/`

Interfaces owned by the application layer.

These ports describe what the use cases need from the outside world.

Examples:

- repositories
- query interfaces
- unit-of-work or transaction boundaries if needed later
- external service interfaces

Ports should be narrow and shaped around use cases, not giant module-wide god interfaces.

### `contexts/<name>/application/services/`

Use cases and orchestration.

Examples:

- `CreateRefill`
- `DeleteRefill`
- `GetHeatDashboard`
- `Login`
- `GetSession`

Application services:

- validate command/query input
- call domain logic
- coordinate ports
- return application results

They should not contain HTTP parsing or SQL details.

### `contexts/<name>/adapters/inbound/`

Entry points into the application.

For now, that mainly means HTTP handlers. Later it could also include:

- Lambda-specific adapters if needed
- async consumers
- CLI commands for bounded contexts

Inbound adapters should stay thin:

- parse request
- map request to use-case input
- call exactly one application service 
- map result to response

### `contexts/<name>/adapters/outbound/`

Implementations of ports.

For now, this is mostly database code. Later it may include:

- external HTTP clients
- message publishers
- cache implementations

Outbound adapters are technical implementations, not owners of business decisions.

## Architectural Rules

### 1. One-way dependency direction

- inbound adapters depend on application
- application depends on domain and application-owned ports
- outbound adapters implement application ports
- domain depends on nothing outside itself and the standard library

### 2. Transport is not the application layer

HTTP handlers must not orchestrate multi-step business workflows. They should translate HTTP to a use case.

### 3. Domain is not a dumping ground

Not every non-HTTP function belongs in `domain/`. Cross-repository coordination, dashboard composition, and authentication workflows usually belong in `application/services/`.

### 4. Queries should not silently mutate state

If defaults or bootstrap data must be created, that behavior should be explicit in the application layer rather than hidden inside read flows.

### 5. Ports should be small

Avoid large interfaces like a single store owning unrelated settings, templates, sheets, entries, and read models.

### 6. Compatibility lives at the boundary

During migration, legacy route aliases or compatibility behavior should live in inbound adapters, not leak into domain logic.

## API Direction

We want cleaner, resource-oriented API routes.

Examples:

- `GET /api/heat/refills`
- `POST /api/heat/refills`
- `DELETE /api/heat/refills/{id}`
- `GET /api/heat/dashboard`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/session`

Rules:

- prefer path identifiers over delete-by-query-param contracts
- keep handlers thin
- one endpoint should call one application use case
- compatibility aliases are temporary and must be removable later

## Deployment Shape

This structure must support multiple runtimes without changing domain or application code.

### Current required runtimes

- local HTTP server
- AWS Lambda

### Future possible runtime

- extracted service per context

That means:

- `cmd/server` and `cmd/lambda` remain thin runtime entrypoints
- runtime-specific concerns stay outside contexts
- contexts should not import server or Lambda packages

## Migration Strategy

This is not a big-bang rewrite.

We migrate one vertical slice at a time.

If we use `apps/server-next/` as a rebuild workspace, it must stay intentionally narrow. It is a learning and design lab, not a second full backend. The rebuild should start with one context, then one use case at a time inside that context. Each use case should begin with explicit contract notes before code: input, output, invariants, side effects, and dependency boundaries. Tests should first prove data flow and business behavior with fakes or in-memory adapters, and only then add database wiring. No feature expansion should happen during this rebuild; the purpose is to gain clarity and ownership over contracts and flow.

Suggested order:

1. Define the target structure and vocabulary.
2. Rebuild `heat` as the first clean context.
3. Keep compatibility at the HTTP boundary so the frontend still works.
4. Use the `heat` slice as the template for the next contexts.
5. Rework auth after the new structure is proven.
6. Leave `expenses` for later because it has the most coupled behavior.

## Current Rebuild Progress (`apps/server-next/`)

The rebuild workspace now has the first runtime foundation in place:

- `cmd/server/main.go` is intentionally thin and only starts the HTTP runtime.
- `internal/app/bootstrap/` owns runtime assembly and router creation.
- `internal/platform/config/` and `internal/platform/logging/` hold shared technical concerns.
- a base `chi` router exists with simple health endpoints (`/health` and `/api/health`).

This is good progress because it proves the basic runtime shape without pulling business logic into `main`.

What is still intentionally missing:

- only the first `heat` vertical slice exists so far: `GET /api/heat/refills`
- the first slice currently uses a fake outbound adapter and a mocked user identity
- no real database wiring exists yet
- no Lambda entrypoint exists yet

The rebuild has now moved past runtime-only scaffolding. The next step is to deepen the `heat` slice use case by use case without adding unrelated runtime infrastructure.

Recommended next move after the first slice:

1. Keep `GET /api/heat/refills` as the reference slice for naming and dependency direction.
2. Add the next heat contract first before coding, likely `POST /api/heat/refills`.
3. Keep auth compatibility at the HTTP boundary; do not let application services read HTTP context.
4. Replace the fake outbound adapter only after the use-case contract and tests are stable.
5. Add database wiring after the service behavior is proven with tests.

The preferred first use case should be the smallest one that still proves the architecture. For this rebuild, `GET /api/heat/refills` is a better starting point than the dashboard because it is simpler, read-focused, and easier to test end-to-end without hiding extra orchestration.

## First Slice: Heat

`heat` is the preferred first rewrite target because it is small enough to reason about but still exercises the full flow:

- HTTP route handling
- authenticated requests
- read model generation
- write-side mutation
- database adapter implementation

The target `heat` context should clearly separate:

- domain rules like refill and season logic
- application use cases like dashboard retrieval and refill creation
- outbound db adapters
- inbound HTTP handlers

## Non-Goals

This target shape does not mean:

- introducing microservices now
- forcing every tiny helper into a separate package
- chasing purity at the expense of delivery
- rewriting the full backend before shipping value

The goal is clarity, boundary control, and incremental ownership.

## Definition Of Better

The backend is in a better place when:

- a developer can open one context and understand the full flow
- route handlers only translate HTTP concerns
- use cases are easy to test without HTTP or SQL
- ports are narrow and obvious
- Lambda and local HTTP share the same assembled application
- future service extraction is possible without rewriting domain logic
