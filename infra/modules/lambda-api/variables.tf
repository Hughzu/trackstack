variable "aws_region" {
  type = string
}

variable "resource_prefix" {
  type = string
}

variable "lambda_function_name" {
  type = string
}

variable "lambda_handler" {
  type = string
}

variable "lambda_runtime" {
  type = string
}

variable "lambda_architectures" {
  type    = list(string)
  default = ["arm64"]
}

variable "lambda_memory_size" {
  type = number
}

variable "lambda_timeout" {
  type = number
}

variable "log_retention_days" {
  type = number
}

variable "artifact_bucket" {
  type = string
}

variable "artifact_key" {
  type = string
}

variable "artifact_path" {
  type    = string
  default = null
}

variable "artifact_version" {
  type    = string
  default = null
}

variable "origin_header_name" {
  type = string
}

variable "origin_header_ssm_name" {
  type = string
}

variable "ssm_prefix" {
  type = string
}

variable "lambda_env" {
  type    = map(string)
  default = {}
}

variable "tags" {
  type = map(string)
}
