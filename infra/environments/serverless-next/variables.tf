variable "aws_region" {
  description = "AWS region for the serverless-next environment."
  type        = string
  default     = "eu-west-1"
}

variable "project_name" {
  description = "Project name tag value for the serverless-next environment."
  type        = string
  default     = "trackstack"
}

variable "resource_prefix" {
  description = "Resource naming prefix for the parallel serverless-next environment."
  type        = string
  default     = "trackstack-next"
}

variable "assets_bucket_name" {
  description = "Optional override for the static assets S3 bucket name."
  type        = string
  default     = null
}

variable "artifacts_bucket_name" {
  description = "Optional override for the Lambda artifacts S3 bucket name."
  type        = string
  default     = null
}

variable "lambda_function_name" {
  description = "Lambda function name for the Go API in the serverless-next environment."
  type        = string
  default     = "trackstack-go-api-next"
}

variable "lambda_artifact_key" {
  description = "Optional override for the S3 key used for the Go Lambda artifact."
  type        = string
  default     = null
}

variable "lambda_artifact_path" {
  description = "Optional local path to a Lambda artifact uploaded during Terraform apply."
  type        = string
  default     = null
}

variable "lambda_artifact_version" {
  description = "Optional existing S3 object version for the Lambda artifact."
  type        = string
  default     = null
}

variable "lambda_handler" {
  description = "Lambda handler value for the Go custom runtime artifact."
  type        = string
  default     = "bootstrap"
}

variable "lambda_runtime" {
  description = "Lambda runtime identifier for the Go custom runtime artifact."
  type        = string
  default     = "provided.al2023"
}

variable "origin_header_name" {
  description = "CloudFront-to-Lambda origin verification header name."
  type        = string
  default     = "X-Origin-Verify"
}

variable "origin_header_ssm_name" {
  description = "SSM parameter name used to store the origin verification secret for serverless-next."
  type        = string
  default     = "/trackstack/serverless-next/cloudfront/origin-verify"
}

variable "ssm_prefix" {
  description = "SSM prefix used for serverless-next runtime and infra parameters."
  type        = string
  default     = "/trackstack/serverless-next"
}

variable "lambda_env" {
  description = "Additional plain Lambda environment variables for serverless-next."
  type        = map(string)
  default     = {}
}

variable "billing_alarm_email" {
  description = "Email address for cost guardrail notifications."
  type        = string
  default     = "hughesstiernon@gmail.com"
}

variable "billing_budget_limit" {
  description = "Monthly AWS budget threshold in USD for serverless-next."
  type        = number
  default     = 3
}

variable "tags" {
  description = "Tags applied to all resources in the serverless-next environment."
  type        = map(string)
  default = {
    Project     = "trackstack"
    ManagedBy   = "terraform"
    Environment = "serverless-next"
    Purpose     = "static-astro-go-api-validation"
  }
}
