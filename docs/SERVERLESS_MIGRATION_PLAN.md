# Serverless Migration Debug Brief

Date: 2026-03-13
Status: Active debugging
Scope: `serverless-next` temporary environment only

## Goal

TrackStack is migrating from the old Astro SSR-style serverless shape to:

- static Astro frontend in S3,
- Go API in Lambda,
- CloudFront as the only public entrypoint.

This file is optimized to give another agent enough context to debug the current `serverless-next` environment quickly.

## Intended `serverless-next` Architecture

- CloudFront default origin: S3
- Lambda origin only for:
  - `/api/*`
  - `/health`
  - `/openapi.yaml`
- Static Astro routes and assets served from S3 via CloudFront
- Extensionless frontend routes rewritten to `index.html` object keys by a CloudFront Function
- Go Lambda runtime config resolved from SSM-backed parameters

Key files:

- `infra/environments/serverless-next/main.tf`
- `infra/modules/static-hosting/main.tf`
- `infra/modules/lambda-api/main.tf`
- `.github/workflows/terraform-serverless.yml`
- `.github/workflows/deploy-serverless.yml`

## What Was Already Completed

### Backend runtime

- Added Lambda entrypoint at `apps/server/cmd/lambda/main.go`
- Shared HTTP/Lambda bootstrap now lives in `apps/server/internal/core/app/app.go`
- `make backend-guard` passed after that change

### Terraform

- Added `infra/environments/serverless-next/`
- Reused and extended shared modules instead of creating one-off infra:
  - `infra/modules/static-hosting`
  - `infra/modules/lambda-api`
  - `infra/modules/cost-guardrails`
- `serverless-next` uses:
  - `resource_prefix = trackstack-next`
  - `ssm_prefix = /trackstack/serverless-next`
  - `lambda_function_name = trackstack-go-api-next`
- `terraform validate` passed locally for `infra/environments/serverless-next`
- GitHub Terraform workflow apply passed for `serverless-next`

### CI/CD

- Terraform workflow now targets `serverless-next` and can build a bootstrap Go Lambda zip:
  - `.github/workflows/terraform-serverless.yml`
- Deploy workflow now targets `serverless-next` and:
  - builds static Astro from `apps/web`
  - builds Go Lambda from `apps/server/cmd/lambda`
  - uploads static files to S3
  - updates Lambda code
  - invalidates CloudFront
  - file: `.github/workflows/deploy-serverless.yml`

## Current Observed Problems

### Problem 1: CloudFront root returns AccessDenied

Observed after Terraform apply on the new `serverless-next` environment:

```xml
<Error>
  <Code>AccessDenied</Code>
  <Message>Access Denied</Message>
</Error>
```

Important context:

- The Terraform workflow creates infra and may bootstrap the Lambda artifact.
- The Terraform workflow does not upload the Astro static site to the assets bucket.
- The static site upload happens in `.github/workflows/deploy-serverless.yml`, not in Terraform.
- `serverless-next` uses S3 as the default CloudFront origin.

This means the most likely immediate explanation is:

- CloudFront is correctly pointing to S3,
- but the assets bucket does not yet contain `index.html` and the rest of the built Astro output,
- so CloudFront root requests resolve to S3 and fail with `AccessDenied`/missing-object behavior.

But this still needs confirmation, because it could also be:

- S3 bucket policy/OAC mismatch,
- wrong CloudFront behavior ordering,
- bad default root object / rewrite behavior,
- or missing deploy step after infra apply.

### Problem 2: Local Terraform plan wants to destroy bootstrap artifact object

Observed locally:

```text
  # module.lambda_api.aws_s3_object.lambda_artifact[0] will be destroyed
  # (because index [0] is out of range for count)
```

Plan summary:

```text
Plan: 0 to add, 0 to change, 1 to destroy.
```

Important context:

- `infra/modules/lambda-api/main.tf` creates `aws_s3_object.lambda_artifact` only when `var.artifact_path != null`
- the Terraform workflow passes `lambda_artifact_path=/tmp/go-api-lambda.zip` only when `bootstrap_artifact=true`
- later local plans normally do not pass `lambda_artifact_path`
- therefore Terraform state contains a bootstrap-only managed S3 object, but the current config no longer declares it

This means the most likely explanation is bootstrap drift:

- first apply created a managed temporary artifact object in S3,
- later plan omitted `lambda_artifact_path`,
- Terraform now wants to destroy only the tracked bootstrap object resource.

This looks expected from the current design, but it is a design smell and may not be what we want long term.

