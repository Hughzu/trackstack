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
        Sid    = "ServiceAdmin"
        Effect = "Allow"
        Action = [
          "s3:*",
          "cloudfront:*",
          "lambda:*",
          "iam:*",
          "ssm:*",
          "route53:*",
          "acm:*",
          "logs:*",
          "dynamodb:*"
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
