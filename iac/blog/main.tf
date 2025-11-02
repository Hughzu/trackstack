provider "aws" {
  region = local.region

  default_tags {
    tags = {
      Project   = local.app-name
      ManagedBy = "terraform"
    }
  }
}

module "blog_budget" {
  source = "../common/modules/budget"

  budget_name          = "trackstack-blog-budget"
  monthly_budget_limit = "3.00"
  alert_emails         = ["hughes.stiernon@gmail.com"]

  cost_filters = {
    "TagKeyValue" = ["user:App$blog"]
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