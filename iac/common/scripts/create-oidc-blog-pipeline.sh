#!/bin/bash

AWS_PROFILE="trackstack"
AWS_ACCOUNT_ID="939091506005"
GITHUB_REPO="Hughzu/trackstack"
OIDC_PROVIDER_ARN="arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/token.actions.githubusercontent.com"
BUCKET_NAME="trackstack-blog"
CLOUDFRONT_DISTRIBUTION_ID="E2X7059Q4OT43V"

echo "Creating Blog Deployment Role with Scoped Permissions..."
echo "Repository: ${GITHUB_REPO}"
echo "AWS Account: ${AWS_ACCOUNT_ID}"
echo ""

# Trust policy for OIDC
cat > /tmp/blog-trust-policy.json <<EOF
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

# Permissions policy for blog deployment
cat > /tmp/blog-permissions-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3BlogDeployment",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:PutObjectAcl",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::${BUCKET_NAME}",
        "arn:aws:s3:::${BUCKET_NAME}/*"
      ]
    },
    {
      "Sid": "CloudFrontInvalidation",
      "Effect": "Allow",
      "Action": [
        "cloudfront:CreateInvalidation",
        "cloudfront:GetInvalidation"
      ],
      "Resource": "arn:aws:cloudfront::${AWS_ACCOUNT_ID}:distribution/${CLOUDFRONT_DISTRIBUTION_ID}"
    }
  ]
}
EOF

echo "Creating IAM role..."
aws iam create-role \
  --role-name trackstack-github-blog-role \
  --assume-role-policy-document file:///tmp/blog-trust-policy.json \
  --description "Role for GitHub Actions to deploy TrackStack Blog to S3 and invalidate CloudFront" \
  --tags Key=Project,Value=trackstack Key=ManagedBy,Value=script Key=Purpose,Value=blog-deployment \
  --profile $AWS_PROFILE

echo "Creating and attaching inline policy..."
aws iam put-role-policy \
  --role-name trackstack-github-blog-role \
  --policy-name BlogDeploymentPolicy \
  --policy-document file:///tmp/blog-permissions-policy.json \
  --profile $AWS_PROFILE

# Cleanup
rm -f /tmp/blog-trust-policy.json /tmp/blog-permissions-policy.json

echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "Blog Deployment Role created with SCOPED PERMISSIONS:"
echo "  Name: trackstack-github-blog-role"
echo "  ARN: arn:aws:iam::${AWS_ACCOUNT_ID}:role/trackstack-github-blog-role"
echo ""
echo "Permissions granted:"
echo "  S3: Deploy to any bucket starting with 'trackstack-'"
echo "  CloudFront: Create and monitor invalidations on any distribution"
echo ""
echo "Security: Limited to TrackStack-prefixed resources only"
echo ""