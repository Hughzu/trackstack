variable "resource_prefix" {
  description = "Prefix applied to CloudFront and S3 resources."
  type        = string
}

variable "assets_bucket_name" {
  description = "Name of the S3 bucket that stores frontend assets and HTML files."
  type        = string
}

variable "origin_header_name" {
  description = "Header name CloudFront sends to the Lambda origin for origin verification."
  type        = string
}

variable "origin_header_value" {
  description = "Header value CloudFront sends to the Lambda origin for origin verification."
  type        = string
}

variable "lambda_function_name" {
  description = "Name of the Lambda function behind the Function URL origin."
  type        = string
}

variable "lambda_function_url" {
  description = "Lambda Function URL used as the CloudFront API origin."
  type        = string
}

variable "default_origin" {
  description = "Default CloudFront origin. Use 'lambda' for SSR-style routing or 's3' for static-first routing."
  type        = string
  default     = "lambda"

  validation {
    condition     = contains(["lambda", "s3"], var.default_origin)
    error_message = "default_origin must be either 'lambda' or 's3'."
  }
}

variable "lambda_path_patterns" {
  description = "CloudFront path patterns that should be routed to the Lambda origin."
  type        = list(string)
  default     = []
}

variable "s3_path_patterns" {
  description = "CloudFront path patterns that should be routed to the S3 origin."
  type        = list(string)
  default = [
    "/assets/*",
    "/sw.js",
    "/manifest.webmanifest",
    "/registerSW.js",
    "/workbox-*.js",
    "/favicon.svg",
    "/icons.svg",
    "/android/*",
    "/ios/*"
  ]
}

variable "enable_directory_index_rewrite" {
  description = "When true, rewrites SPA routes without file extensions to the shared '/index.html' app shell via CloudFront Function."
  type        = bool
  default     = false
}

variable "distribution_comment" {
  description = "Comment applied to the CloudFront distribution."
  type        = string
  default     = "Trackstack application"
}

variable "price_class" {
  description = "CloudFront price class to control edge footprint and cost."
  type        = string
  default     = "PriceClass_100"
}

variable "tags" {
  description = "Tags applied to all resources created by this module."
  type        = map(string)
}
