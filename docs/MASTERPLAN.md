# TrackStack & Platform Engineering: Master Plan

Date: 2026-03-10
Strategy: The Infrastructure Matrix Monorepo
Goal: Build TrackStack as a product and as a platform engineering portfolio piece, with Astro as the frontend shell and Go as the backend source of truth.

## 1. Core Philosophy

This repo is a monorepo where application boundaries make multiple deployment targets possible.

1. The app is split into an Astro frontend and a Go backend.
2. The same backend contracts should run across local containers, serverless production, and later lab environments.
3. Production stays serverless and cost-aware; container and Kubernetes environments remain optional labs with explicit cost controls.

## 2. Current Architecture Direction

TrackStack is in a hybrid migration phase.

- Astro owns pages, layouts, browser interaction, and thin API adapters.
- Go owns business API contracts, domain rules, and transport tests.
- Auth login, logout, and session verification are now Go-backed.
- Astro still owns request-local SSR auth context and some migration-era helpers.

Target end state:

- Astro becomes the UI shell only.
- Go becomes the source of truth for all application data access and mutations.
- Frontend request handling should go through Go endpoints rather than direct database access.

## 3. Repository Shape

```text
trackstack/
├── apps/
│   ├── server/
│   │   ├── cmd/
│   │   │   └── server/          # Current HTTP entrypoint for local/container runs
│   │   └── internal/
│   │       ├── core/            # config, db, logging
│   │       ├── modules/         # domain logic, ports, DTOs, db adapters
│   │       ├── transport/       # HTTP transport, middleware, OpenAPI
│   │       └── wiring/          # composition helpers per domain
│   └── web/                     # Astro frontend shell
├── infra/
│   ├── bootstrap/
│   ├── modules/
│   └── environments/
│       ├── serverless/
│       ├── ecs/
│       └── eks/
└── .github/workflows/
```

Planned later only when needed:

- `apps/server/cmd/lambda/` for the serverless Go entrypoint
- `apps/server/cmd/service/<domain>/` if domains are split into separate services

## 4. Infrastructure Matrix

| Level | Environment | Compute Type | Persistence | Cost Strategy | Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 0 | `local` | `docker-compose` | Turso | Free | Active |
| 1 | `serverless` | Astro + Go on Lambda/CloudFront/S3 | Turso (HTTP-oriented) | Production scale-to-zero | WIP |
| 2 | `ecs` | Containers on Fargate | Turso | Ephemeral lab | Planned |
| 3 | `eks` | Kubernetes | Turso | Ephemeral lab | Planned |
| 4 | `k3s` | Self-hosted VPS | Turso | Cheap fixed-cost lab | Backlog |

Database strategy:

- Turso remains the persistence layer for all domains.
- Runtime selection of HTTP vs WS-style connection settings is environment-driven in Go.
- Frontend runtime should not be the long-term owner of domain DB access.

## 5. Execution Roadmap

### Phase 1: Hybrid Astro + Go Migration

Goal: finish moving the app boundary so Go is the operational backend and Astro is a thin frontend shell.

Current progress:

- Go backend module structure is in place.
- Local HTTP server entrypoint exists in `apps/server/cmd/server/main.go`.
- Calories, expenses, heat, dashboard, and auth contracts are actively served by Go.
- Astro adapters proxy login/logout/session verification and migrated domain mutations to Go.
- Regression guardrails exist at three layers: Go transport tests, frontend contract tests, and Playwright browser flows.

Remaining work in this phase:

1. Remove remaining frontend-side direct DB usage from the application request path.
2. Add the production Go Lambda entrypoint and complete serverless deployment wiring.
3. Keep tightening endpoint contracts and regression coverage as refactors land.
4. Review security, performance, and architecture consistency of the Go backend.

### Phase 2: Serverless Production Completion

Goal: make the hybrid frontend/backend production deployment clean and durable.

Tasks:

1. Add `apps/server/cmd/lambda/main.go`.
2. Finalize Terraform and CI/CD wiring for Astro frontend plus Go backend deployment.
3. Ensure migrations, runtime config, and rollback posture are documented and repeatable.

### Phase 3: Container Lab (ECS)

Goal: prove the same backend contracts run in containerized compute without changing domain logic.

Tasks:

1. Package the Go backend cleanly for long-running container execution.
2. Build ECS environment Terraform and deployment workflow.
3. Add explicit destroy/cleanup automation to preserve FinOps guardrails.

### Phase 4: Orchestration Lab (EKS)

Goal: demonstrate Kubernetes operation with the same backend contracts.

Tasks:

1. Add manifests or Helm packaging.
2. Build EKS infrastructure and deploy workflow.
3. Reuse the same application contracts and keep the lab ephemeral.

## 6. Technical Guardrails

1. Strict hexagonal architecture: Go modules do not import transport or cloud runtime packages.
2. Thin transport adapters: request parsing, auth, error mapping, serialization.
3. Stateless compute: runtime state lives in Turso or object storage, never local disk.
4. Frontend/backend boundary tests are mandatory when contracts change.
5. FinOps first: non-serverless environments must have an explicit cost-control plan.

## 7. Success Criteria

TrackStack is considered structurally successful when:

- Go owns all business API contracts.
- Astro is primarily a frontend shell and adapter layer.
- Local and production environments use the same backend contract surface.
- Regression coverage catches contract breakage before refactors ship.
- The repo demonstrates both product engineering and platform engineering maturity.
