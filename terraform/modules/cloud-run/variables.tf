variable "google_project" {
  type        = string
  description = "The ID of the GCP project"
}

variable "google_region" {
  type        = string
  description = "The GCP region"
}

variable "container_image" {
  description = "Container image name"
  type        = string
  default     = "europe-west6-docker.pkg.dev/rpch-375921/rpch/latency-monitor"
}

variable "container_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}

variable "rpc_server_one_hop_address" {
  description = "Rpc-h server address for one hop"
  type        = string
  default     = "https://rpc-server-one-hop.rpch.tech"
}

variable "rpc_server_zero_hop_address" {
  description = "Rpc-h server address for zero hop"
  type        = string
  default     = "https://rpc-server-zero-hop.rpch.tech"
}
