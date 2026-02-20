#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/config.sh"

echo "=========================================="
echo " STEP 2: Create Terraform State Bucket"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Region: $AWS_REGION"
echo "Bucket: $TFSTATE_BUCKET"
echo "Lock Table: $TFSTATE_LOCK_TABLE"
echo ""

echo "Creating S3 bucket..."
aws s3api create-bucket \
    --bucket "$TFSTATE_BUCKET" \
    --region "$AWS_REGION" \
    --create-bucket-configuration LocationConstraint="$AWS_REGION" \
    --profile "$AWS_PROFILE" \
    2>/dev/null || echo "Bucket already exists, continuing..."

echo "Enabling versioning..."
aws s3api put-bucket-versioning \
    --bucket "$TFSTATE_BUCKET" \
    --versioning-configuration Status=Enabled \
    --profile "$AWS_PROFILE"

echo "Enabling encryption..."
aws s3api put-bucket-encryption \
    --bucket "$TFSTATE_BUCKET" \
    --server-side-encryption-configuration '{
        "Rules": [{
            "ApplyServerSideEncryptionByDefault": {
                "SSEAlgorithm": "AES256"
            }
        }]
    }' \
    --profile "$AWS_PROFILE"

echo "Blocking public access..."
aws s3api put-public-access-block \
    --bucket "$TFSTATE_BUCKET" \
    --public-access-block-configuration \
        BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true \
    --profile "$AWS_PROFILE"

echo "Creating DynamoDB lock table..."
aws dynamodb create-table \
    --table-name "$TFSTATE_LOCK_TABLE" \
    --attribute-definitions AttributeName=LockID,AttributeType=S \
    --key-schema AttributeName=LockID,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --tags Key=Project,Value=trackstack Key=ManagedBy,Value=bootstrap \
    --profile "$AWS_PROFILE" \
    --region "$AWS_REGION" \
    2>/dev/null || echo "Table already exists, continuing..."

echo ""
echo "=========================================="
echo " Setup Complete!"
echo "=========================================="
echo ""
echo "S3 Bucket: $TFSTATE_BUCKET"
echo "DynamoDB Table: $TFSTATE_LOCK_TABLE"
echo ""
echo "Next: ./03-create-oidc-provider.sh"
echo ""
