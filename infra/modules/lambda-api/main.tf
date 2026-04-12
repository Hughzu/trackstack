data "aws_caller_identity" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
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
  name               = "${var.resource_prefix}-lambda-api-exec"
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
  name   = "${var.resource_prefix}-lambda-api-logs"
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
  name   = "${var.resource_prefix}-lambda-api-ssm"
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_ssm.json
}

resource "aws_lambda_function" "ssr" {
  function_name = var.lambda_function_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = var.lambda_handler
  runtime       = var.lambda_runtime
  architectures = var.lambda_architectures

  filename          = var.artifact_path
  source_code_hash  = var.artifact_path != null ? filebase64sha256(var.artifact_path) : null
  s3_bucket         = var.artifact_path == null ? var.artifact_bucket : null
  s3_key            = var.artifact_path == null ? var.artifact_key : null
  s3_object_version = var.artifact_path == null ? var.artifact_version : null

  memory_size = var.lambda_memory_size
  timeout     = var.lambda_timeout

  environment {
    variables = merge(
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
      filename,
      source_code_hash,
      s3_bucket,
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
