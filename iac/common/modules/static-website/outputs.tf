output "bucket_name" {
  description = "Name of the S3 bucket"
  value       = aws_s3_bucket.website.id
}

output "bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.website.arn
}

output "bucket_regional_domain_name" {
  description = "Regional domain name of the S3 bucket"
  value       = aws_s3_bucket.website.bucket_regional_domain_name
}

output "cloudfront_distribution_id" {
  description = "ID of the CloudFront distribution (needed for cache invalidation)"
  value       = aws_cloudfront_distribution.website.id
}

output "cloudfront_distribution_arn" {
  description = "ARN of the CloudFront distribution"
  value       = aws_cloudfront_distribution.website.arn
}

output "cloudfront_domain_name" {
  description = "Domain name of the CloudFront distribution (use this as your website URL)"
  value       = aws_cloudfront_distribution.website.domain_name
}

output "cloudfront_hosted_zone_id" {
  description = "CloudFront Route 53 hosted zone ID (for Route 53 alias records)"
  value       = aws_cloudfront_distribution.website.hosted_zone_id
}

output "website_url" {
  description = "Full HTTPS URL to access the website"
  value       = "https://${aws_cloudfront_distribution.website.domain_name}"
}