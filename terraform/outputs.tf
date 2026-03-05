output "service_url" {
  description = "Internal Cloud Run URL"
  value       = google_cloud_run_v2_service.mcp_server.uri
}

output "public_url" {
  description = "Public HTTPS URL of the MCP server"
  value       = "https://${var.domain}/mcp"
}

output "load_balancer_ip" {
  description = "Global static IP - create a DNS A record pointing your domain to this IP"
  value       = google_compute_global_address.default.address
}

output "artifact_registry" {
  description = "Artifact Registry repository path"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.repo.repository_id}"
}

output "dns_instructions" {
  description = "DNS configuration instructions"
  value       = "Create an A record for ${var.domain} pointing to ${google_compute_global_address.default.address}. SSL certificate provisioning may take up to 60 minutes after DNS propagation."
}
