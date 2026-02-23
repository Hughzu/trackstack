variable "resource_prefix" {
  type = string
}

variable "assets_bucket_name" {
  type = string
}

variable "origin_header_name" {
  type = string
}

variable "origin_header_value" {
  type = string
}

variable "lambda_function_name" {
  type = string
}

variable "lambda_function_url" {
  type = string
}

variable "price_class" {
  type    = string
  default = "PriceClass_100"
}

variable "tags" {
  type = map(string)
}
