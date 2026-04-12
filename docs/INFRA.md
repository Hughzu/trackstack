# Infrastructure

This document explains the infrastructure layout, intent, and operational boundaries. It is designed to be a constraint for future changes, not just a description of the current state.

## Philosophy

- Bootstrap is one-time, account-level setup.
- Modules are reusable building blocks.
- Environments compose modules with minimal glue.
- Serverless is the production baseline with scale-to-zero and tight cost controls.

## Infrastructure Rules

- Bootstrap changes live in `infra/bootstrap/` (account-level, one-time setup).
- Shared infrastructure components live in `infra/modules/`.
- Environment wiring lives in `infra/environments/*/` and should stay thin.
- CI/CD must use OIDC only; no long-lived AWS keys in GitHub secrets.
- Serverless is production and must keep scale-to-zero and low-cost defaults.
- Lab environments must include automated destroy/kill-switches.
- Runtime secrets come from SSM; do not hardcode or inline secrets in Terraform.
- Cost guardrails (budget alert) must remain enabled for production.

## Directory Overview

```
infra/
  bootstrap/
    bootstrap/     # One-time scripts for AWS account setup + OIDC
    iam/           # Terraform for the GitHub deploy role policy
  modules/         # Reusable Terraform components
  environments/    # Composition for each environment
```

## bootstrap/

### Purpose

Account-level setup that must happen before any environment can be deployed. This is intentionally separate from the app stacks to keep credential/bootstrap concerns isolated.

### bootstrap/bootstrap (scripts)

Location: `infra/bootstrap/bootstrap`

What it does:
- Creates the admin IAM user for local CLI use.
- Creates the Terraform state S3 bucket and DynamoDB lock table.
- Creates the GitHub OIDC provider and role.
- Pushes GitHub repo secrets (role ARN, region, state bucket, lock table).
- Syncs the deploy role permissions via Terraform.

Key design choice: GitHub Actions uses OIDC, not long-lived access keys.

### bootstrap/iam (Terraform)

Location: `infra/bootstrap/iam`

What it does:
- Defines the IAM policy attached to the GitHub deploy role.
- Scopes write permissions to Trackstack resources where possible.
- Adds broad read/list permissions needed by Terraform plans.

This is the only place we should grant or tighten CI/CD permissions. All environment stacks assume this policy is already in place.

Local init:
- Uses a local `backend.hcl` (ignored by git) for state backend config.

## modules/

Modules are reusable building blocks shared across environments. They should be stable and environment-agnostic.

### lambda-api

Location: `infra/modules/lambda-api`

Responsibilities:
- Lambda function + Function URL
- IAM execution role and least-privilege policies
- Log group with retention
- Origin verification secret stored in SSM
- Optional first-apply Lambda bootstrap from a local zip
- Runtime environment hydration from SSM-backed parameters

Notes:
- Function URL uses IAM auth and is invoked via CloudFront.
- Runtime config for the parallel environment is sourced from SSM parameters and injected into Lambda environment variables during apply.
- Bootstrap Lambda artifacts are uploaded directly to Lambda during Terraform apply; deploy artifacts in S3 are not tracked as Terraform-managed objects.

### static-hosting

Location: `infra/modules/static-hosting`

Responsibilities:
- Assets S3 bucket (SSE, public access block, ownership)
- CloudFront distribution
- OAC for S3 and Lambda
- Security headers policy
- Bucket policy restricted to the CloudFront distribution
- Lambda permissions for CloudFront invocation

Notes:
- Price class defaults to `PriceClass_100` to reduce cost.
- Origin verification header is passed to Lambda origin.
- Supports both legacy Lambda-default routing and static-first S3-default routing via module inputs.
- Can rewrite SPA frontend routes without file extensions back to the shared `/index.html` shell through a CloudFront Function.

### cost-guardrails

Location: `infra/modules/cost-guardrails`

Responsibilities:
- AWS Budgets monthly cost guardrail
- Email notification at 90% of the budget limit

Notes:
- Budgets is global and must run in `us-east-1`.
- Set `billing_alarm_email` to enable.

## environments/

### serverless

Location: `infra/environments/serverless`

Purpose:
- Production environment with scale-to-zero.
- Composes the modules and keeps only environment glue in the root.

Composition:
- `module.lambda_api`: Go API Lambda + IAM + SSM origin secret.
- `module.static_hosting`: S3 frontend assets + CloudFront + OAC.
- `module.cost_guardrails`: monthly budget.
- `s3.tf`: artifacts bucket used for Lambda package storage.
- `ssm.tf`: publishes infra outputs to SSM for CI/CD.

### serverless-next

Location: `infra/environments/serverless-next`

Purpose:
- Legacy migration validation environment for the static frontend + Go API serverless contract.
- Safe to destroy after the production `serverless` cutover is complete.

