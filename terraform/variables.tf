variable "google_project" {
  type        = string
  description = "The ID of the GCP project"
}

# https://cloud.google.com/compute/docs/regions-zones
# https://www.coordenadas-gps.com/
variable "google_regions" {
  description = "The GCP region details"
  type = map(object({
    name      = string
    latitude  = number
    longitude = number
    start_at  = number
  }))
}

variable "rpc_server_client_tokens" {
  description = "The tokens for discovery platform used by rpc server"
  type = map(string)
  sensitive = true
}

variable "environment" {
  type        = string
  description = "The name of the environment"
}

variable "rpc_server_container_tag" {
  type        = string
  description = "RPC Server container tag"
}
