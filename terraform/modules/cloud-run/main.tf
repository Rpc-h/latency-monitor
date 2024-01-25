resource "google_service_account" "cloud_run" {
  #Not inferred from the provider
  project      = var.google_project
  account_id   = "latency-monitor-${var.environment}"
  display_name = "latency-monitor-${var.environment}"
}

resource "google_project_iam_binding" "cloud_run" {
  members = ["serviceAccount:${google_service_account.cloud_run.email}"]
  role    = "roles/secretmanager.secretAccessor"
  #Not inferred from the provider
  project = var.google_project
}

module "cloud-run-service" {
  for_each                  = var.google_regions
  source                    = "./cloud-run-service"
  google_project            = var.google_project
  environment               = var.environment
  google_region             = each.key
  client_token              = var.rpc_server_client_tokens[each.key]
  latency_container_tag     = var.latency_container_tag
  latency_start_at          = each.value.start_at
  latency_interval_duration = each.value.interval_duration
  rpc_server_container_tag  = var.rpc_server_container_tag
  service_account_email     = google_service_account.cloud_run.email
  location_name             = each.value.name
  location_latitude         = each.value.latitude
  location_longitude        = each.value.longitude
  depends_on                = [google_project_iam_binding.cloud_run]
}
