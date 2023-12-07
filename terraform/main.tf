resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

module "cloud-run" {
  source                   = "./modules/cloud-run"
  google_project           = var.google_project
  google_regions           = var.google_regions
  rpc_server_client_tokens = var.rpc_server_client_tokens
  environment              = var.environment
  latency_container_tag    = "004d54b64b5e06116088bd08bf382875da06fab3632ed2d3b58686e80781fd40"
  rpc_server_container_tag = "0.13.1"
  depends_on               = [google_project_service.run]
}
