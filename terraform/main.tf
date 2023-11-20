resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

module "cloud-run" {
  count           = length(var.google_regions)
  source          = "./modules/cloud-run"
  google_project  = var.google_project
  google_region   = var.google_regions[count.index]
}


# https://cloud.google.com/compute/docs/regions-zones