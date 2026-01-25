variable "bucket_name" {
  description = "Name of the S3 bucket for static website (must be globally unique)"
  type        = string
}

variable "default_root_object" {
  description = "Default root object for CloudFront (usually index.html)"
  type        = string
  default     = "index.html"
}

variable "price_class" {
  description = "CloudFront price class (PriceClass_100, PriceClass_200, PriceClass_All)"
  type        = string
  default     = "PriceClass_100"
}

variable "description" {
  description = "Description for the CloudFront distribution"
  type        = string
  default     = "Static website distribution"
}

variable "enable_versioning" {
  description = "Enable S3 bucket versioning"
  type        = bool
  default     = false
}

variable "cache_policy_id" {
  description = "CloudFront managed cache policy ID (default: CachingOptimized)"
  type        = string
  default     = "658327ea-f89d-4fab-a63d-7e88639e58f6" # Managed-CachingOptimized
}

variable "origin_request_policy_id" {
  description = "CloudFront managed origin request policy ID (default: CORS-S3Origin)"
  type        = string
  default     = "88a5eaf4-2fd4-4709-b370-b4c650ea3fcf" # Managed-CORS-S3Origin
}

variable "custom_error_responses" {
  description = "Custom error responses for CloudFront (useful for SPAs)"
  type = list(object({
    error_code            = number
    response_code         = number
    response_page_path    = string
    error_caching_min_ttl = number
  }))
  default = []
}

variable "tags" {
  description = "Additional tags for resources"
  type        = map(string)
  default     = {}
}