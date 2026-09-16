variable "name" {
  description = "DynamoDB table name"
  type        = string
}

variable "hash_key" {
  description = "DynamoDB partition key"
  type        = string
}

variable "hash_key_type" {
  description = "DynamoDB partition key type"
  type        = string

  validation {
    condition     = contains(["S", "N", "B"], var.hash_key_type)
    error_message = "hash_key_type must be S, N, or B."
  }
}

variable "tags" {
  description = "Tags for DynamoDB"
  type        = map(string)
  default     = {}
}