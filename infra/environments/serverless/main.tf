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
  lambda_artifact_key   = coalesce(var.lambda_artifact_key, "go/api/lambda.zip")
  ssm_prefix            = trim(var.ssm_prefix, "/")
  runtime_ssm_parameters = {
    APP_ENV                 = "/${local.ssm_prefix}/runtime/APP_ENV"
    LOG_LEVEL               = "/${local.ssm_prefix}/runtime/LOG_LEVEL"
    DB_CONNECTION_MODE      = "/${local.ssm_prefix}/runtime/DB_CONNECTION_MODE"
    CORS_ALLOWED_ORIGIN     = "/${local.ssm_prefix}/runtime/CORS_ALLOWED_ORIGIN"
    JWT_SECRET              = "/${local.ssm_prefix}/runtime/JWT_SECRET"
    TURSO_USERS_URL_HTTP    = "/${local.ssm_prefix}/runtime/TURSO_USERS_URL_HTTP"
    TURSO_USERS_TOKEN       = "/${local.ssm_prefix}/runtime/TURSO_USERS_TOKEN"
    TURSO_CALORIES_URL_HTTP = "/${local.ssm_prefix}/runtime/TURSO_CALORIES_URL_HTTP"
    TURSO_CALORIES_TOKEN    = "/${local.ssm_prefix}/runtime/TURSO_CALORIES_TOKEN"
    TURSO_EXPENSES_URL_HTTP = "/${local.ssm_prefix}/runtime/TURSO_EXPENSES_URL_HTTP"
    TURSO_EXPENSES_TOKEN    = "/${local.ssm_prefix}/runtime/TURSO_EXPENSES_TOKEN"
    TURSO_HEAT_URL_HTTP     = "/${local.ssm_prefix}/runtime/TURSO_HEAT_URL_HTTP"
    TURSO_HEAT_TOKEN        = "/${local.ssm_prefix}/runtime/TURSO_HEAT_TOKEN"
  }
}

module "lambda_api" {
  source = "../../modules/lambda-api"

  aws_region             = var.aws_region
  resource_prefix        = var.resource_prefix
  lambda_function_name   = var.lambda_function_name
  lambda_handler         = var.lambda_handler
  lambda_runtime         = var.lambda_runtime
  lambda_memory_size     = 512
  lambda_timeout         = 30
  log_retention_days     = 3
  artifact_bucket        = aws_s3_bucket.artifacts.bucket
  artifact_key           = local.lambda_artifact_key
  artifact_path          = var.lambda_artifact_path
  artifact_version       = var.lambda_artifact_version
  origin_header_name     = var.origin_header_name
  origin_header_ssm_name = var.origin_header_ssm_name
  ssm_prefix             = local.ssm_prefix
  runtime_ssm_parameters = local.runtime_ssm_parameters
  lambda_env             = var.lambda_env

  tags = var.tags
}

module "static_hosting" {
  source = "../../modules/static-hosting"

  resource_prefix                = var.resource_prefix
  assets_bucket_name             = local.assets_bucket_name
  origin_header_name             = var.origin_header_name
  origin_header_value            = module.lambda_api.origin_verify_value
  lambda_function_name           = module.lambda_api.lambda_function_name
  lambda_function_url            = module.lambda_api.lambda_function_url
  default_origin                 = "s3"
  lambda_path_patterns           = ["/api/*", "/health", "/openapi.yaml"]
  enable_directory_index_rewrite = true
  distribution_comment           = "Trackstack serverless static frontend + Go API"
  price_class                    = "PriceClass_100"

  tags = var.tags
}

module "cost_guardrails" {
  source = "../../modules/cost-guardrails"

  billing_alarm_email  = var.billing_alarm_email
  billing_budget_limit = var.billing_budget_limit
  resource_prefix      = var.resource_prefix

  providers = {
    aws = aws.us_east_1
  }
}
