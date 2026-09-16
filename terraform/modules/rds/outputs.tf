output "endpoint" {
  description = "RDS endpoint"
  value       = aws_db_instance.this.address
}

output "port" {
  description = "RDS PostgreSQL port"
  value       = aws_db_instance.this.port
}

output "database_name" {
  description = "Initial database name"
  value       = aws_db_instance.this.db_name
}

output "username" {
  description = "PostgreSQL username"
  value       = aws_db_instance.this.username
}

output "arn" {
  description = "RDS ARN"
  value       = aws_db_instance.this.arn
}

output "instance_id" {
  description = "RDS instance ID"
  value       = aws_db_instance.this.id
}