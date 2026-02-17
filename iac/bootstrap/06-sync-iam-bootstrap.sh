#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

IAM_BOOTSTRAP_DIR="$SCRIPT_DIR/../environments/iam-bootstrap"
AUTO_APPROVE=false

while [ $# -gt 0 ]; do
    case "$1" in
        --auto-approve)
            AUTO_APPROVE=true
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: ./06-sync-iam-bootstrap.sh [--auto-approve]"
            exit 1
            ;;
    esac
    shift
done

if ! command -v terraform >/dev/null 2>&1; then
    echo "ERROR: terraform is not installed."
    echo "Install it and try again."
    exit 1
fi

if [ ! -d "$IAM_BOOTSTRAP_DIR" ]; then
    echo "ERROR: iam-bootstrap environment not found at: $IAM_BOOTSTRAP_DIR"
    exit 1
fi

echo "=========================================="
echo " STEP 6: Sync IAM Bootstrap (Terraform)"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Region: $AWS_REGION"
echo "Role: $GITHUB_ROLE_NAME"
echo "Policy: $TERRAFORM_POLICY_NAME"
echo "State Bucket: $TFSTATE_BUCKET"
echo "Lock Table: $TFSTATE_LOCK_TABLE"
echo ""

terraform -chdir="$IAM_BOOTSTRAP_DIR" init \
  -backend-config="bucket=${TFSTATE_BUCKET}" \
  -backend-config="dynamodb_table=${TFSTATE_LOCK_TABLE}" \
  -backend-config="region=${AWS_REGION}"

APPLY_FLAGS=()
if [ "$AUTO_APPROVE" = true ]; then
    APPLY_FLAGS+=("-auto-approve")
fi

terraform -chdir="$IAM_BOOTSTRAP_DIR" apply \
  "${APPLY_FLAGS[@]}" \
  -var="aws_region=${AWS_REGION}" \
  -var="role_name=${GITHUB_ROLE_NAME}" \
  -var="policy_name=${TERRAFORM_POLICY_NAME}" \
  -var="tfstate_bucket=${TFSTATE_BUCKET}" \
  -var="tfstate_lock_table=${TFSTATE_LOCK_TABLE}"

echo ""
echo "=========================================="
echo " IAM BOOTSTRAP COMPLETE"
echo "=========================================="
echo ""
