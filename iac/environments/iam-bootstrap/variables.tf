variable "aws_region" {
  type    = string
  default = "eu-west-1"
}

variable "role_name" {
  type    = string
  default = "trackstack-github-deploy-role"
}

variable "policy_name" {
  type    = string
  default = "trackstack-terraform-deploy-policy"
}

variable "tfstate_bucket" {
  type = string
}

variable "tfstate_lock_table" {
  type = string
}

variable "resource_prefix" {
  type    = string
  default = "trackstack-"
}

variable "tags" {
  type = map(string)
  default = {
    Project   = "trackstack"
    ManagedBy = "terraform"
    Purpose   = "terraform-deploy"
  }
}
