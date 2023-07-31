# Terraform v1.5.4

locals {
  vpc_name = var.environment == "production" ? "${var.project}-vpc-prod" : "${var.project}-vpc-dev"
}

## Project level VPC
#resource "google_compute_network" "vpc" {
#  name    = local.vpc_name
#  project = var.project
#}

# enable API for serverless VPC access.
resource "google_project_service" "vpc_access_api" {
  project = var.project
  service = "vpcaccess.googleapis.com"
}

module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 7.1"

  project_id              = var.project
  network_name            = local.vpc_name
  routing_mode            = "GLOBAL"
  auto_create_subnetworks = true

  subnets = var.subnets
}