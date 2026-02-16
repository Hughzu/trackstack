# Step 0.5: AWS SSR Deployment (Astro Node Adapter)

**Goal:** Deploy the current Astro SSR app on AWS as a temporary, production-like environment for daily use and UX iteration. This is a stopgap before Step 1 (Go backend + Lambda).

**Constraints:**
- Keep it cheap and minimal.
- Infrastructure must be Terraform-managed.
- Preserve current Astro SSR architecture (Node adapter + API routes + middleware).

---

## 1) Architecture (What We Are Building)

- **CloudFront** as the public edge.
- **Lambda (Node SSR)** for SSR + API routes + auth middleware.
- **S3** for static assets (Astro build output).
- **Secrets/Env** in Lambda for Turso + auth config.

Flow:
- `GET /` or any SSR route -> CloudFront -> Lambda (SSR)
- `GET /assets/*` or static files -> CloudFront -> S3
- `POST /api/*` -> CloudFront -> Lambda (API routes)

---

## 2) Build Outputs (Astro)

- Use current config: `output: "server"` with `@astrojs/node`.
- Build output includes:
  - `dist/server/` (Node SSR entry)
  - `dist/client/` (static assets)

---

## 3) AWS Resources (Terraform Scope)

### Core
- **S3 bucket** for static assets
- **CloudFront distribution**
  - Default origin: Lambda SSR (Function URL or Lambda@Edge, see section 4)
  - Additional origin: S3 static assets
  - Behaviors:
    - `/assets/*` -> S3
    - `/_astro/*` -> S3 (if present)
    - `/*` -> Lambda SSR
- **Lambda function** for SSR (Node 18/20)
- **IAM roles/policies** for Lambda + CloudFront access
- **ACM certificate** (for your domain; optional for Step 0.5)
- **Route 53 record** (optional)

### Observability (minimal)
- CloudWatch Logs for Lambda
- Optional: basic alarms (errors > 0)

---

## 4) Lambda Integration Options

### Option A (Simpler): CloudFront -> Lambda Function URL
- Use Lambda Function URL as CloudFront origin.
- Pros: simpler setup, no Lambda@Edge packaging.
- Cons: slightly less optimal edge latency.

### Option B (More AWS-native): Lambda@Edge
- SSR code deployed to us-east-1 and attached to CloudFront viewer request.
- Pros: classic SSR pattern, edge compute.
- Cons: more complex deployment, slower iteration.

**Recommendation for Step 0.5:** Option A (Function URL). It is faster to set up and iterate.

---

## 5) Deployment Flow (Manual Steps)

1. **Build the app**
   - `pnpm install`
   - `pnpm build`

2. **Upload static assets to S3**
   - Sync `dist/client/` to S3
   - Keep cache-control aggressive for fingerprinted files

3. **Package SSR for Lambda**
   - Zip `dist/server/` + required Node runtime files
   - Ensure `node_modules` is included (or bundle for Lambda)
   - Entry point: `dist/server/entry.mjs`

4. **Set Lambda environment variables**
   - `TURSO_USERS_URL`, `TURSO_USERS_TOKEN`
   - `TURSO_CALORIES_URL`, `TURSO_CALORIES_TOKEN`
   - `TURSO_EXPENSES_URL`, `TURSO_EXPENSES_TOKEN`
   - `TURSO_HEAT_URL`, `TURSO_HEAT_TOKEN`
   - `AUTH_COOKIE_SECURE=true`
   - `AUTH_COOKIE_SAMESITE=lax` (or `strict`)
   - `AUTH_COOKIE_NAME=session`
   - `AUTH_SESSION_*` values if needed

5. **Configure CloudFront**
   - Default behavior -> Lambda origin
   - Static behaviors -> S3 origin
   - Forward headers/cookies for SSR and auth
   - Enable HTTPS

6. **Test**
   - Login flow
   - API routes
   - Session cookie persistence
   - Static assets caching

---

## 6) Terraform Implementation Outline

### Modules/Resources
- `aws_s3_bucket` + `aws_s3_bucket_policy`
- `aws_cloudfront_distribution`
- `aws_lambda_function`
- `aws_iam_role` + `aws_iam_policy` (Lambda execution)
- `aws_lambda_permission` (Function URL + CloudFront)
- `aws_acm_certificate` (optional)
- `aws_route53_record` (optional)

