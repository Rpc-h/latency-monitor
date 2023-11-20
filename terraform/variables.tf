variable "google_project" {
  type        = string
  description = "The ID of the GCP project"
}

variable "google_regions" {
  description = "The GCP region"
  type        = list(string)
  default     = ["europe-west4", "europe-west3"]
}