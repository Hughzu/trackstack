resource "aws_ssm_parameter" "infra_assets_bucket" {
  name  = "/${local.ssm_prefix}/infra/assets_bucket"
  type  = "String"
  value = local.assets_bucket_name
  tags  = var.tags
}

resource "aws_ssm_parameter" "infra_artifacts_bucket" {
  name  = "/${local.ssm_prefix}/infra/artifacts_bucket"
  type  = "String"
  value = local.artifacts_bucket_name
  tags  = var.tags
}

resource "aws_ssm_parameter" "infra_lambda_function_name" {
  name  = "/${local.ssm_prefix}/infra/lambda_function_name"
  type  = "String"
  value = aws_lambda_function.ssr.function_name
  tags  = var.tags
}

resource "aws_ssm_parameter" "infra_lambda_key" {
  name  = "/${local.ssm_prefix}/infra/lambda_key"
  type  = "String"
  value = local.lambda_artifact_key
  tags  = var.tags
}

resource "aws_ssm_parameter" "infra_cloudfront_distribution_id" {
  name  = "/${local.ssm_prefix}/infra/cloudfront_distribution_id"
  type  = "String"
  value = aws_cloudfront_distribution.ssr.id
  tags  = var.tags
}
