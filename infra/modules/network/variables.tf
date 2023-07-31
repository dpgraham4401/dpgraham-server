# network inputs

variable "project" {
  description = "The GCP project to deploy to"
  type        = string
}

variable "region" {
  description = "The region to deploy to"
  type        = string
  default     = "us-east1"
}

variable "environment" {
  description = "The environment to deploy to"
  type        = string
  validation {
    condition     = contains(["development", "production"], var.environment)
    error_message = "Environment must be one of [dev, prod]"
  }
}

variable "subnets" {
  description = "The subnets to deploy to"
  type        = list(object({
    subnet_name   = string
    subnet_ip     = string
    subnet_region = string
  }))
  default = []
}