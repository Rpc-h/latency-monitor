resource "google_cloud_run_service" "default" {
  name     = "cloudrun-srv"
  location = var.google_region

  template {
    spec {
      containers {
        image = format("%s:%s", var.container_image, var.container_tag)
        name  = "latency"
        ports {
          container_port = 80
        }
        env {
          name = "LATENCY_MONITOR_RPC_SERVER_ONE_HOP_ADDRESS"
          value = var.rpc_server_one_hop_address
        }
        env {
          name = "LATENCY_MONITOR_RPC_SERVER_ZERO_HOPE_ADDRESS"
          value = var.rpc_server_zero_hop_address
        }
        resources {
          limits = {
            cpu = "200m"
            memory = "128Mi"
          }
        }
        liveness_probe {
          timeout_seconds = 20
          period_seconds = 120
          failure_threshold = 5
          http_get {
            path = "/metrics"
            port = 80
          }
        }
      }
    }
  }

}