provider "aws" {
  region = local.region

  default_tags {
    tags = {
      Project   = local.app-name
      ManagedBy = "terraform"
    }
  }
}

resource "aws_s3_bucket" "example" {
  bucket = "my-tf-test-bucket-hsthsthst"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}