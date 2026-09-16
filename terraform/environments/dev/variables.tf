variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "project_name" {
  type    = string
  default = "ecommerce-platform"
}

variable "environment" {
  type    = string
  default = "dev"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "availability_zones" {
  type = list(string)

  default = [
    "us-east-1a",
    "us-east-1b"
  ]
}

variable "enable_nat_gateway" {
  type    = bool
  default = false
}

variable "database_username" {
  description = "PostgreSQL username"
  type        = string
  sensitive   = true
}

variable "database_password" {
  description = "PostgreSQL password"
  type        = string
  sensitive   = true
}