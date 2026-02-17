#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

echo "=========================================="
echo " STEP 4: Create GitHub Deployment Role"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Role: $GITHUB_ROLE_NAME"
echo "Repository: $GITHUB_REPO"
echo ""

OIDC_PROVIDER_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/${OIDC_PROVIDER_URL}"
ROLE_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:role/${GITHUB_ROLE_NAME}"

echo "Creating trust policy..."
cat > /tmp/github-trust-policy.json <<EOF
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

echo "Checking if role already exists..."
EXISTING_ROLE=$(aws iam list-roles \
    --query "Roles[?RoleName=='${GITHUB_ROLE_NAME}'].RoleName" \
    --output text \
    --profile "$AWS_PROFILE")

if [ -n "$EXISTING_ROLE" ]; then
    echo "Role already exists. Updating trust policy..."
else
    echo "Creating IAM role..."
    aws iam create-role \
        --role-name "$GITHUB_ROLE_NAME" \
        --assume-role-policy-document file:///tmp/github-trust-policy.json \
        --description "Role for GitHub Actions to deploy TrackStack infrastructure" \
        --tags Key=Project,Value=trackstack Key=ManagedBy,Value=bootstrap Key=Purpose,Value=github-deployment \
        --profile "$AWS_PROFILE"
fi

aws iam update-assume-role-policy \
    --role-name "$GITHUB_ROLE_NAME" \
    --policy-document file:///tmp/github-trust-policy.json \
    --profile "$AWS_PROFILE"

rm -f /tmp/github-trust-policy.json

echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "Role Name: $GITHUB_ROLE_NAME"
echo "Role ARN: $ROLE_ARN"
echo ""
echo "Next: ./05-set-github-secrets.sh"
echo ""
