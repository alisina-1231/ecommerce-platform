

output "repository_arns" {
  description = "ECR repository ARNs"
  value = {
    for name, repo in aws_ecr_repository.this :
    name => repo.arn
  }
}

output "repository_urls" {
  value = {
    for name, repo in aws_ecr_repository.this :
    name => repo.repository_url
  }
}