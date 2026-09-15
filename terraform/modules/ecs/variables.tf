variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "log_retention_days" {
  description = "CloudWatch log retention period"
  type        = number
  default     = 14
}

variable "services" {
  description = "ECS service names"
  type        = list(string)

  default = [
    "product",
    "cart",
    "checkout",
    "order",
    "payment"
  ]
}

variable "service_cpu" {
  description = "CPU units for each ECS task"
  type        = number
  default     = 256
}

variable "service_memory" {
  description = "Memory in MB for each ECS task"
  type        = number
  default     = 512
}

variable "container_port" {
  description = "Container application port"
  type        = number
  default     = 8080
}

variable "container_image_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}
variable "repository_urls" {
  description = "ECR repository URLs by service name"
  type        = map(string)
}
variable "subnet_ids" {
  description = "Private subnet IDs for ECS tasks"
  type        = list(string)
}

variable "security_group_id" {
  description = "Security group for ECS tasks"
  type        = string
}

variable "target_group_arns" {
  description = "ALB target group ARNs by service"
  type        = map(string)
}

variable "desired_count" {
  description = "Number of ECS tasks per service"
  type        = number
  default     = 1
}