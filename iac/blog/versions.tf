terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  backend "s3" {
    bucket  = "trackstack-terraform-state"
    key     = "blog/terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
  }
}