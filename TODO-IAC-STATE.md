# Trackstack IaC Current State

This document summarizes the current infrastructure shape for `step-0-5-astro-ssr` and where code lives in the repo.

## Environment: `iac/environments/step-0-5-astro-ssr`

### High-level topology
- **CloudFront** distribution with two origins:
  - **Lambda Function URL** for SSR and API routes (default behavior).
  - **S3** for static assets (`/_astro/*`, `/assets/*`).
- **Lambda** hosts Astro SSR (Function URL, AWS_IAM, OAC SigV4).
- **S3** serves built client assets via OAC.
- **SSM** stores secrets (origin verify header value, Turso URLs/tokens).

### Key resources (by file)
- `cloudfront.tf`
  - Two origins: `lambda-ssr` and `s3-assets`.
  - OAC for Lambda and S3.
  - Default behavior -> Lambda, assets behaviors -> S3.
  - Origin request policy: `Managed-AllViewerExceptHostHeader` for Lambda.
  - Custom origin header `X-Origin-Verify` from SSM secret.

- `lambda.tf`
  - `aws_lambda_function` with env vars for origin verify + Turso.
  - `aws_lambda_function_url` with `AWS_IAM`.
  - Permissions for CloudFront to invoke the Function URL.

- `s3.tf`
  - Assets bucket for static client assets.

- `artifacts.tf`
  - Uploads the Lambda artifact when `lambda_artifact_path` is provided.

- `ssm.tf`
  - Stores the Lambda artifact key and origin secret.

- `outputs.tf`
  - Exposes CloudFront URL and key infra outputs.

## CI/CD Flow (Step 0.5)

Workflow: `.github/workflows/terraform-step-0-5-astro-ssr.yml`

- Optional build step creates `/tmp/astro-ssr.zip` from `src/web/dist/server`.
- Terraform plan/apply uses `lambda_artifact_path` to upload the new artifact.

## Current Security Model

- Function URL is `AWS_IAM`.
- CloudFront OAC signs origin requests.
- Requests with bodies must include `x-amz-content-sha256` (handled by client `signedFetch`).
- Additional origin guard: `X-Origin-Verify` header validated by Astro middleware.

## Code Boundaries: Where Things Go

This repo uses the **Managed Polylith** architecture (see `AGENTS.md`).

### Backend (Go)
- `internal/core/` **Managed**: Auth, DB setup, server plumbing. Do not modify unless explicitly asked.
- `internal/modules/` **Product**: Business logic for features (Calories, Heat, Expenses).
- `cmd/monolith/` **Assembly**: The universal binary that wires modules.

### Frontend (Astro)
- `src/web/src/core/` **Managed UI**: layouts, shared UI kit, design system.
- `src/web/src/modules/` **Product UI**: feature pages/components for modules.
- `src/web/src/pages/` **Routes**: SSR pages and API endpoints.

### Infrastructure (Terraform)
- `iac/_vendor/` **Managed**: base modules.
- `iac/environments/` **User**: environment-specific config (this is where `step-0-5-astro-ssr` lives).

## Notes / Gotchas

- Lambda handler must be `entry.handler` (zip root contains `entry.mjs`).
- Use `@astro-aws/adapter` for Lambda SSR (node standalone adapter is not Lambda-friendly).
- LibSQL in Lambda uses `@libsql/client/web` (no native binaries).
