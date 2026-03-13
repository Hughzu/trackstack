# Serverless Migration Plan: Old Astro SSR to Static Astro + Go API

Date: 2026-03-13
Status: Planning
Scope: Temporary parallel environment, validation, cutover, and cleanup

## Context

TrackStack is no longer the same application shape as the currently deployed serverless stack.

Local and documented architecture now assume:

- Astro frontend built as a static site.
- Go backend as the source of truth for auth and business APIs.
- CloudFront as the single public entrypoint.
- S3 serving static assets and HTML shell.
- Lambda serving Go API routes and health checks.

The currently deployed production stack is still the older model:

- Astro SSR-style Lambda is the default CloudFront origin.
- S3 is only used for asset paths.
- Deployment automation still packages and deploys Astro server output.
- The deployed frontend expects old route contracts that no longer match the current Go API wiring.

This means an in-place redeploy of the current serverless stack is unsafe. Even if infrastructure updates succeed, the old frontend and new backend contracts are not compatible enough for a simple switch.

## Verified Current State

### Application

- `apps/web/astro.config.mjs` uses `output: "static"`.
- `apps/server/cmd/server/main.go` is the current Go HTTP entrypoint for local and container runs.
- There is currently no dedicated Go Lambda entrypoint under `apps/server/cmd/lambda/`.

### Deployment automation

- `.github/workflows/deploy-serverless.yml` still builds Astro and packages `apps/web/dist/server` as the Lambda artifact.
- The deploy workflow does not yet build or deploy the Go backend as the production Lambda artifact.
- The deploy workflow change detection currently focuses on `apps/web/**` and migrations, not the Go backend deployment contract.

### Terraform

- `infra/modules/static-hosting/main.tf` sends default CloudFront traffic to Lambda.
- `infra/modules/static-hosting/main.tf` only routes known asset paths to S3.
- `infra/modules/lambda-api/main.tf` is parameterized like an Astro SSR Lambda runtime.
- `infra/environments/serverless/variables.tf` still defaults to Node runtime values such as `nodejs20.x` and `entry.handler`.

## Problem Statement

There is a deployment-contract gap between:

1. The application that now exists locally.
2. The infrastructure and workflow that still deploy the old application model.

The main risk is not code cleanliness. The main risk is deploying a mixed architecture where:

- old frontend code reaches removed routes,
- new frontend behavior reaches the wrong origin,
- CloudFront caches stale paths,
- S3 sync deletes old assets without a clean rollback path,
- Terraform mutates live resources whose purpose has fundamentally changed.

## Decision

Use a temporary parallel serverless environment to validate the new architecture before touching the old production stack.

This is the preferred approach because:

- the frontend/backend contract changed materially,
- the routing model changed materially,
- the deployment artifact changed materially,
- rollback is simpler if the old production stack stays intact until validation is complete.

This temporary environment should be structurally equivalent to the intended final architecture, not a one-off workaround.

## Target Architecture

The target serverless production shape is:

- CloudFront as the only public entrypoint.
- S3 as the default origin for the static Astro frontend.
- Go Lambda as the origin for `/api/*` and `/health`.
- Turso as the persistence layer.
- Runtime secrets and config loaded from SSM.

Routing target:

- `/api/*` -> Go Lambda
- `/health` -> Go Lambda
- static assets and frontend routes -> S3 via CloudFront

## Proposed Temporary Environment

Create a second serverless environment using a distinct identity, for example:

- environment name: `serverless-next`
- resource prefix: `trackstack-next`
- SSM prefix: `/trackstack/serverless-next`
- Lambda function name: `trackstack-go-api-next`

Resources should be separate from the current environment:

- separate CloudFront distribution,
- separate S3 assets bucket,
- separate artifacts bucket,
- separate Lambda function,
- separate SSM output and runtime paths.

The existing production stack should remain untouched until the temporary stack passes validation.

## Implementation Plan

### Phase 1: Define the new serverless contract

Goal: make the architecture explicit before changing infra and CI/CD.

Tasks:

1. Confirm the final public routing contract.
2. Confirm the Go runtime environment variable contract for Lambda.
3. Confirm which frontend routes must resolve from S3 and which API paths must resolve from Lambda.
4. Confirm whether any legacy route aliases are still required during migration.
5. Confirm artifact naming, SSM paths, and resource naming for the temporary environment.

Deliverables:

- final route map,
- runtime variable list,
- temporary environment naming plan.

### Phase 2: Add the Go Lambda entrypoint

Goal: make the backend deployable in serverless form.

Tasks:

1. Add `apps/server/cmd/lambda/main.go`.
2. Reuse the existing transport and wiring so Lambda and HTTP server share the same application composition.
3. Ensure Lambda startup reads the same configuration model as the local server.
4. Keep runtime-specific code isolated to the Lambda command package.
5. Add or update backend tests if the new entrypoint introduces behavior differences.

Validation:

- `make backend-guard`
- focused transport tests if needed

