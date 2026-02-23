terraform {
  backend "s3" {
    key = "environments/iam-bootstrap/terraform.tfstate"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_caller_identity" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
}

data "aws_iam_role" "terraform" {
  name = var.role_name
}

resource "aws_iam_policy" "terraform_deploy" {
  name        = var.policy_name
  description = "Least-privilege policy for Trackstack Terraform deployments"
  tags        = var.tags

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3Scoped"
        Effect = "Allow"
        Action = [
          "s3:*"
        ]
        Resource = [
          "arn:aws:s3:::${var.resource_prefix}*",
          "arn:aws:s3:::${var.resource_prefix}*/*"
        ]
      },
      {
        Sid    = "LambdaScoped"
        Effect = "Allow"
        Action = [
          "lambda:*"
        ]
        Resource = [
          "arn:aws:lambda:${var.aws_region}:${local.account_id}:function:${var.resource_prefix}*"
        ]
      },
      {
        Sid    = "SSMScoped"
        Effect = "Allow"
        Action = [
          "ssm:*"
        ]
        Resource = [
          "arn:aws:ssm:${var.aws_region}:${local.account_id}:parameter/${trim(var.ssm_parameter_prefix, "/")}*"
        ]
      },
      {
        Sid    = "SSMDescribe"
        Effect = "Allow"
        Action = [
          "ssm:DescribeParameters"
        ]
        Resource = "*"
      },
      {
        Sid    = "LogsScoped"
        Effect = "Allow"
        Action = [
          "logs:*"
        ]
        Resource = [
          "arn:aws:logs:${var.aws_region}:${local.account_id}:log-group:/aws/lambda/${var.resource_prefix}*:*"
        ]
      },
      {
        Sid    = "LogsDescribe"
        Effect = "Allow"
        Action = [
          "logs:DescribeLogGroups",
          "logs:ListTagsForResource"
        ]
        Resource = "*"
      },
      {
        Sid    = "DynamoDbScoped"
        Effect = "Allow"
        Action = [
          "dynamodb:*"
        ]
        Resource = [
          "arn:aws:dynamodb:${var.aws_region}:${local.account_id}:table/${var.resource_prefix}*"
        ]
      },
      {
        Sid    = "IamScoped"
        Effect = "Allow"
        Action = [
          "iam:CreateRole",
          "iam:DeleteRole",
          "iam:UpdateRole",
          "iam:AttachRolePolicy",
          "iam:DetachRolePolicy",
          "iam:PutRolePolicy",
          "iam:DeleteRolePolicy",
          "iam:CreatePolicy",
          "iam:DeletePolicy",
          "iam:CreatePolicyVersion",
          "iam:DeletePolicyVersion",
          "iam:SetDefaultPolicyVersion",
          "iam:PassRole",
          "iam:TagRole",
          "iam:UntagRole",
          "iam:TagPolicy",
          "iam:UntagPolicy"
        ]
        Resource = [
          "arn:aws:iam::${local.account_id}:role/${var.resource_prefix}*",
          "arn:aws:iam::${local.account_id}:policy/${var.resource_prefix}*"
        ]
      },
      {
        Sid    = "IamReadOnly"
        Effect = "Allow"
        Action = [
          "iam:Get*",
          "iam:List*"
        ]
        Resource = "*"
      },
      {
        Sid    = "TerraformReadOnly"
        Effect = "Allow"
        Action = [
          "s3:Get*",
          "s3:List*",
          "lambda:Get*",
          "lambda:List*",
          "cloudfront:Get*",
          "cloudfront:List*",
          "logs:Describe*",
          "logs:Get*",
          "logs:List*",
          "ssm:Describe*",
          "ssm:Get*",
          "ssm:List*",
          "dynamodb:Describe*",
          "dynamodb:Get*",
          "dynamodb:List*",
          "budgets:Describe*",
          "budgets:List*",
          "budgets:ViewBudget"
        ]
        Resource = "*"
      },
      {
        Sid    = "CloudFrontAll"
        Effect = "Allow"
        Action = [
          "cloudfront:*"
        ]
        Resource = "*"
      },
      {
        Sid    = "BudgetsAll"
        Effect = "Allow"
        Action = [
          "budgets:CreateBudget",
          "budgets:ModifyBudget",
          "budgets:DeleteBudget",
          "budgets:DescribeBudgets",
          "budgets:DescribeBudget",
          "budgets:DescribeNotificationsForBudget",
          "budgets:DescribeSubscribersForBudget",
          "budgets:ListTagsForResource",
          "budgets:ViewBudget"
        ]
        Resource = "*"
      },
      {
        Sid      = "STS"
        Effect   = "Allow"
        Action   = "sts:GetCallerIdentity"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "terraform_deploy" {
  role       = data.aws_iam_role.terraform.name
  policy_arn = aws_iam_policy.terraform_deploy.arn
}
