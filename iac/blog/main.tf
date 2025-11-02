provider "aws" {
  region = local.region

  default_tags {
    tags = {
      Project   = local.app-name
      ManagedBy = "terraform"
    }
  }
}

module "blog" {
  source = "../common/modules/static-website"
  bucket_name = "trackstack-blog"

  description       = "TrackStack - Blog"
  enable_versioning = true
  price_class       = "PriceClass_100"

  tags = {
    Environment = "production"
    App     = local.app-name
  }
}