### Notes
- Lambda Function URL needs IAM or `NONE` auth; if `NONE`, lock access to CloudFront origin only.
- CloudFront must forward `Cookie`, `Authorization` (if used), and `Host` headers.
- Ensure `Set-Cookie` passes back through CloudFront.

---

## 6.1) OIDC Setup for Terraform (GitHub Actions)

Use the archive scripts as the starting point, but tighten them for least privilege.

### What worked before (archive)
- OIDC provider: `archive/iac/common/scripts/create-open-id-connector-provider.sh`
- Terraform role: `archive/iac/common/scripts/create-oidc-iac-role.sh`
- TF state bucket: `archive/iac/common/scripts/create-s3-tfstate-bucket.sh`
- CI workflow: `archive/.github/workflows/deploy-blog-iac.yml`

### What to change (hardening)
- **Do not use `AdministratorAccess`** for Terraform in Step 0.5.
  - Replace with a scoped policy limited to the resources you create (S3, CloudFront, Lambda, IAM for Lambda only).
- **Restrict OIDC `sub`** to your repo and branch.
  - Example: `repo:Hughzu/trackstack:ref:refs/heads/main`.
- **Use environment protection** in GitHub Actions for `apply`.
  - Manual approval on `prod`/`step-0-5` environment.

### OIDC Bootstrap Steps (manual, one-time)
1. Create OIDC provider (as in archive script).
2. Create IAM role for Terraform with a least-privilege policy.
3. Store role ARN as GitHub secret: `AWS_TERRAFORM_ROLE_ARN`.
4. Use `aws-actions/configure-aws-credentials` in CI with OIDC.

---

## 6.2) Terraform State (Safe Remote Backend)

Use an S3 backend with DynamoDB locking. The archive script creates the bucket but does not configure DynamoDB locks.

Minimum:
- **S3 bucket** with versioning + encryption + public access block.
- **DynamoDB table** for state locking (PK: `LockID`).
- Backend config:
  - `bucket`: `trackstack-terraform-state`
  - `key`: `step-0-5/terraform.tfstate`
  - `region`: your AWS region
  - `dynamodb_table`: `trackstack-terraform-locks`

---

## 6.3) CI/CD Pattern (Safe Deploy)

Adapt from `archive/.github/workflows/deploy-blog-iac.yml`:

- **Plan on PRs and push** to `main`.
- **Apply only on manual approval** (workflow_dispatch) or protected environment.
- **Destroy only via manual workflow**.

Recommended tweaks:
- Pin Terraform >= 1.5 (current script uses ~1.0).
- Separate `plan` and `apply` jobs.
- Upload `tfplan` as artifact, only apply that exact plan.

**DB migrations:** `/.github/workflows/db-migrations.yml` currently auto-applies on push to `main`. This is risky if infra deploys require manual approval, because schema changes can run ahead of infra or code. Rework it to run `atlas migrate diff`/`plan` on push (or PR) and only apply via manual workflow_dispatch or a protected environment gate. Align its trigger with the new CI/CD flow: apply migrations only after infra/app deploy approval and only from `main`.

---

## 7) Security Checklist (Minimum)

- `.env` is local only, never committed.
- Rotate Turso tokens if you’ve shared or logged them.
- Use HTTPS only; `AUTH_COOKIE_SECURE=true` in prod.
- Restrict Lambda URL to CloudFront (origin access control or IP allowlist).
- Consider adding CSRF protection later (optional for single-user Step 0.5).

---

## 7.1) Challenge Notes (From the Archive)

- **AdminAccess on the Terraform role is too broad.** Use least privilege.
- **OIDC `sub` should be locked to `main`** (not `repo:*`).
- **State locking is missing.** Add DynamoDB table.
- **Auto-apply on push is risky.** Prefer manual apply.

---

## 8) What This Enables

- Daily use from any device
- Fast iteration without touching Go or full Golden Paths
- Clean exit to Step 1 (Go backend + Lambda) later

---

## 9) Step 1 Transition (Later)

- Replace Lambda SSR backend with Go API (Auth + DB)
- Keep CloudFront + S3 for static assets
- Migrate auth/session handling to Go

---

## 10) Open Decisions

- Domain + DNS (Route 53 or external provider)
- Function URL vs Lambda@Edge
- Cache policy for SSR responses
