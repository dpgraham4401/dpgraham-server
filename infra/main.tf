terraform {
  backend "gcs" {
    bucket = "dpgraham-terraform-state"
    prefix = "terraform-prod"
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region
  zone    = "us-east1-b"
}

module "network" {
  source      = "./modules/network"
  project     = var.project
  environment = "production"
}

module "database" {
  source      = "./modules/sql"
  name        = var.db_name
  db_password = var.db_password
  db_username = var.db_username
  environment = "production"
  vpc         = module.network.vpc
}

module "load_balancer" {
  source           = "./modules/gcp-load-balancer"
  name             = "${var.project}-frontend"
  backend_service  = module.server-service.name
  frontend_service = module.frontend-service.name
  environment      = "production"
}


# The domain modules is used to provision resources, such as DNS zones and record sets for our domain
module "domain" {
  source       = "./modules/domain"
  project_id   = var.project
  domain_name  = var.domain_name
  ipv4_address = module.load_balancer.ip_address
}

resource "google_artifact_registry_repository" "dpgraham_com" {
  location      = var.region
  repository_id = var.repo_id
  description   = "Repository for dpgraham.com"
  format        = "DOCKER"
}

module "frontend-service" {
  source        = "./modules/cloud-run"
  name          = "${var.project}-frontend"
  image         = format("%s-docker.pkg.dev/%s/%s/%s:latest", google_artifact_registry_repository.dpgraham_com.location, var.project, google_artifact_registry_repository.dpgraham_com.repository_id, var.client_image_name)
  vpc_connector = module.database.vpc_connector
  port          = "3000"
  environment   = "production"
}
module "server-service" {
  source        = "./modules/cloud-run"
  name          = "${var.project}-server"
  image         = format("%s-docker.pkg.dev/%s/%s/%s:latest", google_artifact_registry_repository.dpgraham_com.location, var.project, google_artifact_registry_repository.dpgraham_com.repository_id, var.server_image_name)
  vpc_connector = module.database.vpc_connector
  port          = "8080"
  environment   = "production"
  env           = [
    {
      name  = "DB_PORT"
      value = "5432"
    },
    {
      name  = "DB_NAME"
      value = module.database.db_name
    },
    {
      name  = "DB_USER"
      value = module.database.db_user
    },
    {
      name  = "DB_PASSWORD"
      value = module.database.db_password
    },
    {
      name  = "DB_HOST"
      value = module.database.db_host
    }
  ]
}
