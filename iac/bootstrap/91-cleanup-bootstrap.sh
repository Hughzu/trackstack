#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

FORCE_EMPTY_BUCKETS=false
YES=false

while [ $# -gt 0 ]; do
    case "$1" in
        --force-empty-buckets)
            FORCE_EMPTY_BUCKETS=true
            ;;
        --yes)
            YES=true
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: ./99-cleanup-bootstrap.sh [--force-empty-buckets] [--yes]"
            exit 1
            ;;
    esac
    shift
done

echo "=========================================="
echo " BOOTSTRAP CLEANUP"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Region: $AWS_REGION"
echo "Account (config): $AWS_ACCOUNT_ID"
echo ""
echo "Resources to delete:"
echo "- IAM User: $ADMIN_USER_NAME"
echo "- IAM Role: $GITHUB_ROLE_NAME"
echo "- OIDC Provider: $OIDC_PROVIDER_URL"
echo "- S3 Bucket: $TFSTATE_BUCKET"
echo "- DynamoDB Table: $TFSTATE_LOCK_TABLE"
echo "- GitHub Secrets in: $GITHUB_REPO"
echo ""

ACCOUNT_ID=$(aws sts get-caller-identity \
    --query Account \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

if [ -z "$ACCOUNT_ID" ]; then
    echo "ERROR: Unable to resolve AWS account identity."
    exit 1
fi

if [ "$ACCOUNT_ID" != "$AWS_ACCOUNT_ID" ]; then
    echo "WARNING: Config account ($AWS_ACCOUNT_ID) does not match AWS profile account ($ACCOUNT_ID)."
fi

if [ "$YES" = false ]; then
    echo "This will permanently delete the resources above."
    read -r -p "Type DELETE to continue: " CONFIRM
    if [ "$CONFIRM" != "DELETE" ]; then
        echo "Aborted."
        exit 0
    fi
fi

echo ""
echo "Deleting GitHub secrets..."
if gh auth status &>/dev/null; then
    gh secret delete AWS_ROLE_ARN --repo "$GITHUB_REPO" 2>/dev/null || true
    gh secret delete AWS_REGION --repo "$GITHUB_REPO" 2>/dev/null || true
    gh secret delete TFSTATE_BUCKET --repo "$GITHUB_REPO" 2>/dev/null || true
    gh secret delete TFSTATE_LOCK_TABLE --repo "$GITHUB_REPO" 2>/dev/null || true
else
    echo "GitHub CLI not authenticated. Skipping secret deletion."
fi

echo ""
echo "Deleting IAM role policies and role..."
ATTACHED_POLICIES=$(aws iam list-attached-role-policies \
    --role-name "$GITHUB_ROLE_NAME" \
    --query 'AttachedPolicies[].PolicyArn' \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

for POLICY_ARN in $ATTACHED_POLICIES; do
    aws iam detach-role-policy \
        --role-name "$GITHUB_ROLE_NAME" \
        --policy-arn "$POLICY_ARN" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
done

INLINE_POLICIES=$(aws iam list-role-policies \
    --role-name "$GITHUB_ROLE_NAME" \
    --query 'PolicyNames[]' \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

for POLICY_NAME in $INLINE_POLICIES; do
    aws iam delete-role-policy \
        --role-name "$GITHUB_ROLE_NAME" \
        --policy-name "$POLICY_NAME" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
done

aws iam delete-role \
    --role-name "$GITHUB_ROLE_NAME" \
    --profile "$AWS_PROFILE" 2>/dev/null || true

echo ""
echo "Deleting Terraform IAM policy..."
TERRAFORM_POLICY_ARN=$(aws iam list-policies \
    --scope Local \
    --query "Policies[?PolicyName=='${TERRAFORM_POLICY_NAME}'].Arn" \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

if [ -n "$TERRAFORM_POLICY_ARN" ]; then
    POLICY_VERSIONS=$(aws iam list-policy-versions \
        --policy-arn "$TERRAFORM_POLICY_ARN" \
        --query 'Versions[?IsDefaultVersion==`false`].VersionId' \
        --output text \
        --profile "$AWS_PROFILE" 2>/dev/null || true)

    for VERSION_ID in $POLICY_VERSIONS; do
        aws iam delete-policy-version \
            --policy-arn "$TERRAFORM_POLICY_ARN" \
            --version-id "$VERSION_ID" \
            --profile "$AWS_PROFILE" 2>/dev/null || true
    done

    aws iam delete-policy \
        --policy-arn "$TERRAFORM_POLICY_ARN" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
fi

echo ""
echo "Deleting OIDC provider..."
OIDC_PROVIDER_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/${OIDC_PROVIDER_URL}"
aws iam delete-open-id-connect-provider \
    --open-id-connect-provider-arn "$OIDC_PROVIDER_ARN" \
    --profile "$AWS_PROFILE" 2>/dev/null || true

echo ""
echo "Deleting Terraform state bucket and lock table..."
if aws s3api head-bucket --bucket "$TFSTATE_BUCKET" --profile "$AWS_PROFILE" 2>/dev/null; then
    if [ "$FORCE_EMPTY_BUCKETS" = true ]; then
        aws s3 rm "s3://$TFSTATE_BUCKET" --recursive --profile "$AWS_PROFILE" 2>/dev/null || true
    fi

    aws s3api delete-bucket \
        --bucket "$TFSTATE_BUCKET" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
else
    echo "Bucket not found: $TFSTATE_BUCKET"
fi

aws dynamodb delete-table \
    --table-name "$TFSTATE_LOCK_TABLE" \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" 2>/dev/null || true

echo ""
echo "Deleting admin user..."
ACCESS_KEYS=$(aws iam list-access-keys \
    --user-name "$ADMIN_USER_NAME" \
    --query 'AccessKeyMetadata[].AccessKeyId' \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

for KEY_ID in $ACCESS_KEYS; do
    aws iam delete-access-key \
        --user-name "$ADMIN_USER_NAME" \
        --access-key-id "$KEY_ID" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
done

USER_POLICIES=$(aws iam list-attached-user-policies \
    --user-name "$ADMIN_USER_NAME" \
    --query 'AttachedPolicies[].PolicyArn' \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

for POLICY_ARN in $USER_POLICIES; do
    aws iam detach-user-policy \
        --user-name "$ADMIN_USER_NAME" \
        --policy-arn "$POLICY_ARN" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
done

INLINE_USER_POLICIES=$(aws iam list-user-policies \
    --user-name "$ADMIN_USER_NAME" \
    --query 'PolicyNames[]' \
    --output text \
    --profile "$AWS_PROFILE" 2>/dev/null || true)

for POLICY_NAME in $INLINE_USER_POLICIES; do
    aws iam delete-user-policy \
        --user-name "$ADMIN_USER_NAME" \
        --policy-name "$POLICY_NAME" \
        --profile "$AWS_PROFILE" 2>/dev/null || true
done

aws iam delete-login-profile \
    --user-name "$ADMIN_USER_NAME" \
    --profile "$AWS_PROFILE" 2>/dev/null || true

aws iam delete-user \
    --user-name "$ADMIN_USER_NAME" \
    --profile "$AWS_PROFILE" 2>/dev/null || true

echo ""
echo "=========================================="
echo " CLEANUP COMPLETE"
echo "=========================================="
echo ""
