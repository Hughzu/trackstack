provider "aws" {
  region = local.region

  default_tags {
    tags = {
      Project   = local.app-name
      ManagedBy = "terraform"
    }
  }
}
