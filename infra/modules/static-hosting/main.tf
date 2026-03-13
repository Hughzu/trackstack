locals {
  lambda_origin_id = "lambda-ssr"
  assets_origin_id = "s3-assets"
  lambda_origin_domain = replace(
    replace(var.lambda_function_url, "https://", ""),
    "/",
    ""
  )
  default_origin_is_s3             = var.default_origin == "s3"
  default_target_origin_id         = local.default_origin_is_s3 ? local.assets_origin_id : local.lambda_origin_id
  default_allowed_methods          = local.default_origin_is_s3 ? ["GET", "HEAD", "OPTIONS"] : ["GET", "HEAD", "OPTIONS", "PUT", "PATCH", "POST", "DELETE"]
  default_cached_methods           = local.default_origin_is_s3 ? ["GET", "HEAD", "OPTIONS"] : ["GET", "HEAD"]
  default_cache_policy_id          = local.default_origin_is_s3 ? data.aws_cloudfront_cache_policy.caching_optimized.id : data.aws_cloudfront_cache_policy.caching_disabled.id
  default_origin_request_policy_id = local.default_origin_is_s3 ? data.aws_cloudfront_origin_request_policy.cors_s3.id : data.aws_cloudfront_origin_request_policy.all_viewer.id
}

resource "aws_s3_bucket" "assets" {
  bucket = var.assets_bucket_name
  tags   = var.tags
}

resource "aws_s3_bucket_versioning" "assets" {
  bucket = aws_s3_bucket.assets.id
  versioning_configuration {
    status = "Suspended"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "assets" {
  bucket = aws_s3_bucket.assets.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "assets" {
  bucket                  = aws_s3_bucket.assets.id
  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_ownership_controls" "assets" {
  bucket = aws_s3_bucket.assets.id
  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

data "aws_cloudfront_cache_policy" "caching_disabled" {
  name = "Managed-CachingDisabled"
}

data "aws_cloudfront_cache_policy" "caching_optimized" {
  name = "Managed-CachingOptimized"
}

data "aws_cloudfront_origin_request_policy" "all_viewer" {
  name = "Managed-AllViewerExceptHostHeader"
}

data "aws_cloudfront_origin_request_policy" "cors_s3" {
  name = "Managed-CORS-S3Origin"
}

resource "aws_cloudfront_response_headers_policy" "security_headers" {
  name = "${var.resource_prefix}-security-headers"

  security_headers_config {
    strict_transport_security {
      access_control_max_age_sec = 31536000
      include_subdomains         = true
      override                   = true
      preload                    = true
    }
    xss_protection {
      mode_block = true
      protection = true
      override   = true
    }
    frame_options {
      frame_option = "SAMEORIGIN"
      override     = true
    }
    content_type_options {
      override = true
    }
    referrer_policy {
      referrer_policy = "strict-origin-when-cross-origin"
      override        = true
    }
  }
}

resource "aws_cloudfront_origin_access_control" "assets" {
  name                              = "${var.resource_prefix}-assets-oac"
  description                       = "OAC for S3 assets"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_origin_access_control" "lambda" {
  name                              = "${var.resource_prefix}-lambda-oac"
  description                       = "OAC for Lambda function URL"
  origin_access_control_origin_type = "lambda"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_function" "directory_index_rewrite" {
  count   = var.enable_directory_index_rewrite ? 1 : 0
  name    = "${var.resource_prefix}-directory-index-rewrite"
  runtime = "cloudfront-js-2.0"
  publish = true
  comment = "Rewrite extensionless routes to index.html objects"
  code    = <<-EOF
  function handler(event) {
    var request = event.request;
    var uri = request.uri || "/";

    if (uri === "/") {
      request.uri = "/index.html";
      return request;
    }

    if (uri.endsWith("/")) {
      request.uri = uri + "index.html";
      return request;
    }

    if (uri.indexOf(".") === -1) {
      request.uri = uri + "/index.html";
    }

    return request;
  }
  EOF
}

resource "aws_cloudfront_distribution" "ssr" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = var.distribution_comment
  price_class         = var.price_class
  default_root_object = local.default_origin_is_s3 ? "index.html" : null

  origin {
    domain_name              = aws_s3_bucket.assets.bucket_regional_domain_name
    origin_id                = local.assets_origin_id
    origin_access_control_id = aws_cloudfront_origin_access_control.assets.id
  }

  origin {
    domain_name              = local.lambda_origin_domain
    origin_id                = local.lambda_origin_id
    origin_access_control_id = aws_cloudfront_origin_access_control.lambda.id

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }

    custom_header {
      name  = var.origin_header_name
      value = var.origin_header_value
    }
  }

  default_cache_behavior {
    target_origin_id           = local.default_target_origin_id
    viewer_protocol_policy     = "redirect-to-https"
    allowed_methods            = local.default_allowed_methods
    cached_methods             = local.default_cached_methods
    cache_policy_id            = local.default_cache_policy_id
    origin_request_policy_id   = local.default_origin_request_policy_id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers.id
    compress                   = true

    dynamic "function_association" {
      for_each = var.enable_directory_index_rewrite ? [1] : []
      content {
        event_type   = "viewer-request"
        function_arn = aws_cloudfront_function.directory_index_rewrite[0].arn
      }
    }
  }

  dynamic "ordered_cache_behavior" {
    for_each = var.lambda_path_patterns
    content {
      path_pattern               = ordered_cache_behavior.value
      target_origin_id           = local.lambda_origin_id
      viewer_protocol_policy     = "redirect-to-https"
      allowed_methods            = ["GET", "HEAD", "OPTIONS", "PUT", "PATCH", "POST", "DELETE"]
      cached_methods             = ["GET", "HEAD"]
      cache_policy_id            = data.aws_cloudfront_cache_policy.caching_disabled.id
      origin_request_policy_id   = data.aws_cloudfront_origin_request_policy.all_viewer.id
      response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers.id
      compress                   = true
    }
  }

  dynamic "ordered_cache_behavior" {
    for_each = var.s3_path_patterns
    content {
      path_pattern               = ordered_cache_behavior.value
      target_origin_id           = local.assets_origin_id
      viewer_protocol_policy     = "redirect-to-https"
      allowed_methods            = ["GET", "HEAD", "OPTIONS"]
      cached_methods             = ["GET", "HEAD", "OPTIONS"]
      cache_policy_id            = data.aws_cloudfront_cache_policy.caching_optimized.id
      origin_request_policy_id   = data.aws_cloudfront_origin_request_policy.cors_s3.id
      response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers.id
      compress                   = true
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = var.tags
}

data "aws_iam_policy_document" "assets_bucket" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.assets.arn}/*"]
    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"
      values   = [aws_cloudfront_distribution.ssr.arn]
    }
  }
}

resource "aws_s3_bucket_policy" "assets" {
  bucket = aws_s3_bucket.assets.id
  policy = data.aws_iam_policy_document.assets_bucket.json
}

resource "aws_lambda_permission" "allow_cloudfront" {
  statement_id           = "AllowCloudFrontInvokeUrl"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = var.lambda_function_name
  principal              = "cloudfront.amazonaws.com"
  source_arn             = aws_cloudfront_distribution.ssr.arn
  function_url_auth_type = "AWS_IAM"
}

resource "aws_lambda_permission" "allow_cloudfront_invoke" {
  statement_id  = "AllowCloudFrontInvokeFunction"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name
  principal     = "cloudfront.amazonaws.com"
  source_arn    = aws_cloudfront_distribution.ssr.arn
}
