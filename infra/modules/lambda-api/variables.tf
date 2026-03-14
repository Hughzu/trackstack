variable "aws_region" {
  description = "AWS region where the Lambda function and SSM parameters live."
  type        = string
}

variable "resource_prefix" {
  description = "Prefix applied to IAM, logs, and related Lambda resources."
  type        = string
}

variable "lambda_function_name" {
  description = "Name of the Lambda function to create or update."
  type        = string
}

variable "lambda_handler" {
  description = "Lambda handler value for the deployed runtime."
  type        = string
}

variable "lambda_runtime" {
  description = "Lambda runtime identifier, for example provided.al2023 or nodejs20.x."
  type        = string
}

variable "lambda_architectures" {
  description = "Instruction set architectures supported by the Lambda function."
  type        = list(string)
  default     = ["arm64"]
}

variable "lambda_memory_size" {
  description = "Lambda memory size in MB."
  type        = number
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds."
  type        = number
}

variable "log_retention_days" {
  description = "CloudWatch log retention period in days."
  type        = number
}

variable "artifact_bucket" {
  description = "S3 bucket that stores deployable Lambda artifacts."
  type        = string
}

variable "artifact_key" {
  description = "S3 key used for the Lambda artifact object."
  type        = string
}

variable "artifact_path" {
  description = "Optional local artifact path uploaded directly to Lambda during Terraform apply."
  type        = string
  default     = null
}

variable "artifact_version" {
  description = "Optional existing S3 object version to deploy when artifact_path is not used."
  type        = string
  default     = null
}

variable "origin_header_name" {
  description = "Header name expected by the Lambda origin for CloudFront verification."
  type        = string
}

variable "origin_header_ssm_name" {
  description = "SSM parameter name that stores the shared CloudFront origin verification secret."
  type        = string
}

variable "ssm_prefix" {
  description = "Base SSM prefix for runtime parameters used by the Lambda function."
  type        = string
}

variable "runtime_ssm_parameters" {
  description = "Map of Lambda environment variable names to SSM parameter names that should be resolved at apply time."
  type        = map(string)
  default     = {}
}

variable "lambda_env" {
  description = "Additional plain Lambda environment variables merged into the runtime configuration."
  type        = map(string)
  default     = {}
}

variable "tags" {
  description = "Tags applied to all resources created by this module."
  type        = map(string)
}
