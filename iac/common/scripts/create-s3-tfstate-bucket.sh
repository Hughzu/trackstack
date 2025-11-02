#!/bin/bash

BUCKET_NAME="trackstack-terraform-state"
REGION="eu-central-1"
AWS_PROFILE="trackstack"

echo "Creating S3 bucket..."
aws s3api create-bucket \
  --bucket $BUCKET_NAME \
  --region $REGION \
  --create-bucket-configuration LocationConstraint=$REGION \
  --profile $AWS_PROFILE

echo "Enabling versioning..."
aws s3api put-bucket-versioning \
  --bucket $BUCKET_NAME \
  --versioning-configuration Status=Enabled \
  --profile $AWS_PROFILE

echo "Enabling encryption..."
aws s3api put-bucket-encryption \
  --bucket $BUCKET_NAME \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }' \
  --profile $AWS_PROFILE

echo "Blocking public access..."
aws s3api put-public-access-block \
  --bucket $BUCKET_NAME \
  --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true \
  --profile $AWS_PROFILE

echo "Setup complete!"
echo "S3 Bucket: $BUCKET_NAME"
echo "AWS Profile: $AWS_PROFILE"