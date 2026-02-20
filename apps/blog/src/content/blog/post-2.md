---
title: '2. Deploying with OIDC to AWS'
description: 'OIDC implementation'
pubDate: '2025-10-26'
heroImage: '../../assets/hero-gradient.svg'
---

## Security First: Why OIDC Matters

Before deploying a single line of infrastructure code, I needed to solve the authentication problem. How does GitHub Actions deploy to AWS without storing long-lived credentials?

**The old way:** Store AWS access keys as GitHub secrets. Rotate them manually. Hope nobody steals them. Pray you remember to revoke them when needed.

**The modern way:** OpenID Connect (OIDC). No secrets. Temporary credentials. Cryptographically verified trust.

If you're deploying infrastructure or application code from CI/CD pipelines and still using long-lived credentials, there's a better way. OIDC provides stronger security with less operational overhead. Let me show you how it works! 

---

## The Problem with Long-Lived Credentials

Here's what happens when you use AWS access keys in GitHub Actions:

**Security risks:**
- Keys exist permanently in GitHub secrets (even if encrypted, they're still there)
- If GitHub is compromised, attackers get persistent AWS access
- Keys don't expire automatically - you must remember to rotate them
- If someone leaves your team, you need to manually revoke and regenerate
- No audit trail showing *which specific GitHub workflow run* used the credentials

**Operational overhead:**
- Manual key rotation processes
- Multiple keys to manage (dev, staging, prod)
- Secret sprawl across repositories
- Keys in CI/CD, keys in local env files, keys in documentation (yikes)

For a personal project like TrackStack, this might seem acceptable. But if I wouldn't use long-lived credentials at a company with millions in revenue at stake, why would I use them here?

---

## How OIDC Actually Works

OIDC eliminates stored credentials through a trust relationship and temporary tokens.

**The flow:**

1. **GitHub Action starts:** "Hey AWS, I'm workflow XYZ from repository Hughzu/trackstack"
2. **GitHub generates a token:** A cryptographically signed JWT containing:
   - Repository name
   - Workflow details
   - Branch/tag information
   - Actor (who triggered it)
   - Expiration time (usually 1 hour)
3. **GitHub Action presents token to AWS:** "Here's my ID token from GitHub"
4. **AWS validates the token:**
   - Checks GitHub's signature (using the OIDC provider's public keys)
   - Verifies token isn't expired
   - Confirms claims match the trust policy (correct repo, branch, etc.)
5. **AWS issues temporary credentials:** Valid for the session only (typically 1 hour)
6. **Workflow runs with temporary access:** Deploys infrastructure
7. **Credentials expire automatically:** No cleanup needed

There are no credentials to steal. Even if someone intercepts the JWT token, it's only valid for minutes and can't be reused after expiration.

---

## Setting Up OIDC for TrackStack

I created three setup scripts to configure OIDC authentication. Let's walk through each one.

### Step 1: Create the OIDC Provider

First, AWS needs to trust GitHub as an identity provider.

```bash
#!/bin/bash

AWS_PROFILE="trackstack"

aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1 \
  --profile $AWS_PROFILE
```

- **url:** GitHub's OIDC token endpoint (where AWS fetches GitHub's public keys)
- **client-id-list:** We're using AWS STS as the audience
- **thumbprint-list:** SHA-1 fingerprint of GitHub's SSL certificate (validates we're talking to the real GitHub)

Run this once per AWS account. It creates a trust relationship that says "AWS, you can verify tokens from GitHub."

### Step 2: Create the Terraform State Bucket

Terraform needs somewhere to store its state file. This is standard practice - you want centralized state, not local files that can get out of sync.

