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
      "arn:aws:ssm:${var.aws_region}:${local.account_id}:parameter/${local.ssm_prefix}/runtime/*"
    ]
  }
}

resource "aws_iam_role_policy" "lambda_ssm" {
  name   = "${var.resource_prefix}-astro-ssr-ssm"
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.lambda_ssm.json
}
