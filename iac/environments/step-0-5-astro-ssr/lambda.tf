resource "random_password" "origin_secret" {
  length  = 32
  special = false
}

resource "aws_ssm_parameter" "origin_secret" {
  name  = var.origin_header_ssm_name
  type  = "SecureString"
  value = random_password.origin_secret.result
  tags  = var.tags
}

data "aws_ssm_parameter" "origin_secret" {
  name            = aws_ssm_parameter.origin_secret.name
  with_decryption = true
}

data "aws_ssm_parameter" "turso_users_url" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_USERS_URL"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_users_token" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_USERS_TOKEN"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_calories_url" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_CALORIES_URL"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_calories_token" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_CALORIES_TOKEN"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_expenses_url" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_EXPENSES_URL"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_expenses_token" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_EXPENSES_TOKEN"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_heat_url" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_HEAT_URL"
  with_decryption = true
}

data "aws_ssm_parameter" "turso_heat_token" {
  name            = "/${local.ssm_prefix}/runtime/TURSO_HEAT_TOKEN"
  with_decryption = true
}

resource "aws_lambda_function" "ssr" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = var.lambda_handler
  runtime       = var.lambda_runtime
  architectures = ["arm64"]

  s3_bucket         = aws_s3_bucket.artifacts.bucket
  s3_key            = local.lambda_artifact_key
  s3_object_version = var.lambda_artifact_path != null ? aws_s3_object.lambda_artifact[0].version_id : var.lambda_artifact_version

  memory_size = 1024
  timeout     = 30

  environment {
    variables = merge(var.lambda_env, {
      ORIGIN_VERIFY_HEADER = var.origin_header_name
      ORIGIN_VERIFY_VALUE  = data.aws_ssm_parameter.origin_secret.value
      TURSO_USERS_URL      = data.aws_ssm_parameter.turso_users_url.value
      TURSO_USERS_TOKEN    = data.aws_ssm_parameter.turso_users_token.value
      TURSO_CALORIES_URL   = data.aws_ssm_parameter.turso_calories_url.value
      TURSO_CALORIES_TOKEN = data.aws_ssm_parameter.turso_calories_token.value
      TURSO_EXPENSES_URL   = data.aws_ssm_parameter.turso_expenses_url.value
      TURSO_EXPENSES_TOKEN = data.aws_ssm_parameter.turso_expenses_token.value
      TURSO_HEAT_URL       = data.aws_ssm_parameter.turso_heat_url.value
      TURSO_HEAT_TOKEN     = data.aws_ssm_parameter.turso_heat_token.value
      AUTH_COOKIE_SECURE   = "true"
      AUTH_COOKIE_SAMESITE = "lax"
      AUTH_COOKIE_NAME     = "session"
    })
  }

  tags = var.tags
}

resource "aws_lambda_function_url" "ssr" {
  function_name      = aws_lambda_function.ssr.function_name
  authorization_type = "AWS_IAM"
}

resource "aws_lambda_permission" "allow_cloudfront" {
  statement_id           = "AllowCloudFrontInvokeUrl"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.ssr.function_name
  principal              = "cloudfront.amazonaws.com"
  source_arn             = aws_cloudfront_distribution.ssr.arn
  function_url_auth_type = "AWS_IAM"
}

resource "aws_lambda_permission" "allow_cloudfront_invoke" {
  statement_id  = "AllowCloudFrontInvokeFunction"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ssr.function_name
  principal     = "cloudfront.amazonaws.com"
  source_arn    = aws_cloudfront_distribution.ssr.arn
}
