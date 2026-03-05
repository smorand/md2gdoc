output "service_url" {
  description = "URL of the deployed Cloud Run MCP server"
  value       = google_cloud_run_v2_service.mcp_server.uri
}

output "artifact_registry" {
  description = "Artifact Registry repository path"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.repo.repository_id}"
}
