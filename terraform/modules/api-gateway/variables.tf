variable "name" {
  description = "API Gateway HTTP API name"
  type        = string
}

variable "alb_listener_arn" {
  description = "Existing ALB HTTP listener ARN"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID used by the existing ALB"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for the API Gateway VPC Link"
  type        = list(string)
}

variable "security_group_ids" {
  description = "Security groups for the API Gateway VPC Link"
  type        = list(string)
}

variable "jwt_issuer" {
  description = "Cognito JWT issuer URL"
  type        = string
}

variable "jwt_audience" {
  description = "Cognito app client ID"
  type        = list(string)
}

variable "throttling_burst_limit" {
  description = "API Gateway throttling burst limit"
  type        = number
  default     = 100
}

variable "throttling_rate_limit" {
  description = "API Gateway throttling rate limit"
  type        = number
  default     = 50
}

variable "tags" {
  description = "Tags applied to API Gateway resources"
  type        = map(string)
  default     = {}
}