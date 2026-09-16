variable "name" {
  description = "RDS instance identifier"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for the RDS subnet group"
  type        = list(string)
}

variable "security_group_ids" {
  description = "Security groups attached to RDS"
  type        = list(string)
}

variable "instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "Initial storage in GB"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "Maximum storage in GB"
  type        = number
  default     = 50
}

variable "database_name" {
  description = "Initial PostgreSQL database"
  type        = string
  default     = "productdb"
}

variable "username" {
  description = "PostgreSQL username"
  type        = string
}

variable "password" {
  description = "PostgreSQL password"
  type        = string
  sensitive   = true
}

variable "backup_retention_period" {
  description = "Number of days to retain backups"
  type        = number
  default     = 1
}

variable "tags" {
  description = "Tags for RDS resources"
  type        = map(string)
  default     = {}
}