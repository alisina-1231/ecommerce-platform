output "user_pool_id" {
  description = "Cognito User Pool ID"
  value       = aws_cognito_user_pool.this.id
}

output "user_pool_arn" {
  description = "Cognito User Pool ARN"
  value       = aws_cognito_user_pool.this.arn
}

output "user_pool_name" {
  description = "Cognito User Pool name"
  value       = aws_cognito_user_pool.this.name
}

output "client_id" {
  description = "Cognito app client ID"
  value       = aws_cognito_user_pool_client.this.id
}

output "issuer_url" {
  description = "JWT issuer URL"
  value       = "https://cognito-idp.${data.aws_region.current.region}.amazonaws.com/${aws_cognito_user_pool.this.id}"
}

output "customers_group" {
  description = "Customer group name"
  value       = aws_cognito_user_group.customers.name
}

output "admins_group" {
  description = "Admin group name"
  value       = aws_cognito_user_group.admins.name
}