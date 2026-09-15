output "cluster_id" {
  description = "ECS cluster ID"
  value       = aws_ecs_cluster.this.id
}

output "cluster_arn" {
  description = "ECS cluster ARN"
  value       = aws_ecs_cluster.this.arn
}

output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.this.name
}

output "log_group_name" {
  description = "ECS CloudWatch log group"
  value       = aws_cloudwatch_log_group.ecs.name
}

output "task_execution_role_arn" {
  description = "ECS task execution role ARN"
  value       = aws_iam_role.ecs_task_execution.arn
}

output "task_role_arn" {
  description = "ECS application task role ARN"
  value       = aws_iam_role.ecs_task.arn
}

output "task_definition_arns" {
  description = "ECS task definition ARNs"
  value = {
    for service, task in aws_ecs_task_definition.this :
    service => task.arn
  }
}

output "task_definition_families" {
  description = "ECS task definition families"
  value = {
    for service, task in aws_ecs_task_definition.this :
    service => task.family
  }
}
output "service_arns" {
  description = "ECS service ARNs by service"
  value = {
    for service, ecs_service in aws_ecs_service.this :
    service => ecs_service.id
  }
}

output "service_names" {
  description = "ECS service names by service"
  value = {
    for service, ecs_service in aws_ecs_service.this :
    service => ecs_service.name
  }
}