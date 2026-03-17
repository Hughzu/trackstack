# Decisions (ADR Timeline)

Accepted architectural decisions for Trackstack. Date-stamped to prevent re-litigating the same choices.

Update guidelines:
- Add a new line for each durable architectural decision or reversal.
- Use the decision date (or the merge date if uncertain).
- Never delete old entries; if reversed, add a new entry that supersedes it.
- Keep each line to one sentence: decision + short rationale or constraint.

2026-02-23 — Use Turso (LibSQL) as the only persistent store to get distributed SQLite with low ops overhead and serverless-friendly access.
2026-02-23 — Split data into four Turso databases (users, calories, expenses, heat) to isolate domains and migrations per feature area.
2026-02-23 — Treat AWS serverless as production baseline (Lambda + CloudFront + S3) to enforce scale-to-zero and cost control.
2026-02-23 — Make CloudFront the only public ingress; route assets to S3 and all other paths to a Lambda Function URL secured with IAM and an origin verification header.
2026-02-23 — Store runtime secrets and connection strings in SSM Parameter Store and resolve at runtime; never inline secrets in code or Terraform.
2026-02-23 — Require GitHub Actions to use OIDC for AWS access; prohibit long-lived AWS keys in CI/CD.
2026-02-23 — Keep cost guardrails enabled (Budgets, price class 100, short log retention, artifacts lifecycle) to maintain $0/mo serverless posture.
2026-02-23 — Use Astro SSR for dynamic routes and API endpoints; ship static assets to S3 with long-lived caching.
2026-02-23 — Enforce Astro boundaries: pages handle routing and request context, modules handle domain logic, components/layouts are presentational only.
2026-02-23 — Use a global vanilla JS runtime (ClientRuntime) for UI interactivity; avoid React/Vue and other heavy client frameworks.
2026-02-23 — Standardize serverless mutations via SigV4 form interception using `data-api-form` and ApiFormHandler to avoid unsigned POSTs.
2026-02-23 — Centralize DB access through `sqlite.ts` and prohibit direct LibSQL client instantiation anywhere else.
2026-02-23 — Follow strict hexagonal architecture in the Go backend: business logic isolated in modules, transport concerns in cmd.
2026-02-23 — Make Turso connection mode configurable via `DB_CONNECTION_MODE` (HTTP for serverless, WebSockets for containers).
2026-02-23 — Gate deploys with migrations-first CI/CD; no rollback in CI until Atlas supports down migrations in community CLI.
2026-02-23 — Keep testing lightweight in the Astro phase: fast unit tests for form contracts, defer heavy backend testing until Go migration.
2026-03-10 — Use Astro as the frontend/auth adapter and Go as the business API runtime during the migration so browser contracts stay stable while domains move to Go incrementally.
2026-03-10 — Protect the Astro -> Go boundary with layered regression checks: Go transport tests, frontend form contract tests, and Playwright browser flows.
2026-03-10 — Make Go the source of truth for auth login, logout, and session verification while Astro keeps request-local auth context for SSR during the migration.
2026-03-10 — Keep Go as an in-process modular monolith first; only introduce `cmd/lambda` or domain-scoped services when deployment shape requires them.
2026-03-10 — Treat consumer-owned ports and outbound adapters as the cross-domain rule so future service extraction does not change module boundaries.
2026-03-11 — Treat the Go backend as a frontend-agnostic JSON API and keep browser form/redirect adaptation in the web client so runtime or frontend changes do not alter backend contracts.
2026-03-11 — Build Astro as a static frontend/PWA served from S3 and route all auth, health, and business API traffic directly to the Go runtime.
2026-03-13 — Make Turso transport selection explicit by requiring `*_URL_HTTP` or `*_URL_WS` for the chosen `DB_CONNECTION_MODE`, so benchmarks and serverless behavior do not silently fall back to `libsql://` URLs.
