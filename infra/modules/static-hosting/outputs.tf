output "assets_bucket_name" {
  value = aws_s3_bucket.assets.bucket
}

output "cloudfront_distribution_id" {
  value = aws_cloudfront_distribution.ssr.id
}

output "cloudfront_domain_name" {
  value = aws_cloudfront_distribution.ssr.domain_name
}
