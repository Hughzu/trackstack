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

echo "Creating permissions policy..."
cat > /tmp/github-permissions-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "TerraformState",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::${TFSTATE_BUCKET}",
        "arn:aws:s3:::${TFSTATE_BUCKET}/*"
      ]
    },
    {
      "Sid": "DynamoDBLocks",
      "Effect": "Allow",
      "Action": [
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:DeleteItem"
      ],
      "Resource": "arn:aws:dynamodb:${AWS_REGION}:${AWS_ACCOUNT_ID}:table/${TFSTATE_LOCK_TABLE}"
    },
    {
      "Sid": "Lambda",
      "Effect": "Allow",
      "Action": [
        "lambda:CreateFunction",
        "lambda:UpdateFunctionCode",
        "lambda:UpdateFunctionConfiguration",
        "lambda:GetFunction",
        "lambda:DeleteFunction",
        "lambda:AddPermission",
        "lambda:RemovePermission",
        "lambda:GetFunctionConfiguration",
        "lambda:ListTags",
        "lambda:TagResource",
        "lambda:UntagResource"
      ],
      "Resource": [
        "arn:aws:lambda:${AWS_REGION}:${AWS_ACCOUNT_ID}:function:trackstack-*"
      ]
    },
    {
      "Sid": "APIGateway",
      "Effect": "Allow",
      "Action": [
        "apigateway:GET",
        "apigateway:POST",
        "apigateway:PUT",
        "apigateway:PATCH",
        "apigateway:DELETE"
      ],
      "Resource": [
        "arn:aws:apigateway:${AWS_REGION}::/restapis",
        "arn:aws:apigateway:${AWS_REGION}::/restapis/*"
      ]
    },
    {
      "Sid": "CloudFront",
      "Effect": "Allow",
      "Action": [
        "cloudfront:CreateDistribution",
        "cloudfront:UpdateDistribution",
        "cloudfront:GetDistribution",
        "cloudfront:DeleteDistribution",
        "cloudfront:CreateInvalidation",
        "cloudfront:GetInvalidation",
        "cloudfront:ListTagsForResource",
        "cloudfront:TagResource"
      ],
      "Resource": "*"
    },
    {
      "Sid": "S3AppAssets",
      "Effect": "Allow",
      "Action": [
        "s3:CreateBucket",
        "s3:DeleteBucket",
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket",
        "s3:PutBucketPolicy",
        "s3:GetBucketPolicy",
        "s3:DeleteBucketPolicy",
        "s3:PutBucketWebsite",
        "s3:PutBucketCors",
        "s3:PutPublicAccessBlock"
      ],
      "Resource": "arn:aws:s3:::trackstack-*"
    },
    {
      "Sid": "IAMForLambda",
      "Effect": "Allow",
      "Action": [
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:GetRole",
        "iam:PassRole",
        "iam:PutRolePolicy",
        "iam:DeleteRolePolicy",
        "iam:GetRolePolicy",
        "iam:AttachRolePolicy",
        "iam:DetachRolePolicy",
        "iam:TagRole",
        "iam:UntagRole"
      ],
      "Resource": "arn:aws:iam::${AWS_ACCOUNT_ID}:role/trackstack-*"
    },
    {
      "Sid": "CloudWatchLogs",
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:DeleteLogGroup",
        "logs:PutRetentionPolicy",
        "logs:TagResource"
      ],
      "Resource": "arn:aws:logs:${AWS_REGION}:${AWS_ACCOUNT_ID}:log-group:/aws/lambda/trackstack-*"
    },
    {
      "Sid": "Route53",
      "Effect": "Allow",
      "Action": [
        "route53:ListHostedZones",
        "route53:GetHostedZone",
        "route53:ChangeResourceRecordSets"
      ],
      "Resource": "*"
    },
    {
      "Sid": "ACM",
      "Effect": "Allow",
      "Action": [
        "acm:ListCertificates",
        "acm:DescribeCertificate",
        "acm:RequestCertificate",
        "acm:DeleteCertificate"
      ],
      "Resource": "*"
    },
    {
      "Sid": "STS",
      "Effect": "Allow",
      "Action": "sts:GetCallerIdentity",
      "Resource": "*"
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
    echo "Role already exists. Updating policies..."
    aws iam delete-role-policy \
        --role-name "$GITHUB_ROLE_NAME" \
        --policy-name "TrackstackDeployPolicy" \
        --profile "$AWS_PROFILE" \
        2>/dev/null || true
else
    echo "Creating IAM role..."
    aws iam create-role \
        --role-name "$GITHUB_ROLE_NAME" \
        --assume-role-policy-document file:///tmp/github-trust-policy.json \
        --description "Role for GitHub Actions to deploy TrackStack infrastructure" \
        --tags Key=Project,Value=trackstack Key=ManagedBy,Value=bootstrap Key=Purpose,Value=github-deployment \
        --profile "$AWS_PROFILE"
fi

echo "Attaching permissions policy..."
aws iam put-role-policy \
    --role-name "$GITHUB_ROLE_NAME" \
    --policy-name "TrackstackDeployPolicy" \
    --policy-document file:///tmp/github-permissions-policy.json \
    --profile "$AWS_PROFILE"

rm -f /tmp/github-trust-policy.json /tmp/github-permissions-policy.json

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
