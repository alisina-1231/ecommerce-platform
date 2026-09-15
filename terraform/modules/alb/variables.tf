variable "project_name" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public subnet IDs"
  type        = list(string)
}

variable "security_group_id" {
  description = "ALB security group ID"
  type        = string
}

variable "services" {
  description = "Application services"
  type        = list(string)

  default = [
    "product",
    "cart",
    "checkout",
    "order",
    "payment"
  ]
}

variable "container_port" {
  description = "Container port"
  type        = number
  default     = 8080
}