resource "google_project_service" "run" {
  service = "run.googleapis.com"
}

module "cloud-run" {
  source                   = "./modules/cloud-run"
  google_project           = var.google_project
  google_regions           = var.google_regions
  rpc_server_client_tokens = var.rpc_server_client_tokens
  environment              = var.environment
  latency_container_tag    = "1d033450ff611e61f63ea655e595105f046a1067715a27562ef9090739991ddc"
  rpc_server_container_tag = "9876d4d300e31a33ee5ef094c0f97eb489f6a034a2f5040361264ab25abb7f1b"
  depends_on               = [google_project_service.run]
}
