variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region for Cloud Run deployment"
  type        = string
  default     = "europe-west1"
}

variable "service_name" {
  description = "Cloud Run service name"
  type        = string
  default     = "md2gdoc"
}

variable "image" {
  description = "Container image to deploy (e.g. gcr.io/PROJECT/md2gdoc:latest)"
  type        = string
}
