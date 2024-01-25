variable "google_project" {
  type        = string
  description = "The ID of the GCP project"
}

variable "environment" {
  type        = string
  description = "The name of the environment"
}

variable "google_regions" {
  type = map(object({
    name      = string
    latitude  = number
    longitude = number
    start_at  = number
    interval_duration = number
  }))
  description = "The GCP region"
}

variable "rpc_server_client_tokens" {
  description = "The tokens for discovery platform used by rpc server"
  type = map(string)
  sensitive = true
}

variable "latency_container_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}

variable "rpc_server_container_tag" {
  description = "Container image tag"
  type        = string
}

variable "rpc_server_one_hop_address" {
  description = "Rpc-h server address for one hop"
  type        = string
  default     = "https://rpc-server-one-hop.staging.rpch.tech/?provider=https://gnosis-provider.rpch.tech"
}

variable "rpc_server_zero_hop_address" {
  description = "Rpc-h server address for zero hop"
  type        = string
  default     = "https://rpc-server-zero-hop.staging.rpch.tech/?provider=https://gnosis-provider.rpch.tech"
}
