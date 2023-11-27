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
  rpc_server_container_tag = "30e2f62fc6f9ce171e267fee6eeb4f8643946608d34f7c59747eb49f191fc44b"
  depends_on               = [google_project_service.run]
}
