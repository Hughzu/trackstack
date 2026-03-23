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
- Can rewrite extensionless frontend paths to `index.html` object keys through a CloudFront Function for static Astro deployments.

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
- Parallel validation environment for the static Astro + Go API serverless contract.
- Keeps the existing production serverless environment untouched during migration testing.

Composition:
- `module.lambda_api`: Go API Lambda using the custom Go runtime artifact.
- `module.static_hosting`: S3 default origin, Lambda only for `/api/*`, `/health`, and `/openapi.yaml`.
- `module.cost_guardrails`: monthly budget.
- `s3.tf`: artifacts bucket used for Lambda package storage.
- `ssm.tf`: publishes infra outputs to a dedicated `/trackstack/serverless-next` prefix.
- `01-set-runtime-ssm.sh`: seeds the serverless-next runtime contract into SSM and now needs to provide `JWT_SECRET` for the bearer-auth rebuild runtime.

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
- Lambda Function URL is IAM-authenticated and restricted to CloudFront.

## FinOps Guardrails

- CloudFront price class set to `PriceClass_100`.
- Lambda memory size set to 512 MB and log retention to 3 days.
- Assets bucket versioning is suspended (immutable assets).
- Artifacts bucket lifecycle expires old versions after 7 days.
- Monthly budget alert via AWS Budgets.

## Manual Steps (Not in Terraform)

- Run bootstrap scripts in `infra/bootstrap/bootstrap` once per account.
- Create local `backend.hcl` files for Terraform init (not committed).
- Set runtime secrets in SSM using `infra/environments/serverless/01-set-runtime-ssm.sh`.
- For the migration environment, seed runtime config with `infra/environments/serverless-next/01-set-runtime-ssm.sh` before deploy validation; it should align with `apps/server-next/.env`, auto-load `TURSO_*_URL_HTTP` and `TURSO_*_TOKEN` when present, and exported env vars still win.
- The migration runtime seed script must now include `JWT_SECRET`; legacy cookie/session settings are no longer part of the `apps/server-next` runtime contract.

## CI/CD Touchpoints

- `.github/workflows/terraform-serverless.yml` currently manages the temporary `serverless-next` environment with OIDC, can build a bootstrap Go Lambda artifact for first apply, and only runs its job when the ref is `main`.
- `.github/workflows/deploy-serverless.yml` currently deploys the temporary `serverless-next` environment: it builds static Astro assets, packages the Go Lambda custom runtime artifact, uploads both to S3, updates Lambda, invalidates CloudFront, and only runs jobs when the ref is `main`.
