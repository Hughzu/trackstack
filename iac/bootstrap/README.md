# AWS Bootstrap Scripts

This guide sets up secure CI/CD for TrackStack using OIDC (OpenID Connect) authentication with AWS. This is the foundation that enables GitHub Actions to deploy infrastructure without storing any long-lived AWS credentials.

## What is TrackStack?

TrackStack is a "Managed Polylith" application - a framework-based approach where:
- **Backend:** Go (Hexagonal Architecture)
- **Frontend:** Astro (Static Output) 
- **Database:** Turso (LibSQL over HTTP)
- **Deployment:** AWS (Lambda + CloudFront for serverless, EKS for scale)

This bootstrap creates the security foundation for automated deployments via GitHub Actions.

## Why OIDC Instead of Access Keys?

### The Problem with Long-Lived Credentials

Traditional approach: Store AWS access keys as GitHub secrets

❌ **Security risks:**
- Keys exist permanently in GitHub secrets
- If GitHub is compromised, attackers get persistent AWS access  
- Keys don't expire automatically
- No audit trail showing *which specific GitHub workflow run* used the credentials

✅ **The OIDC Solution:**
- **Zero stored credentials** - GitHub Actions authenticate via cryptographically signed tokens
- **Temporary access** - Credentials expire automatically (typically 1 hour)
- **Repository-scoped** - Only workflows from `Hughzu/trackstack` can authenticate
- **Full audit trail** - CloudTrail shows exact workflow runs

### How It Works

1. GitHub Action starts: "Hey AWS, I'm workflow XYZ from repository Hughzu/trackstack"
2. GitHub generates a JWT token (cryptographically signed, expires in 1 hour)
3. GitHub Action presents token to AWS
4. AWS validates: checks signature, expiration, and that the repo matches our trust policy
5. AWS issues temporary credentials valid only for this session
6. Workflow runs with temporary access
7. Credentials expire automatically - no cleanup needed

**Key benefit:** Even if someone intercepts the JWT token, it's only valid for minutes and can't be reused after expiration. No credentials to steal, rotate, or revoke.

---

## Prerequisites

Install the required tools:

| Tool | Purpose | Install Command | Verify |
|------|---------|-----------------|--------|
| AWS CLI | Interact with AWS | `sudo pacman -S aws-cli` | `aws --version` |
| GitHub CLI | Manage GitHub secrets | `sudo pacman -S github-cli` | `gh --version` |
| jq | Parse JSON output | `sudo pacman -S jq` | `jq --version` |

### GitHub CLI Authentication

The bootstrap scripts need to set secrets in your GitHub repository. First, authenticate the CLI:

```bash
gh auth login
# Select: GitHub.com
# Select: HTTPS  
# Select: Login with a web browser
# Follow the browser prompts to authorize

# Verify you're logged in:
gh auth status
```

---

## Bootstrap Execution

### Step 1: Create Admin IAM User (Requires Root Credentials)

**Why this step exists:** You need a dedicated IAM user for local CLI work. Using your root account for daily operations is a security risk. This user will have AdministratorAccess but is scoped to your specific project.

**Prerequisites:**
- You need root AWS account credentials (temporary, just for this step)

**Get root credentials:**
1. Log into AWS Console using your root account (the email tied to your Amazon account)
2. Go to IAM → Users → (your root user) → Security Credentials tab
3. Click "Create access key"
4. Select "Command Line Interface (CLI)"
5. Save the Access Key ID and Secret Access Key

**Configure root credentials (temporary):**
```bash
aws configure
# AWS Access Key ID: <paste root Access Key ID>
# AWS Secret Access Key: <paste root Secret Access Key>
# Default region name: eu-west-1
# Default output format: json
```

**Run the bootstrap script:**
```bash
# Ensure we're using the default (root) credentials, not any profile
unset AWS_PROFILE

# Navigate to the bootstrap directory
cd iac/bootstrap

# Run the admin user creation script
./01-create-admin-user.sh
```

**What happens:**
- Creates IAM user `trackstack-admin`
- Attaches AdministratorAccess policy
- Generates access keys for the new user
- Outputs the keys for you to save

**Save the credentials:**
```bash
aws configure --profile trackstack
# AWS Access Key ID: <paste the key from script output>
# AWS Secret Access Key: <paste the secret from script output>
# Default region name: eu-west-1
# Default output format: json
```

**IMPORTANT: Clean up root credentials:**
After saving the `trackstack` profile credentials, remove the root credentials from your AWS config:

```bash
# Edit the credentials file
nano ~/.aws/credentials

# Delete the [default] section or rename it to [root-temp]
# Only keep the [trackstack] section
```

**Verify the admin user works:**
```bash
export AWS_PROFILE=trackstack
aws sts get-caller-identity

# Expected output:
# {
#     "UserId": "AID...",
#     "Account": "939091506005",
#     "Arn": "arn:aws:iam::939091506005:user/trackstack-admin"
# }
```

---

### Steps 2-5: Bootstrap Infrastructure (Uses Admin Profile)

Now that you have the `trackstack` admin user configured, run the remaining scripts:

