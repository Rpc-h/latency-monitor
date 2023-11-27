variable "google_project" {
  type        = string
  description = "The ID of the GCP project"
}

variable "environment" {
  type        = string
  description = "The name of the environment"
}

variable "google_region" {
  type        = string
  description = "The GCP region"
}

variable "client_token" {
  type        = string
  description = "The rpc client token"
  sensitive = true
}

variable "service_account_email" {
  type        = string
  description = "The service account email of the Cloud Run"
}

variable "latency_container_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}

variable "rpc_server_container_tag" {
  description = "Container image tag"
  type        = string
  default     = "0.11.5"
}

variable "location_name" {
  type        = string
  description = "The location name of the GCP region"
}

variable "location_latitude" {
  type        = string
  description = "The location latitude of the GCP region"
}

variable "location_longitude" {
  type        = string
  description = "The location longitude of the GCP region"
}

variable "latency_start_at" {
  type        = number
  description = "The time at which the latency monitor should start"
}

variable "latency_interval_duration" {
  type        = number
  description = "The interval duration between requests"
}
