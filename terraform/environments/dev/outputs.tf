
output "ecr_repository_urls" {
  description = "ECR repository URLs"
  value       = module.ecr.repository_urls
}

output "ecs_cluster_name" {
  value = module.ecs.cluster_name
}

output "ecs_cluster_arn" {
  value = module.ecs.cluster_arn
}

output "ecs_log_group_name" {
  value = module.ecs.log_group_name
}

output "ecs_task_execution_role_arn" {
  value = module.ecs.task_execution_role_arn
}

output "ecs_task_role_arn" {
  value = module.ecs.task_role_arn
}

output "ecs_task_definition_arns" {
  value = module.ecs.task_definition_arns
}

output "ecs_task_definition_families" {
  value = module.ecs.task_definition_families
}

output "alb_security_group_id" {
  value = module.security_groups.alb_security_group_id
}

output "ecs_security_group_id" {
  value = module.security_groups.ecs_security_group_id
}
output "ecr_api_endpoint_id" {
  value = module.vpc_endpoints.ecr_api_endpoint_id
}

output "ecr_dkr_endpoint_id" {
  value = module.vpc_endpoints.ecr_dkr_endpoint_id
}

output "logs_endpoint_id" {
  value = module.vpc_endpoints.logs_endpoint_id
}

output "s3_endpoint_id" {
  value = module.vpc_endpoints.s3_endpoint_id
}

output "vpc_endpoint_security_group_id" {
  value = module.vpc_endpoints.endpoint_security_group_id
}
output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = module.alb.alb_dns_name
}

output "alb_arn" {
  description = "Application Load Balancer ARN"
  value       = module.alb.alb_arn
}

output "alb_listener_arn" {
  description = "ALB HTTP listener ARN"
  value       = module.alb.listener_arn
}

output "alb_target_group_arns" {
  description = "ALB target group ARNs"
  value       = module.alb.target_group_arns
}

output "ecs_service_arns" {
  description = "ECS service ARNs"
  value       = module.ecs.service_arns
}

output "ecs_service_names" {
  description = "ECS service names"
  value       = module.ecs.service_names
}

output "postgres_endpoint" {
  description = "PostgreSQL RDS endpoint"
  value       = module.rds.endpoint
}

output "postgres_port" {
  description = "PostgreSQL port"
  value       = module.rds.port
}

output "postgres_database" {
  description = "PostgreSQL database"
  value       = module.rds.database_name
}

output "cart_dynamodb_table" {
  description = "Cart DynamoDB table"
  value       = module.cart_dynamodb.name
}

output "cart_dynamodb_arn" {
  description = "Cart DynamoDB table ARN"
  value       = module.cart_dynamodb.arn
}