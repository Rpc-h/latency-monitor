resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

module "cloud-run" {
  count           = length(var.google_regions)
  source          = "./modules/cloud-run"
  google_project  = var.google_project
  google_region   = var.google_regions[count.index]
  depends_on      = [ google_project_service.run ]
  container_tag   = "0b20a4f"
}


# https://cloud.google.com/compute/docs/regions-zones