Composition:
- `module.lambda_api`: Go API Lambda using the custom Go runtime artifact.
- `module.static_hosting`: S3 default origin, Lambda only for `/api/*`, `/health`, and `/openapi.yaml`.
- `module.cost_guardrails`: monthly budget.
- `s3.tf`: artifacts bucket used for Lambda package storage.
- `ssm.tf`: publishes infra outputs to a dedicated `/trackstack/serverless-next` prefix.
- `01-set-runtime-ssm.sh`: seeds the serverless-next rebuild runtime into SSM with `APP_ENV`, `LOG_LEVEL`, `DB_CONNECTION_MODE`, `CORS_ALLOWED_ORIGIN`, `JWT_SECRET`, and `TURSO_*_{URL_HTTP,TOKEN}` values.

Local init:
- Uses a local `backend.hcl` (ignored by git) for state backend config.

## Resource Inventory (Serverless)

- 1 CloudFront distribution (frontend assets + API origin)
- 2 S3 buckets (assets, artifacts)
- 1 Go Lambda function + Function URL
- 1 IAM role for Lambda execution
- SSM parameters for origin verification and deployment outputs
- AWS Budget for cost guardrail

## Resource Inventory (Serverless-Next)

- 1 CloudFront distribution with S3 as default origin
- 2 S3 buckets (assets, artifacts)
- 1 Go Lambda function + Function URL for API and health traffic
- 1 IAM role for Lambda execution
- SSM parameters for runtime config, origin verification, and deployment outputs
- AWS Budget for cost guardrail

## Security Boundaries

- GitHub deploy role is scoped to Trackstack resources and OIDC is restricted to the main branch.
- Lambda execution role can only write logs and read runtime SSM parameters.
- Assets bucket is not public and is only accessible via CloudFront OAC.
- Lambda Function URL is IAM-authenticated and restricted to CloudFront. Because CloudFront uses `Authorization` for SigV4 signing to the Function URL origin, browser bearer tokens must be forwarded in `X-Trackstack-Authorization` instead of `Authorization` on deployed app requests.

### Lambda URL SigV4 Contract

- CloudFront OAC signs origin requests to the Lambda Function URL with SigV4 against the Lambda URL host, not the public CloudFront host.
- Browser requests must treat `Authorization` as reserved for CloudFront signing. App auth travels in `X-Trackstack-Authorization` so CloudFront can keep its own SigV4 header intact.
- Browser write requests that travel through the CloudFront -> Lambda Function URL path must include `x-amz-content-sha256` for the exact final request body.
- Empty-body write requests still need the SHA256 of the empty payload. AWS Lambda Function URLs do not accept unsigned payloads for `POST`, `PUT`, or `PATCH`, and missing hashes can also break other write-style requests routed through the same signed path.
- When this header is missing or wrong, AWS rejects the request before app code runs. The usual symptom is `403 InvalidSignatureException` with a message about the calculated request signature not matching.
- `GET` and `HEAD` can still work while writes fail, so a healthy `/health` endpoint does not prove login or mutations are wired correctly.

## FinOps Guardrails

- CloudFront price class set to `PriceClass_100`.
- Lambda memory size set to 512 MB and log retention to 3 days.
- Assets bucket versioning is suspended (immutable assets).
- Artifacts bucket lifecycle expires old versions after 7 days.
- Monthly budget alert via AWS Budgets.

## Manual Steps (Not in Terraform)

- Run bootstrap scripts in `infra/bootstrap/bootstrap` once per account.
- Create local `backend.hcl` files for Terraform init (not committed).
- Seed production runtime config with `infra/environments/serverless/01-set-runtime-ssm.sh` before the first Terraform plan/apply; it should align with `apps/server/.env`, auto-load `JWT_SECRET` plus `TURSO_*_URL_HTTP` and `TURSO_*_TOKEN` when present, and exported env vars still win.
- If you still need the old migration stack, seed its runtime config with `infra/environments/serverless-next/01-set-runtime-ssm.sh` before plan/apply.
- The shared Lambda module no longer injects legacy app-specific runtime defaults automatically; each environment must define its own runtime contract explicitly.

## Local Compose

- The local frontend compose service is `web-frontend` and serves the Solid/Vite app from `http://localhost:5173`.
- `web-frontend` keeps `node_modules` and the pnpm store in named volumes so bind-mounting `apps/web` does not nuke installed dependencies on container start.
- Local browser-to-Go traffic stays same-origin through the Vite `/api` proxy; set `API_PROXY_URL` for the container network target and use `VITE_API_BASE_URL` only for direct static deployment origins.

## CI/CD Touchpoints

- `.github/workflows/terraform-serverless.yml` manages the production `serverless` environment with OIDC, can build a bootstrap Go Lambda artifact for first apply, and only runs its job when the ref is `main`.
- `.github/workflows/deploy-serverless.yml` deploys the production `serverless` environment: it builds the Solid/Vite frontend in `apps/web`, packages the Go Lambda custom runtime artifact, uploads both to S3, updates Lambda, invalidates CloudFront, and only runs jobs when the ref is `main`.
