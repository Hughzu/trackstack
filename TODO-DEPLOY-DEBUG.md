# Trackstack SSR Debug Log (Step 0.5 Astro SSR)

This is a concise record of the issues encountered and the fixes applied while bringing the CloudFront + Lambda Function URL Astro SSR stack online.

## Overview
We debugged the CloudFront distribution, Lambda Function URL authorization, Astro SSR handler setup, and packaging pipeline. Most issues were configuration mismatches or Lambda artifact contents not matching the handler or runtime expectations.

## Issues and Resolutions

### 1) CloudFront 403 from Origin Verify
**Symptom**: `403 Forbidden` from CloudFront, with response `{"Message":null}`.
**Cause**: Lambda middleware enforced an `X-Origin-Verify` header, but CloudFront was not reliably sending the expected header.
**Fix**:
- Ensure CloudFront origin custom header is set to `X-Origin-Verify` with SSM secret value.
- Ensure Lambda env has `ORIGIN_VERIFY_HEADER` and `ORIGIN_VERIFY_VALUE` populated.

### 2) CloudFront signing mismatch (OAC + Lambda Function URL)
**Symptom**: `InvalidSignatureException` on POST to `/api/auth/login` from CloudFront.
**Cause**: Lambda Function URL with `AWS_IAM` requires `x-amz-content-sha256` payload hash for requests with bodies. Browsers do not send this header by default; CloudFront OAC does not compute it for viewer payloads.
**Fix** (two options):
- **Preferred (kept AWS_IAM)**: add `x-amz-content-sha256` in client requests for any POST/DELETE with a body.
- **Alternative**: set Function URL auth to `NONE` and keep `X-Origin-Verify` as the gate (not chosen for final config).

### 3) Handler mismatch in Lambda
**Symptom**: `Runtime.ImportModuleError: Cannot find module 'entry'`.
**Cause**: Lambda handler set to `dist/server/entry.mjs` but the zip places `entry.mjs` at the archive root. Lambda handlers must be `file.export`, not a path.
**Fix**:
- Set handler to `entry.handler`.

### 4) Lambda init timeout in standalone adapter
**Symptom**: init timeout during Lambda cold start.
**Cause**: Astro `@astrojs/node` standalone adapter attempted to resolve `dist/server` path baked into `import.meta.url` (CI path), causing a loop during init.
**Fix**:
- Switched to `@astro-aws/adapter` for Lambda SSR (AWS-native), `mode: "ssr"`.

### 5) Missing `@libsql/client` at runtime
**Symptom**: `Cannot find package '@libsql/client'`.
**Cause**: SSR bundle did not include the dependency.
**Fix**:
- Bundle the dependency via Vite SSR config:
  - `vite.ssr.noExternal: ["@libsql/client", "@libsql/client/web", "@libsql/client/http"]`.

### 6) Missing native `@libsql/linux-x64-gnu`
**Symptom**: `Could not dynamically require "@libsql/linux-x64-gnu"`.
**Cause**: The Node client tries to load native binaries at runtime.
**Fix**:
- Use the web client for Lambda:
  - `import { createClient } from "@libsql/client/web"`.
- Remove file-based SQLite fallback in Lambda (Turso URLs are required).

### 7) Function URL auth drift + CloudFront origin mismatch
**Symptom**: All routes returned `Forbidden` from Function URL, no Lambda logs.
**Cause**: Function URL was deleted/recreated via CLI, generating a new URL; CloudFront origin still pointed to the old URL until `terraform apply`.
**Fix**:
- Reconcile with Terraform apply so CloudFront uses the current Function URL.

## Configuration Changes (IaC)

### CloudFront
- Use `Managed-AllViewerExceptHostHeader` origin request policy for Lambda OAC.
- Ensure Lambda origin custom header is `X-Origin-Verify` with SSM secret.

### Lambda Function URL
- Keep `authorization_type = "AWS_IAM"`.
- Keep both permissions:
  - `lambda:InvokeFunctionUrl` for CloudFront
  - `lambda:InvokeFunction` for CloudFront (safe to include)

### Handler
- `lambda_handler = "entry.handler"`

## Astro Changes

### Adapter
- Use `@astro-aws/adapter` with `mode: "ssr"`.
- Remove `@astrojs/node` adapter.

### LibSQL in Lambda
- Use `@libsql/client/web`.
- Enforce `TURSO_*_URL` presence (no file fallback in Lambda).

### Signed Fetch (POST/DELETE)
**Reason**: AWS_IAM Function URL requires payload hash for body requests.
**Change**: Global `window.signedFetch` computes `x-amz-content-sha256` and is used for all POST/DELETE fetch calls in the app, including `/api/auth/login`.

## CI/CD Notes

### Artifact Build
The GH Actions workflow builds and zips `dist/server`:

```
pnpm build
cd dist/server
zip -r /tmp/astro-ssr.zip .
```

### Deploy
Terraform uses `lambda_artifact_path` if provided to upload a fresh artifact:

```
terraform apply -var="lambda_artifact_path=/tmp/astro-ssr.zip"
```

### Sharp Warning
`@astro-aws/adapter` warns about `sharp` support. Build still completes; warning is safe to ignore unless image optimization is required.

## Current Expected State

- CloudFront in front of Lambda Function URL with OAC (SigV4)
- Lambda Function URL auth type: `AWS_IAM`
- Requests with bodies include `x-amz-content-sha256` (handled by `window.signedFetch`)
- Login works at `/login` and POSTs to `/api/auth/login` succeed

## Remaining Action Items

- Run `terraform apply` to reconcile any drift from CLI actions.
- Re-test login and all POST flows.
- Confirm CloudFront origin points to the active Function URL.
