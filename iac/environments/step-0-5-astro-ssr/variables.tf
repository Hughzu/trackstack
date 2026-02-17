variable "aws_region" {
  type    = string
  default = "eu-west-1"
}

variable "project_name" {
  type    = string
  default = "trackstack"
}

variable "resource_prefix" {
  type    = string
  default = "trackstack"
}

variable "assets_bucket_name" {
  type    = string
  default = null
}

variable "artifacts_bucket_name" {
  type    = string
  default = null
}

variable "lambda_function_name" {
  type    = string
  default = "trackstack-astro-ssr"
}

variable "lambda_artifact_key" {
  type    = string
  default = null
}

variable "lambda_artifact_path" {
  type    = string
  default = null
}

variable "lambda_artifact_version" {
  type    = string
  default = null
}

variable "lambda_handler" {
  type    = string
  default = "dist/server/entry.handler"
}

variable "lambda_runtime" {
  type    = string
  default = "nodejs20.x"
}

variable "origin_header_name" {
  type    = string
  default = "X-Origin-Verify"
}

variable "origin_header_ssm_name" {
  type    = string
  default = "/trackstack/cloudfront/origin-verify"
}

variable "ssm_prefix" {
  type    = string
  default = "/trackstack/step-0-5"
}

variable "lambda_env" {
  type    = map(string)
  default = {}
}

variable "tags" {
  type = map(string)
  default = {
    Project   = "trackstack"
    ManagedBy = "terraform"
    Purpose   = "astro-ssr"
  }
}
