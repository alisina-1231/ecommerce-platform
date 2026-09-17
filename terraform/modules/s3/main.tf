
# 2. Create the S3 Bucket
resource "aws_s3_bucket" "website" {
  bucket        = var.bucket_name
  force_destroy = true # Allows terraform destroy to delete the bucket even if it has files
}

# 3. Configure Website Settings
resource "aws_s3_bucket_website_configuration" "website_config" {
  bucket = aws_s3_bucket.website.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}

# 4. Disable "Block Public Access" (Required for public static websites)
resource "aws_s3_bucket_public_access_block" "public_block" {
  bucket = aws_s3_bucket.website.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# 5. Attach a Public Read Policy to the Bucket
resource "aws_s3_bucket_policy" "public_read_policy" {
  depends_on = [aws_s3_bucket_public_access_block.public_block]
  bucket     = aws_s3_bucket.website.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.website.arn}/*"
      }
    ]
  })
}

# 6. Upload Your Website Files
# Uploads index.html
resource "aws_s3_object" "index" {
  bucket       = aws_s3_bucket.website.id
  key          = "index.html"
  source       = "/home/ali/project/ecommerce-platform/frontend/index.html" # Path to your local file
  content_type = "text/html"                                                # Vital so browsers render it instead of downloading it
}

