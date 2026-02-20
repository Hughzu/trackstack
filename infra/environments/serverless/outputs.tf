output "cloudfront_domain_name" {
  value = aws_cloudfront_distribution.ssr.domain_name
}

output "assets_bucket" {
  value = local.assets_bucket_name
}

output "artifacts_bucket" {
  value = local.artifacts_bucket_name
}

output "lambda_function_name" {
  value = aws_lambda_function.ssr.function_name
}

output "lambda_function_url" {
  value = aws_lambda_function_url.ssr.function_url
}

output "cloudfront_distribution_id" {
  value = aws_cloudfront_distribution.ssr.id
}
