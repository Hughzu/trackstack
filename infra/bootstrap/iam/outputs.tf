output "role_name" {
  value = data.aws_iam_role.terraform.name
}

output "policy_arn" {
  value = aws_iam_policy.terraform_deploy.arn
}
