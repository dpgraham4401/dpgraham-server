# network inputs

variable "project" {
  description = "The GCP project to deploy to"
  type        = string
}

variable "environment" {
  description = "The environment to deploy to"
  type        = string
  validation {
    condition     = contains(["development", "production"], var.environment)
    error_message = "Environment must be one of [dev, prod]"
  }
}