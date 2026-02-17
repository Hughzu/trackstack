resource "aws_s3_object" "lambda_artifact" {
  count = var.lambda_artifact_path == null ? 0 : 1

  bucket = aws_s3_bucket.artifacts.id
  key    = local.lambda_artifact_key
  source = var.lambda_artifact_path

  etag = filemd5(var.lambda_artifact_path)
}
