variable "name" {
  description = "Frontend resource name"
  type        = string
}

variable "price_class" {
  description = "CloudFront price class"
  type        = string
  default     = "PriceClass_100"
}

variable "tags" {
  description = "Tags applied to frontend resources"
  type        = map(string)
  default     = {}
}