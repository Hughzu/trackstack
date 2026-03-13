data "aws_caller_identity" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
  legacy_runtime_env = {
    TURSO_USERS_URL      = "/${var.ssm_prefix}/runtime/TURSO_USERS_URL"
    TURSO_USERS_TOKEN    = "/${var.ssm_prefix}/runtime/TURSO_USERS_TOKEN"
    TURSO_CALORIES_URL   = "/${var.ssm_prefix}/runtime/TURSO_CALORIES_URL"
    TURSO_CALORIES_TOKEN = "/${var.ssm_prefix}/runtime/TURSO_CALORIES_TOKEN"
    TURSO_EXPENSES_URL   = "/${var.ssm_prefix}/runtime/TURSO_EXPENSES_URL"
    TURSO_EXPENSES_TOKEN = "/${var.ssm_prefix}/runtime/TURSO_EXPENSES_TOKEN"
    TURSO_HEAT_URL       = "/${var.ssm_prefix}/runtime/TURSO_HEAT_URL"
    TURSO_HEAT_TOKEN     = "/${var.ssm_prefix}/runtime/TURSO_HEAT_TOKEN"
    AUTH_COOKIE_SECURE   = "true"
    AUTH_COOKIE_SAMESITE = "lax"
    AUTH_COOKIE_NAME     = "session"
  }
  runtime_env_from_ssm = {
    for key, parameter in data.aws_ssm_parameter.runtime : key => parameter.value
  }
}

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

data "aws_ssm_parameter" "runtime" {
  for_each        = var.runtime_ssm_parameters
  name            = each.value
  with_decryption = true
}

resource "aws_cloudwatch_log_group" "ssr" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = var.log_retention_days
  tags              = var.tags
}

data "aws_iam_policy_document" "lambda_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name               = "${var.resource_prefix}-astro-ssr-exec"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  tags               = var.tags
}

data "aws_iam_policy_document" "lambda_exec" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${var.aws_region}:${local.account_id}:log-group:/aws/lambda/${var.lambda_function_name}:*"]
  }
}

resource "aws_iam_role_policy" "lambda_exec" {
  name   = "${var.resource_prefix}-astro-ssr-logs"
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_exec.json
}

data "aws_iam_policy_document" "lambda_ssm" {
  statement {
    effect = "Allow"
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]
    resources = [
      "arn:aws:ssm:${var.aws_region}:${local.account_id}:parameter/${var.ssm_prefix}/runtime/*"
    ]
  }
}

resource "aws_iam_role_policy" "lambda_ssm" {
  name   = "${var.resource_prefix}-astro-ssr-ssm"
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_ssm.json
}

resource "aws_s3_object" "lambda_artifact" {
  count = var.artifact_path == null ? 0 : 1

  bucket = var.artifact_bucket
  key    = var.artifact_key
  source = var.artifact_path

  etag = filemd5(var.artifact_path)
}

resource "aws_lambda_function" "ssr" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = var.lambda_handler
  runtime       = var.lambda_runtime
  architectures = var.lambda_architectures

  s3_bucket         = var.artifact_bucket
  s3_key            = var.artifact_key
  s3_object_version = var.artifact_path != null ? aws_s3_object.lambda_artifact[0].version_id : var.artifact_version

  memory_size = var.lambda_memory_size
  timeout     = var.lambda_timeout

  environment {
    variables = merge(
      local.legacy_runtime_env,
      var.lambda_env,
      local.runtime_env_from_ssm,
      {
        ORIGIN_VERIFY_HEADER = var.origin_header_name
        ORIGIN_VERIFY_VALUE  = data.aws_ssm_parameter.origin_secret.value
      }
    )
  }

  lifecycle {
    ignore_changes = [
      s3_key,
      s3_object_version,
    ]
  }

  tags = var.tags
}

resource "aws_lambda_function_url" "ssr" {
  function_name      = aws_lambda_function.ssr.function_name
  authorization_type = "AWS_IAM"
}
