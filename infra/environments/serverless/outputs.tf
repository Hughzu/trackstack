output "cloudfront_domain_name" {
  value = module.static_hosting.cloudfront_domain_name
}

output "assets_bucket" {
  value = module.static_hosting.assets_bucket_name
}

output "artifacts_bucket" {
  value = local.artifacts_bucket_name
}

output "lambda_function_name" {
  value = module.lambda_api.lambda_function_name
}

output "lambda_function_url" {
  value = module.lambda_api.lambda_function_url
}

output "cloudfront_distribution_id" {
  value = module.static_hosting.cloudfront_distribution_id
}