### Phase 3: Create the temporary serverless environment

Goal: provision a parallel environment without mutating the old production stack.

Tasks:

1. Create a new Terraform environment directory, for example `infra/environments/serverless-next`.
2. Reuse the same modules where possible, but parameterize them for the new architecture.
3. Change the CloudFront routing model so S3 is the default origin and Go Lambda only handles `/api/*` and `/health`.
4. Update Lambda runtime defaults to Go for the new environment.
5. Publish outputs to a separate SSM prefix.
6. Keep budget and cost guardrails enabled.

Validation:

- `terraform validate`
- `terraform plan`
- manual review of origin behaviors, bucket names, Lambda runtime, and SSM paths

### Phase 4: Update CI/CD for the new deployment model

Goal: make deployment automation match the current application architecture.

Tasks:

1. Update the deploy workflow to detect backend changes as deployment-relevant.
2. Build Astro static output from `apps/web`.
3. Build the Go Lambda binary from `apps/server`.
4. Package the Go Lambda artifact and upload it to the artifacts bucket.
5. Sync static assets to the temporary environment S3 bucket.
6. Update the Go Lambda function code.
7. Invalidate the temporary CloudFront distribution.
8. Ensure migrations run before deployment where appropriate.

Validation:

- workflow dry review
- manual `workflow_dispatch`
- successful artifact upload and Lambda update

### Phase 5: End-to-end validation in the temporary environment

Goal: prove the new stack works as a whole.

Tasks:

1. Verify the login flow against `/api/auth/login`.
2. Verify auth bootstrap against `/api/auth/session`.
3. Verify protected page loading.
4. Verify dashboard reads for calories, expenses, and heat.
5. Verify at least one mutation flow per domain.
6. Verify logout.
7. Verify CloudFront route behavior for static pages, assets, API, and health.
8. Verify runtime configuration and origin protection behavior.

Recommended checks:

- `make backend-guard`
- `pnpm test` in `apps/web`
- `pnpm test:e2e` in `apps/web`
- manual browser smoke test against the temporary CloudFront URL

Exit criteria:

- no old-route dependencies remain in the deployed frontend,
- API and auth behavior match local expectations,
- no critical CloudFront routing or caching errors remain.

### Phase 6: Cutover strategy

Goal: replace the old production stack with the validated architecture.

Preferred sequence:

1. Freeze non-essential changes.
2. Re-run validation on the temporary environment.
3. Choose cutover method:
   - promote the temporary stack by switching DNS/domain, or
   - recreate the final production stack from the now-proven configuration.
4. Run post-cutover smoke tests immediately.
5. Keep the old stack available briefly as rollback insurance.

Rollback rule:

- If login, auth bootstrap, or primary dashboard flows fail after cutover, revert traffic to the old stack immediately.

### Phase 7: Cleanup

Goal: remove obsolete infrastructure only after successful cutover.

Tasks:

1. Destroy the old Astro SSR-style environment once the new stack has soaked successfully.
2. Remove old deployment assumptions from Terraform and CI/CD.
3. Delete obsolete SSM paths, Lambda artifacts, and route assumptions.
4. Update repository docs to reflect the final steady state.

Docs to update during or after completion:

- `docs/INFRA.md`
- `docs/ARCHITECTURE.md`
- `docs/APPLICATION.md`
- `docs/TESTING.md`
- `docs/DECISIONS.md` if the temporary-environment cutover choice should be recorded as an architecture decision

## Recommended Work Order

The most efficient sequence is:

1. Add Go Lambda entrypoint.
2. Create temporary Terraform environment.
3. Update deployment workflow.
4. Deploy temporary environment.
5. Run automated and manual validation.
6. Cut over.
7. Remove the old stack.

## Risks and Mitigations

### Risk: mixed frontend/backend contract

Mitigation:

- never point the old deployed frontend at the new Go API as an intermediate step,
- deploy the new frontend and new backend together in the temporary environment.

### Risk: CloudFront routing mistakes

Mitigation:

- validate path patterns explicitly,
- keep `/api/*` and `/health` isolated to Lambda,
- keep frontend routes and static assets isolated to S3.

### Risk: rollback difficulty after asset deletion

Mitigation:

- do not mutate the old environment during validation,
- keep rollback as a traffic decision rather than a rebuild decision.

### Risk: infrastructure drift during migration

Mitigation:

- introduce the new environment separately,
- only collapse or replace old resources after the new environment is proven.

## Immediate Next Actions

1. Design `apps/server/cmd/lambda/main.go`.
2. Add a parallel `serverless-next` Terraform environment.
3. Update the deploy workflow to build Go instead of Astro SSR for Lambda.
4. Deploy and validate the temporary environment.

## Success Criteria

This migration is complete when:

- the deployed frontend is the current static Astro app,
- the deployed backend is the current Go API,
- the temporary environment validates all critical user flows,
- production cutover completes without depending on old route contracts,
- the old Astro SSR stack is retired.