```bash
# Ensure you're using the trackstack profile
export AWS_PROFILE=trackstack

# Verify authentication
aws sts get-caller-identity

# Run remaining bootstrap scripts in order
./02-create-tfstate-bucket.sh    # S3 bucket + DynamoDB table for Terraform state
./03-create-oidc-provider.sh     # GitHub OIDC identity provider  
./04-create-github-role.sh       # IAM role GitHub Actions will assume
./05-set-github-secrets.sh       # Push AWS config to GitHub repository secrets
```

**What each script creates:**

| Script | Creates | Purpose |
|--------|---------|---------|
| `02-create-tfstate-bucket.sh` | S3 bucket `trackstack-tfstate-939091506005` + DynamoDB table `trackstack-terraform-locks` | Stores Terraform state with locking to prevent concurrent modifications |
| `03-create-oidc-provider.sh` | OIDC provider `token.actions.githubusercontent.com` | Establishes trust relationship with GitHub |
| `04-create-github-role.sh` | IAM role `trackstack-github-deploy-role` | Role GitHub Actions assumes via OIDC |
| `05-set-github-secrets.sh` | GitHub repository secrets | Stores role ARN, region, and bucket names for workflows |

---

### Step 6: Sync IAM Bootstrap (Terraform)

The bootstrap scripts create the OIDC trust and the IAM role, but **permissions are managed by Terraform** so you can evolve them without re-running bootstrap scripts.

```bash
cd iac/bootstrap
./06-sync-iam-bootstrap.sh
```

For non-interactive runs:

```bash
cd iac/bootstrap
./06-sync-iam-bootstrap.sh --auto-approve
```

### Step 7: Destroy IAM Bootstrap (Terraform)

If you need to tear down the IAM policy managed by Terraform:

```bash
cd iac/bootstrap
./07-destroy-iam-bootstrap.sh
```

---

## Step 0.5: Astro SSR (Terraform)

After IAM bootstrap is synced, provision the Astro SSR stack:

```bash
cd iac/environments/step-0-5-astro-ssr
terraform init \
  -backend-config="bucket=${TFSTATE_BUCKET}" \
  -backend-config="dynamodb_table=${TFSTATE_LOCK_TABLE}" \
  -backend-config="region=${AWS_REGION}"

terraform apply \
  -var="aws_region=${AWS_REGION}"
```

Before first deploy, set runtime secrets in SSM:

```bash
cd iac/environments/step-0-5-astro-ssr
./01-set-runtime-ssm.sh
```

---

## Resources Created

| Resource | Name | Purpose |
|----------|------|---------|
| IAM User | `trackstack-admin` | Your local CLI access (human user) |
| S3 Bucket | `trackstack-tfstate-939091506005` | Terraform state storage |
| DynamoDB Table | `trackstack-terraform-locks` | Prevents concurrent Terraform runs |
| OIDC Provider | `token.actions.githubusercontent.com` | Trusts GitHub as identity provider |
| IAM Role | `trackstack-github-deploy-role` | Role GitHub Actions assumes |
| IAM Policy | `trackstack-terraform-deploy-policy` | Permissions for Terraform deploys (managed by Terraform) |

## GitHub Secrets Set

| Secret | Value | Purpose |
|--------|-------|---------|
| `AWS_ROLE_ARN` | `arn:aws:iam::939091506005:role/trackstack-github-deploy-role` | Role for GitHub Actions to assume |
| `AWS_REGION` | `eu-west-1` | AWS region for deployments |
| `TFSTATE_BUCKET` | `trackstack-tfstate-939091506005` | S3 bucket for Terraform state |
| `TFSTATE_LOCK_TABLE` | `trackstack-terraform-locks` | DynamoDB table for state locking |

**Note:** These are NOT AWS credentials. They are configuration values that tell GitHub Actions which AWS role to assume via OIDC.

---

## How GitHub Actions Uses This

Once bootstrap is complete, your GitHub Actions workflows can authenticate to AWS without any stored credentials:

```yaml
# .github/workflows/deploy.yml example
name: Deploy

on:
  push:
    branches: [main]

permissions:
  id-token: write   # Required for OIDC token generation
  contents: read    # Required for checkout

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Configure AWS credentials via OIDC
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
          aws-region: ${{ secrets.AWS_REGION }}
      
      # Now AWS CLI and Terraform commands work with temporary credentials
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
      
      - name: Terraform Init
        run: terraform init -backend-config="bucket=${{ secrets.TFSTATE_BUCKET }}"
```

**The flow:**
1. GitHub Actions generates a JWT token (via `id-token: write` permission)
2. `configure-aws-credentials` action exchanges the JWT for temporary AWS credentials
3. Terraform/ AWS CLI commands run with those credentials
4. Credentials expire when the workflow ends

---

## Cleanup (Rerun Bootstrap)

If you need to rerun bootstrap from scratch, use the cleanup script:

```bash
cd iac/bootstrap
./90-destroy-iam-bootstrap.sh
./91-cleanup-bootstrap.sh --force-empty-buckets
```

This deletes the admin IAM user, Terraform state bucket, DynamoDB lock table, OIDC provider, Terraform IAM policy, and GitHub secrets. Use with care.
