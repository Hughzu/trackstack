# TrackStack & Platform Engineering: Master Plan

**Date:** 2026-04-12
**Strategy:** The Infrastructure Matrix Monorepo
**Goal:** Build TrackStack as a product and platform engineering showcase. Solid.js for the UI, Go for the absolute source of truth.
About me : Platform engineer with an SRE title and a software engineering background, focused on reliability, automation, and end-to-end system design.

## 1. Core Philosophy

1. **Write Once, Run Anywhere:** The exact same backend code must run flawlessly across local containers, a zero-cost serverless production environment, and microservice Kubernetes clusters.
2. **Zero Domain Coupling:** Domains are strictly isolated. No bleeding boundaries.
3. **Scale-to-Zero:** We don't bankrupt ourselves for a side project. Production is cheap; labs are ephemeral.

## 2. The Infrastructure Matrix



| Level | Environment | Compute Type | Persistence | Cost Strategy |
| :--- | :--- | :--- | :--- | :--- |
| 0 | `local` | `docker-compose` | Turso | Free |
| 1 | `serverless` | Lambda, CloudFront, S3 | Turso | Scale-to-zero ($0 target) |
| 2 | `k3s` | Self-hosted VPS | Turso | Cheap fixed-cost lab |
| 3 | `eks` / `ecs` | Cloud Kubernetes / Fargate | Turso | Ephemeral lab only |

*Turso (SQLite at the edge) remains the persistence layer across all environments.*

## 3. Technical Guardrails

1. **Infrastructure Agnostic Architecture:** The software architecture must support running identically as a monolith locally, a serverless binary in Lambda, or distributed microservices in K8s. If we have to rewrite business logic to change the deploy target, we fucked up.
2. **FinOps First:** Production must remain serverless and aggressively cost-aware. Labs must have explicit destroy automation so I don't set my wallet on fire.
3. **No Compromise on Performance:** We stay in the $0 tier, but the app still needs to scream. Performance dictates engineering decisions within those cost boundaries.

## 4. Success Criteria

- Domain logic is completely portable across deployment targets without modification.
- The Solid frontend is violently fast and entirely detached from backend implementation quirks.
- We can tear down and rebuild any infrastructure tier from scratch predictably.
- The AWS bill stays firmly at $0.
