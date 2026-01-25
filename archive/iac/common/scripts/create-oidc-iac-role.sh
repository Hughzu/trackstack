#!/bin/bash

AWS_PROFILE="trackstack"
AWS_ACCOUNT_ID="939091506005"
GITHUB_REPO="Hughzu/trackstack"
OIDC_PROVIDER_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/token.actions.githubusercontent.com"

echo "Creating Terraform Deployment Role with Admin Access..."
echo "Repository: ${GITHUB_REPO}"
echo "AWS Account: ${AWS_ACCOUNT_ID}"
echo ""

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

echo "Creating IAM role..."
aws iam create-role \
  --role-name trackstack-github-terraform-role \
  --assume-role-policy-document file:///tmp/terraform-trust-policy.json \
  --description "Role for GitHub Actions to deploy TrackStack infrastructure via Terraform" \
  --tags Key=Project,Value=trackstack Key=ManagedBy,Value=script Key=Purpose,Value=terraform-deployment \
  --profile $AWS_PROFILE

echo "Attaching AdministratorAccess policy..."
aws iam attach-role-policy \
  --role-name trackstack-github-terraform-role \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess \
  --profile $AWS_PROFILE

rm -f /tmp/terraform-trust-policy.json

echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "Terraform Deployment Role created with FULL ADMIN ACCESS:"
echo "  Name: trackstack-github-terraform-role"
echo "  ARN: arn:aws:iam::${AWS_ACCOUNT_ID}:role/trackstack-github-terraform-role"
echo "  Policy: AdministratorAccess (full AWS access)"
echo ""
echo "    Security note: This role has full AWS access."
echo "    Only use in trusted personal projects."
echo "    The role is restricted to your GitHub repo via OIDC."
echo ""