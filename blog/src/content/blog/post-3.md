---
title: '3. Deploying the Blog: Infrastructure as Code with OIDC'
description: 'Setting up S3, CloudFront, and GitHub Actions OIDC for secure deployments'
pubDate: '2025-11-15'
heroImage: '../../assets/hero-gradient.svg'
---

## From Localhost to Production

The blog you're reading right now is deployed on AWS infrastructure that I provisioned with Terraform and deploy automatically via GitHub Actions. This post walks through the infrastructure setup, the routing challenges I solved, and the OIDC implementation that eliminates the need for long-lived AWS credentials.

**Spoiler:** This is the boring, pragmatic foundation that everything else will build on. No microservices yet. Just solid infrastructure that costs pennies per month.

## The Infrastructure: S3 + CloudFront

The setup is straightforward: static site hosting with S3 as the origin and CloudFront as the CDN. Here's what I deployed:

### S3 Bucket Configuration

The S3 bucket is locked down tight - no public access at all. CloudFront accesses it via Origin Access Control (OAC), which is AWS's modern replacement for Origin Access Identity.

```hcl
resource "aws_s3_bucket" "website" {
  bucket = var.bucket_name
  
  tags = merge(
    var.tags,
    {
      Name    = var.bucket_name
      Purpose = "Static Website Hosting"
    }
  )
}

resource "aws_s3_bucket_public_access_block" "website" {
  bucket = aws_s3_bucket.website.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
```

Security essentials:
- **Versioning** enabled (optional, but good practice)
- **AES256 encryption** at rest
- **Completely private** - CloudFront is the only way in

### CloudFront Distribution

