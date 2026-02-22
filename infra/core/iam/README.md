# IAM Bootstrap Backend Config

Terraform backends are evaluated before variables and providers, so the backend
cannot read values from `variables.tf`. That is why the S3 backend only has a
static key in `main.tf` and expects the bucket/region/lock table at init time.

We keep a local `backend.hcl` file (ignored by git) to avoid hardcoding values
in commands and to keep account-specific details out of the repo.

Usage:

```bash
terraform init -backend-config=backend.hcl
```

If you need to change accounts or regions, update `backend.hcl` locally.
