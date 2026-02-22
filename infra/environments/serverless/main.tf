terraform {
  backend "s3" {
    key            = "environments/serverless/terraform.tfstate"
    encrypt        = true
    dynamodb_table = "terraform-lock"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

data "aws_caller_identity" "current" {}

locals {
  account_id            = data.aws_caller_identity.current.account_id
  assets_bucket_name    = coalesce(var.assets_bucket_name, "${var.resource_prefix}-assets-${local.account_id}")
  artifacts_bucket_name = coalesce(var.artifacts_bucket_name, "${var.resource_prefix}-artifacts-${local.account_id}")
  lambda_artifact_key   = coalesce(var.lambda_artifact_key, "astro/ssr/lambda.zip")
  ssm_prefix            = trim(var.ssm_prefix, "/")
}
