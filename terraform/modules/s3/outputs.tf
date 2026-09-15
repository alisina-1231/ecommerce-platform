# 7. Output the Website URL
output "website_url" {
  value       = aws_s3_bucket_website_configuration.website_config.website_endpoint
  description = "The public URL of your static website"
}
