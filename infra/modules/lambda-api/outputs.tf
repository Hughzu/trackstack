output "lambda_function_name" {
  value = aws_lambda_function.ssr.function_name
}

output "lambda_function_url" {
  value = aws_lambda_function_url.ssr.function_url
}

output "origin_verify_value" {
  value     = data.aws_ssm_parameter.origin_secret.value
  sensitive = true
}