```bash
#!/bin/bash

BUCKET_NAME="trackstack-terraform-state"
REGION="eu-central-1"
AWS_PROFILE="trackstack"

echo "Creating S3 bucket..."
aws s3api create-bucket \
  --bucket $BUCKET_NAME \
  --region $REGION \
  --create-bucket-configuration LocationConstraint=$REGION \
  --profile $AWS_PROFILE

echo "Enabling versioning..."
aws s3api put-bucket-versioning \
  --bucket $BUCKET_NAME \
  --versioning-configuration Status=Enabled \
  --profile $AWS_PROFILE

echo "Enabling encryption..."
aws s3api put-bucket-encryption \
  --bucket $BUCKET_NAME \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }' \
  --profile $AWS_PROFILE

echo "Blocking public access..."
aws s3api put-public-access-block \
  --bucket $BUCKET_NAME \
  --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true \
  --profile $AWS_PROFILE
```

- **Versioning enabled:** If Terraform state gets corrupted, I can roll back
- **Encryption at rest:** State files contain sensitive data (resource IDs, configurations)
- **Public access blocked:** No one should ever access this bucket except Terraform
- **Region-specific:** eu-central-1 because I'm in Europe and GDPR matters

This bucket is the source of truth for my infrastructure. Treat it like a database.

### Step 3: Create the IAM Role for GitHub Actions

Now we create a role that GitHub Actions can assume via OIDC.

```bash
#!/bin/bash

AWS_PROFILE="trackstack"
AWS_ACCOUNT_ID="***"
GITHUB_REPO="Hughzu/trackstack"
OIDC_PROVIDER_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/token.actions.githubusercontent.com"

echo "Creating Terraform Deployment Role with Admin Access..."

cat > /tmp/terraform-trust-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "${OIDC_PROVIDER_ARN}"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:${GITHUB_REPO}:*"
        }
      }
    }
  ]
}
EOF

aws iam create-role \
  --role-name trackstack-github-terraform-role \
  --assume-role-policy-document file:///tmp/terraform-trust-policy.json \
  --description "Role for GitHub Actions to deploy TrackStack infrastructure via Terraform" \
  --tags Key=Project,Value=trackstack Key=ManagedBy,Value=script Key=Purpose,Value=terraform-deployment \
  --profile $AWS_PROFILE

aws iam attach-role-policy \
  --role-name trackstack-github-terraform-role \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess \
  --profile $AWS_PROFILE

rm -f /tmp/terraform-trust-policy.json
```

**The trust policy:**

The conditions are the security boundary:

```json
"Condition": {
  "StringEquals": {
    "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
  },
  "StringLike": {
    "token.actions.githubusercontent.com:sub": "repo:${GITHUB_REPO}:*"
  }
}
```

- **aud check:** Token must be intended for AWS STS (prevents token reuse for other services)
- **sub check:** Token must come from `repo:Hughzu/trackstack:*` (the wildcard allows any branch/tag)

Even if someone steals a token, they can't use it unless they're running from my exact GitHub repository.

**About AdministratorAccess:**

Yes, I'm using full admin access. I do not know at this point all the AWS resources I will use, I need to change that in the future once the dust settles. In a production environment with multiple team members, you'd use least-privilege policies scoped to specific resources.

The key security control is the OIDC trust policy - only my GitHub repository can assume this role.

---

## The GitHub Actions Workflow

Here's how it all comes together in `.github/workflows/deploy-blog-iac.yml`:

