resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

module "cloud-run" {
  source                   = "./modules/cloud-run"
  google_project           = var.google_project
  google_regions           = var.google_regions
  rpc_server_client_tokens = var.rpc_server_client_tokens
  environment              = var.environment
  latency_container_tag    = "69d1c571dd6d5e0d07b76d5fe18648c5b6d8db0e9040fa915b172d9a92abce09"
  rpc_server_container_tag = var.rpc_server_container_tag
  depends_on               = [google_project_service.run]
}
