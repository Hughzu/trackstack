#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

echo "=========================================="
echo " STEP 3: Create OIDC Provider"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Provider: $OIDC_PROVIDER_URL"
echo ""

OIDC_PROVIDER_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/${OIDC_PROVIDER_URL}"

echo "Checking if OIDC provider already exists..."
EXISTING_PROVIDER=$(aws iam list-open-id-connect-providers \
    --query "OpenIDConnectProviderList[?contains(Arn, '${OIDC_PROVIDER_URL}')].Arn" \
    --output text \
    --profile "$AWS_PROFILE")

if [ -n "$EXISTING_PROVIDER" ]; then
    echo "OIDC provider already exists: $EXISTING_PROVIDER"
    echo "Skipping creation."
else
    echo "Creating OIDC provider..."
    OIDC_ARN=$(aws iam create-open-id-connect-provider \
        --url "https://${OIDC_PROVIDER_URL}" \
        --client-id-list sts.amazonaws.com \
        --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1 \
        --tags Key=Project,Value=trackstack Key=ManagedBy,Value=bootstrap \
        --profile "$AWS_PROFILE" \
        --query "OpenIDConnectProviderArn" \
        --output text)
    
    echo "Created: $OIDC_ARN"
fi

echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "OIDC Provider ARN: $OIDC_PROVIDER_ARN"
echo ""
echo "Next: ./04-create-github-role.sh"
echo ""
