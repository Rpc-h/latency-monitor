resource "google_cloud_run_service" "default" {
  name     = "cloudrun-srv"
  location = var.google_region

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
      }
      labels = {
        "run.googleapis.com/startupProbeType" = "Custom"
      }
    }
    spec {
      timeout_seconds = 300
      containers {
        image = format("%s:%s", var.container_image, var.container_tag)
        name  = "latency"
        ports {
          container_port = 80
          name           = "http1" 
        }
        args = []
        command = []
        env {
          name = "LATENCY_MONITOR_RPC_SERVER_ONE_HOP_ADDRESS"
          value = var.rpc_server_one_hop_address
        }
        env {
          name = "LATENCY_MONITOR_RPC_SERVER_ZERO_HOP_ADDRESS"
          value = var.rpc_server_zero_hop_address
        }
        resources {
          requests = {}
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

        startup_probe {
          initial_delay_seconds = 65
          period_seconds = 60
          failure_threshold = 2
          timeout_seconds = 60
          http_get {
            path = "/metrics"
            port = 80
          }
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

