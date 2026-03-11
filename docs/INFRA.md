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
- Lambda artifact upload (optional)

Notes:
- Function URL uses IAM auth and is invoked via CloudFront.
- Runtime config reads Turso secrets from SSM paths.

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

Local init:
- Uses a local `backend.hcl` (ignored by git) for state backend config.

## Resource Inventory (Serverless)

- 1 CloudFront distribution (frontend assets + API origin)
- 2 S3 buckets (assets, artifacts)
- 1 Go Lambda function + Function URL
- 1 IAM role for Lambda execution
- SSM parameters for origin verification and deployment outputs
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

## CI/CD Touchpoints

- `.github/workflows/terraform-serverless.yml` runs Terraform with OIDC.
- `.github/workflows/deploy-serverless.yml` builds static Astro assets, uploads them to S3, updates the Go Lambda, and invalidates CloudFront.