```yaml
name: Deploy Infrastructure

on:
  push:
    branches: [main]
    paths:
      - 'iac/blog/**'
      - '.github/workflows/deploy-blog-iac.yml'
  workflow_dispatch:
    inputs:
      action:
        description: 'Terraform action to perform'
        required: true
        default: 'plan'
        type: choice
        options:
          - plan
          - apply
          - destroy

permissions:
  id-token: write   # Required for OIDC
  contents: read    # Required for checkout

jobs:
  terraform:
    runs-on: ubuntu-latest
    
    defaults:
      run:
        working-directory: iac/blog
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: ~1.0
          terraform_wrapper: false

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
          role-session-name: GitHubActions-Terraform-${{ github.run_number }}
          aws-region: ${{ secrets.AWS_REGION }}

      - name: Terraform Format Check
        id: fmt
        run: terraform fmt -check -recursive
        continue-on-error: true

      - name: Terraform Init
        id: init
        run: |
          terraform init
          echo "✅ Terraform initialized successfully"

      - name: Terraform Validate
        id: validate
        run: |
          terraform validate
          echo "✅ Terraform configuration is valid"

      - name: Terraform Plan
        id: plan
        run: |
          echo "🔍 Planning infrastructure changes..."
          set +e
          terraform plan -detailed-exitcode -no-color -out=tfplan
          PLAN_EXIT_CODE=$?
          set -e
          
          if [ $PLAN_EXIT_CODE -eq 0 ]; then
            echo "✅ No changes needed"
            echo "has_changes=false" >> $GITHUB_OUTPUT
          elif [ $PLAN_EXIT_CODE -eq 2 ]; then
            echo "📋 Changes detected, ready to apply"
            echo "has_changes=true" >> $GITHUB_OUTPUT
          else
            echo "❌ Plan failed"
            exit 1
          fi

      - name: Terraform Apply (Auto on Push to Main)
        if: github.ref == 'refs/heads/main' && github.event_name == 'push' && steps.plan.outputs.has_changes == 'true'
        run: |
          echo "🚀 Auto-applying infrastructure changes..."
          terraform apply tfplan
          echo "✅ Infrastructure deployed successfully!"

      - name: Manual Apply
        if: github.event_name == 'workflow_dispatch' && github.event.inputs.action == 'apply'
        run: |
          echo "🚀 Manually applying Terraform changes..."
          terraform apply tfplan

      - name: Terraform Destroy
        if: github.event_name == 'workflow_dispatch' && github.event.inputs.action == 'destroy'
        run: |
          echo "🗑️ Destroying infrastructure..."
          terraform destroy -auto-approve
```

**Key elements:**

**Permissions block:**
```yaml
permissions:
  id-token: write   # Allows GitHub to generate OIDC tokens
  contents: read    # Allows checkout of repository code
```

Without `id-token: write`, the workflow can't generate the JWT token needed for OIDC authentication.

**AWS credential configuration:**
```yaml
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v4
  with:
    role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
    role-session-name: GitHubActions-Terraform-${{ github.run_number }}
    aws-region: ${{ secrets.AWS_REGION }}
```

This step does all the OIDC magic:
1. Requests an OIDC token from GitHub
2. Exchanges it with AWS STS for temporary credentials
3. Sets AWS environment variables for subsequent steps
4. Credentials expire after the workflow completes

**Terraform plan with exit code handling:**

Terraform's exit codes matter:
- `0` = No changes needed
- `2` = Changes detected, ready to apply
- `1` = Error occurred

I capture the exit code to conditionally apply changes only when needed.

**Automatic apply on main branch:**

When code is pushed to main and changes are detected, Terraform automatically applies. For a personal project where I'm the only contributor, this speeds up deployments. In a team environment, you'd want pull request reviews and manual approval gates.

---

## GitHub Repository Secrets

Two secrets need to be configured in your GitHub repository settings:

**AWS_ROLE_ARN:**
```
arn:aws:iam::***:role/trackstack-github-terraform-role
```

**AWS_REGION:**
```
eu-central-1
```

That's it. No access keys. No secret access keys. Just the role ARN and region.

---

## Why This Approach Wins

Let's compare the security posture:

### Traditional Access Keys:

❌ Permanent credentials stored in GitHub  
❌ Manual rotation required  
❌ Credentials work from anywhere if leaked  
❌ No automatic expiration  
❌ Difficult to audit which workflow run used them  
❌ Revocation requires regenerating and updating secrets  

### OIDC with Temporary Credentials:

✅ Zero stored credentials  
✅ Automatic expiration (typically 1 hour)  
✅ Credentials only work from specified GitHub repository  
✅ Cryptographic verification via JWT signatures  
✅ Full audit trail in CloudTrail showing exact workflow runs  
✅ Instant revocation by modifying IAM role trust policy  
✅ Principle of least privilege by default (each workflow gets fresh credentials)  

---

## Cost Impact: Zero

OIDC is free. IAM is free. The only costs are for resources Terraform creates (S3, CloudFront, etc.).

The S3 bucket for Terraform state costs is really cheep (approximately €0.02 per month) for storage (assuming <1GB state file) plus negligible request costs.