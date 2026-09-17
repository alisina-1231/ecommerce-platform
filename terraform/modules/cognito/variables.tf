variable "name" {
  description = "Cognito User Pool name"
  type        = string
}

variable "client_name" {
  description = "Cognito app client name"
  type        = string
}

variable "tags" {
  description = "Tags applied to Cognito resources"
  type        = map(string)
  default     = {}
}