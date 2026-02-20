#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

echo "=========================================="
echo " STEP 5: Set GitHub Secrets"
echo "=========================================="
echo ""
echo "Repository: $GITHUB_REPO"
echo ""

ROLE_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:role/${GITHUB_ROLE_NAME}"

echo "Checking GitHub CLI authentication..."
if ! gh auth status &>/dev/null; then
    echo "ERROR: GitHub CLI is not authenticated."
    echo "Please run: gh auth login"
    exit 1
fi

echo "Setting repository secrets..."
echo ""

echo "Setting AWS_ROLE_ARN..."
gh secret set AWS_ROLE_ARN \
    --repo "$GITHUB_REPO" \
    --body "$ROLE_ARN"

echo "Setting AWS_REGION..."
gh secret set AWS_REGION \
    --repo "$GITHUB_REPO" \
    --body "$AWS_REGION"

echo "Setting TFSTATE_BUCKET..."
gh secret set TFSTATE_BUCKET \
    --repo "$GITHUB_REPO" \
    --body "$TFSTATE_BUCKET"

echo "Setting TFSTATE_LOCK_TABLE..."
gh secret set TFSTATE_LOCK_TABLE \
    --repo "$GITHUB_REPO" \
    --body "$TFSTATE_LOCK_TABLE"

echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "Secrets configured:"
echo ""
gh secret list --repo "$GITHUB_REPO"
echo ""
echo "=========================================="
echo " BOOTSTRAP COMPLETE!"
echo "=========================================="
echo ""
echo "Your GitHub Actions can now authenticate to AWS via OIDC."
echo ""
echo "Example workflow usage:"
echo ""
cat << 'WORKFLOW'
name: Deploy

on:
  push:
    branches: [main]

permissions:
  id-token: write
  contents: read

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: \${{ secrets.AWS_ROLE_ARN }}
          aws-region: \${{ secrets.AWS_REGION }}
      
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v3
      
      - name: Terraform Init
        run: terraform init
          -backend-config="bucket=\${{ secrets.TFSTATE_BUCKET }}"
WORKFLOW
echo ""