## Current Relevant Infrastructure Facts

### CloudFront routing

From `infra/environments/serverless-next/main.tf`:

- `default_origin = "s3"`
- `lambda_path_patterns = ["/api/*", "/health", "/openapi.yaml"]`
- `enable_directory_index_rewrite = true`

From `infra/modules/static-hosting/main.tf`:

- `default_root_object = "index.html"` when S3 is default
- viewer-request CloudFront Function rewrites:
  - `/` -> `/index.html`
  - `/foo/` -> `/foo/index.html`
  - `/foo` -> `/foo/index.html`
- assets bucket policy allows `s3:GetObject` only to the specific CloudFront distribution ARN

### Lambda artifact behavior

From `infra/modules/lambda-api/main.tf`:

- Lambda always points at:
  - `s3_bucket = var.artifact_bucket`
  - `s3_key = var.artifact_key`
- `aws_s3_object.lambda_artifact` is only created when `artifact_path` is provided
- Lambda ignores changes to:
  - `s3_key`
  - `s3_object_version`

### Workflow split

Terraform workflow:

- file: `.github/workflows/terraform-serverless.yml`
- purpose: infra creation/update
- optional first-apply bootstrap of Go Lambda artifact
- does not upload Astro static site

Deploy workflow:

- file: `.github/workflows/deploy-serverless.yml`
- purpose: app deployment
- uploads Astro build output to assets bucket
- uploads Go Lambda artifact to artifacts bucket
- updates Lambda code
- invalidates CloudFront

## Most Likely Debug Directions

### For CloudFront AccessDenied

Check these in order:

1. Was `.github/workflows/deploy-serverless.yml` run successfully after Terraform apply?
2. Does the assets bucket contain:
   - `index.html`
   - `login/index.html`
   - `calories/index.html`
   - `expenses/index.html`
   - `heat/index.html`
   - `/_astro/*`
3. Does the CloudFront distribution default origin point to the assets bucket, not Lambda?
4. Does the assets bucket policy match the actual CloudFront distribution ARN?
5. Does `https://<cloudfront>/index.html` also return `AccessDenied`, or only `/`?
6. Do `/api/health` or `/health` reach Lambda successfully?

### For Terraform bootstrap drift

Check these in order:

1. Was the first apply run with `bootstrap_artifact=true`?
2. Are later plans intentionally omitting `lambda_artifact_path`?
3. Should the bootstrap object be treated as:
   - a one-time temporary artifact not managed after first apply, or
   - a persistent Terraform-managed object?
4. If it is one-time-only, should the module avoid managing that uploaded object in state at all?
5. If it should stay managed, should CI always pass `lambda_artifact_path`, or should artifact versioning be modeled differently?

## Recommended Investigation Commands

Useful Terraform checks:

```bash
terraform -chdir=infra/environments/serverless-next validate
terraform -chdir=infra/environments/serverless-next plan
terraform -chdir=infra/environments/serverless-next state show module.static_hosting.aws_cloudfront_distribution.ssr
terraform -chdir=infra/environments/serverless-next state show module.static_hosting.aws_s3_bucket_policy.assets
terraform -chdir=infra/environments/serverless-next state show module.lambda_api.aws_lambda_function.ssr
```

Useful AWS checks:

```bash
aws s3 ls s3://trackstack-next-assets-<account-id>/
aws s3 ls s3://trackstack-next-assets-<account-id>/login/
aws cloudfront get-distribution --id <distribution-id>
aws lambda get-function --function-name trackstack-go-api-next
curl -i https://<cloudfront-domain>/
curl -i https://<cloudfront-domain>/index.html
curl -i https://<cloudfront-domain>/health
curl -i https://<cloudfront-domain>/api/auth/session
```

## Probable Short Answer

If Terraform apply succeeded but the app deploy workflow has not run yet, then the `AccessDenied` at the CloudFront URL is probably expected because the S3 default origin has no uploaded static site yet.

Separately, the local plan wanting to destroy `module.lambda_api.aws_s3_object.lambda_artifact[0]` is probably caused by the current bootstrap-artifact design: first apply created a Terraform-managed temporary object, and later plans no longer declare it.

## What Another Agent Should Answer

1. Is the CloudFront `AccessDenied` simply because static assets were never deployed, or is there an infra policy/routing bug?
2. Is the bootstrap artifact drift acceptable, or should `infra/modules/lambda-api/main.tf` be changed so later plans are clean?
3. What is the safest minimal fix to make `serverless-next` stable before Phase 5 validation?