CloudFront provides:
- Global CDN with HTTPS (using CloudFront's default certificate)
- Gzip compression
- Geographic restrictions (whitelisted to EU countries for now)
- Cache optimization via managed policies

```hcl
resource "aws_cloudfront_distribution" "website" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = "PriceClass_100"  # US, Canada, Europe

  origin {
    domain_name              = aws_s3_bucket.website.bucket_regional_domain_name
    origin_id                = "S3-${var.bucket_name}"
    origin_access_control_id = aws_cloudfront_origin_access_control.website.id
  }

  default_cache_behavior {
    allowed_methods          = ["GET", "HEAD", "OPTIONS"]
    cached_methods           = ["GET", "HEAD"]
    target_origin_id         = "S3-${var.bucket_name}"
    viewer_protocol_policy   = "redirect-to-https"
    compress                 = true
    cache_policy_id          = var.cache_policy_id
    origin_request_policy_id = var.origin_request_policy_id
  }
  
  # More config...
}
```

**Cost breakdown:**
- S3 storage: ~€0.021/GB/month
- CloudFront requests: First 10TB free tier, then ~€0.078/GB
- **Total expected cost: <€1/month**

## The Astro Routing Problem

Here's where things got interesting. Astro generates clean URLs without `.html` extensions. So `/blog/post-1` should serve `/blog/post-1/index.html`. But CloudFront doesn't know this out of the box.

When you request `/blog/post-1`:
- CloudFront looks for an S3 object at that exact key
- S3 says "404, that doesn't exist"
- User sees an error page

The solution? **CloudFront Functions** - lightweight JavaScript that runs on every request at CloudFront edge locations.

### The URL Rewrite Function

I created a CloudFront Function that appends `index.html` to directory-style requests:

```javascript
function handler(event) {
    var request = event.request;
    var uri = request.uri;
    
    // If URL doesn't have a file extension and doesn't end with /
    if (!uri.includes('.') && !uri.endsWith('/')) {
        request.uri = uri + '/index.html';
    } 
    // If URL ends with /, append index.html
    else if (uri.endsWith('/')) {
        request.uri = uri + 'index.html';
    }
    
    return request;
}
```

This function:
- Runs on every viewer request (before hitting the cache)
- Adds negligible latency (sub-millisecond)
- Costs virtually nothing (first 2 million invocations free)
- Makes Astro routing "just work"

Now when you visit `/blog/post-1`, CloudFront requests `/blog/post-1/index.html` from S3, and everything works beautifully.

```hcl
resource "aws_cloudfront_function" "url_rewrite" {
  name    = "${var.bucket_name}-url-rewrite"
  runtime = "cloudfront-js-1.0"
  comment = "Append index.html to directory requests for Astro static site"
  publish = true
  code    = file("${path.module}/url-rewrite.js")
}
```

Attached to the distribution:

```hcl
default_cache_behavior {
  # ... other config ...
  
  function_association {
    event_type   = "viewer-request"
    function_arn = aws_cloudfront_function.url_rewrite.arn
  }
}
```

**Why not Lambda@Edge?** CloudFront Functions are cheaper, faster, and sufficient for simple URL rewrites. Lambda@Edge is overkill here.

## Making It Reusable: Terraform Modules

This infrastructure isn't just for the blog - it's a template I'll use for every static site in TrackStack. So I extracted it into a reusable Terraform module.

The module lives in `iac/common/modules/static-website/` and encapsulates all the S3 + CloudFront logic. Using it is dead simple:

```hcl
module "blog" {
  source = "../common/modules/static-website"
  
  bucket_name       = "trackstack-blog"
  description       = "TrackStack - Blog"
  enable_versioning = true
  price_class       = "PriceClass_100"
  
  tags = {
    Environment = "production"
    App         = "trackstack"
  }
}
```

**Benefits:**
- **DRY principle** - Define the infrastructure pattern once, reuse everywhere
- **Consistency** - Every static site gets the same security posture and optimizations
- **Easy iteration** - Improve the module, all sites benefit automatically
- **Clear separation** - Application-specific config (bucket name, tags) separated from infrastructure logic

When I deploy the Angular micro-frontends later, each one gets its own S3+CloudFront setup with just 10 lines of Terraform. No copy-paste, no drift.

## Secure Deployments with OIDC

Here's where modern DevOps gets good: **no AWS access keys stored in GitHub**.

Traditional approach:
- Create IAM user, generate access keys
- Store in GitHub Secrets, rotate manually every 90 days
- Hope nobody leaked them

OIDC approach:
- GitHub Actions assumes an IAM role temporarily during workflow runs
- No long-lived credentials anywhere
- Principle of least privilege with scoped permissions

### Scoped IAM Role for Blog Deployment

I created a dedicated role that can **only** deploy the blog - nothing else:

```bash
# Trust policy: Only this GitHub repo can assume the role
{
  "Effect": "Allow",
  "Principal": {
    "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
  },
  "Action": "sts:AssumeRoleWithWebIdentity",
  "Condition": {
    "StringLike": {
      "token.actions.githubusercontent.com:sub": "repo:Hughzu/trackstack:*"
    }
  }
}
```

```bash
# Permissions policy: Only S3 + CloudFront for the blog
{
  "Statement": [
    {
      "Sid": "S3BlogDeployment",
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:GetObject", "s3:DeleteObject", "s3:ListBucket"],
      "Resource": ["arn:aws:s3:::trackstack-blog", "arn:aws:s3:::trackstack-blog/*"]
    },
    {
      "Sid": "CloudFrontInvalidation",
      "Effect": "Allow",
      "Action": ["cloudfront:CreateInvalidation", "cloudfront:GetInvalidation"],
      "Resource": "arn:aws:cloudfront::ACCOUNT_ID:distribution/DISTRIBUTION_ID"
    }
  ]
}
```

**Security wins:**
- Role is locked to my specific GitHub repository
- Can only touch the blog's S3 bucket and CloudFront distribution
- Can't accidentally delete production databases or modify other infrastructure
- Separate role for Terraform deployments (with broader permissions)

## The Deployment Workflow

Two separate workflows handle infrastructure and application deployments:

### 1. Infrastructure Deployment (`deploy-blog-iac.yml`)

Handles Terraform changes:

```yaml
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v4
  with:
    role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
    role-session-name: GitHubActions-Terraform-${{ github.run_number }}
    aws-region: ${{ secrets.AWS_REGION }}

- name: Terraform Plan
  run: |
    terraform plan -detailed-exitcode -out=tfplan

- name: Terraform Apply (Auto on Push to Main)
  if: github.ref == 'refs/heads/main' && github.event_name == 'push'
  run: |
    terraform apply tfplan
```

**Workflow logic:**
- Always runs `terraform plan` on PR or push
- Auto-applies on push to `main` branch
- Manual trigger available for `apply` or `destroy`

### 2. Application Deployment (`deploy-blog-to-s3.yml`)

Handles blog content updates:

```yaml
- name: Build Astro site
  run: npm run build

- name: Deploy to S3
  run: |
    aws s3 sync dist/ s3://${{ secrets.BLOG_S3_BUCKET_NAME }} --delete

- name: Invalidate CloudFront cache
  run: |
    aws cloudfront create-invalidation \
      --distribution-id ${{ secrets.BLOG_CLOUDFRONT_DISTRIBUTION_ID }} \
      --paths "/*"
```

**Key operations:**
1. Build the static site with Astro
2. Sync to S3 (using `--delete` to remove old files)
3. Invalidate CloudFront cache so changes are immediate
4. Wait for invalidation to complete

**Cache invalidation cost:** First 1,000 paths per month are free. After that, it's €0.0046 per path. My `/*` wildcard invalidation counts as 1 path.

## Terraform Outputs

The infrastructure exports everything needed for the application deployment:

```hcl
output "cloudfront_domain_name" {
  description = "Domain name of the CloudFront distribution"
  value       = aws_cloudfront_distribution.website.domain_name
}

output "cloudfront_distribution_id" {
  description = "ID of the CloudFront distribution (needed for cache invalidation)"
  value       = aws_cloudfront_distribution.website.id
}

output "website_url" {
  description = "Full HTTPS URL to access the website"
  value       = "https://${aws_cloudfront_distribution.website.domain_name}"
}
```

These outputs feed into GitHub Secrets for the application deployment workflow.

## What's Next

Right now, I'm using CloudFront's default domain (`d1234abcd.cloudfront.net`). Eventually I'll add Route 53 for DNS and AWS Certificate Manager for a custom domain with proper SSL. But that's future work.

**Current priority:** Building the actual TrackStack application. The blog infrastructure is deployed and working - time to focus on the features that matter.

---

**Cost to date:**
- Infrastructure setup: €0 (one-time manual OIDC setup)
- Monthly running costs: <€1 (S3 + CloudFront)
- Domain + SSL: €0 (not implemented yet)

**Total: Cheaper than a coffee.**