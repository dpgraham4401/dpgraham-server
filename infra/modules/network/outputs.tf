# network outputs

output "database_subnet_cidr" {
  value = google_compute_subnetwork.database_subnet.ip_cidr_range
}

output "network" {
  value = module.vpc.network
}