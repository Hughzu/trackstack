#!/bin/bash
set -e
#AWS_PROFILE="trackstack"
AWS_ACCOUNT_ID="939091506005"

# 1. Delete IAM roles
echo "Deleting IAM roles..."
aws iam delete-role-policy --role-name trackstack-github-blog-role --policy-name BlogDeploymentPolicy || true
aws iam delete-role --role-name trackstack-github-blog-role || true
aws iam detach-role-policy --role-name trackstack-github-terraform-role --policy-arn arn:aws:iam::aws:policy/AdministratorAccess  || true
aws iam delete-role --role-name trackstack-github-terraform-role  || true

# 2. Delete OIDC provider
echo "Deleting OIDC provider..."
aws iam delete-open-id-connect-provider \
  --open-id-connect-provider-arn "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/token.actions.githubusercontent.com" \
   || true

# 3. Empty and delete tfstate bucket
echo "Emptying S3 tfstate bucket..."
aws s3 rm s3://trackstack-terraform-state --recursive  || true
echo "Deleting S3 tfstate bucket..."
aws s3api delete-bucket --bucket trackstack-terraform-state  || true
echo "=== Cleanup complete ==="