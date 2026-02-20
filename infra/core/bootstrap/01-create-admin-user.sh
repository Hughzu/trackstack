#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Check BEFORE sourcing config.sh (which sets AWS_PROFILE)
if [ -n "$AWS_PROFILE" ]; then
    echo "ERROR: AWS_PROFILE is set to '$AWS_PROFILE'. Please unset it:"
    echo "  unset AWS_PROFILE"
    echo "  Then run this script again."
    exit 1
fi

source "$SCRIPT_DIR/config.sh"

# Save profile name before unsetting (this script uses root credentials)
PROFILE_NAME="$AWS_PROFILE"
unset AWS_PROFILE

echo "=========================================="
echo " STEP 1: Create Admin IAM User"
echo "=========================================="
echo ""
echo "WARNING: This script must be run with root AWS credentials."
echo "         Do NOT set AWS_PROFILE for this script."
echo ""

echo "Creating IAM user: $ADMIN_USER_NAME..."

aws iam create-user \
    --user-name "$ADMIN_USER_NAME" \
    --tags Key=Project,Value=trackstack Key=ManagedBy,Value=bootstrap Key=Purpose,Value=admin-access \
    2>/dev/null || echo "User already exists, continuing..."

echo "Attaching AdministratorAccess policy..."
aws iam attach-user-policy \
    --user-name "$ADMIN_USER_NAME" \
    --policy-arn arn:aws:iam::aws:policy/AdministratorAccess

echo "Creating access keys..."
ACCESS_KEY_OUTPUT=$(aws iam create-access-key --user-name "$ADMIN_USER_NAME")
ACCESS_KEY_ID=$(echo "$ACCESS_KEY_OUTPUT" | jq -r '.AccessKey.AccessKeyId')
SECRET_ACCESS_KEY=$(echo "$ACCESS_KEY_OUTPUT" | jq -r '.AccessKey.SecretAccessKey')

echo ""
echo "=========================================="
echo " SAVE THESE CREDENTIALS SECURELY!"
echo "=========================================="
echo ""
echo "Access Key ID: $ACCESS_KEY_ID"
echo "Secret Access Key: $SECRET_ACCESS_KEY"
echo ""
echo "=========================================="
echo " NEXT: Run this command to save credentials:"
echo "=========================================="
echo ""
echo "  aws configure --profile trackstack"
echo ""
echo "When prompted:"
echo "  - Access Key ID: $ACCESS_KEY_ID"
echo "  - Secret Access Key: $SECRET_ACCESS_KEY"
echo "  - Region: $AWS_REGION"
echo "  - Output: json"
echo ""
echo "Then continue with:"
echo ""
echo "  export AWS_PROFILE=trackstack"
echo "  aws sts get-caller-identity"
echo "  ./02-create-tfstate-bucket.sh"
echo ""
