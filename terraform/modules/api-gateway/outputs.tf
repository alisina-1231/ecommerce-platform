output "api_id" {
  description = "API Gateway HTTP API ID"
  value       = aws_apigatewayv2_api.this.id
}

output "api_endpoint" {
  description = "API Gateway endpoint"
  value       = aws_apigatewayv2_api.this.api_endpoint
}

output "api_arn" {
  description = "API Gateway ARN"
  value       = aws_apigatewayv2_api.this.arn
}

output "vpc_link_id" {
  description = "API Gateway VPC Link ID"
  value       = aws_apigatewayv2_vpc_link.this.id
}

output "authorizer_id" {
  description = "Cognito JWT authorizer ID"
  value       = aws_apigatewayv2_authorizer.cognito.id
}