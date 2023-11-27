
resource "google_cloud_run_service" "latency_monitor" {
  name     = "latency-monitor-${var.environment}"
  project  = var.google_project
  location = var.google_region

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"          = "1"
        "run.googleapis.com/cpu-throttling"         = "false"
        "run.googleapis.com/container-dependencies" = jsonencode({ latency = ["rpc-server-zero-hop"] })
      }
      labels = {
        "run.googleapis.com/startupProbeType" = "Custom"
      }
    }
    spec {
      timeout_seconds       = 30
      service_account_name  = var.service_account_email
      container_concurrency = 1
      containers {
        image = length(var.latency_container_tag) == 64 ? format("europe-west6-docker.pkg.dev/rpch-375921/rpch/latency-monitor@sha256:%s", var.latency_container_tag) : format("europe-west6-docker.pkg.dev/rpch-375921/rpch/latency-monitor:%s", var.latency_container_tag)
        name  = "latency"
        ports {
          container_port = 80
          name           = "http1"
        }
        args    = []
        command = []
        env {
          name  = "LATENCY_MONITOR_REQUEST_INTERVAL_DURATION"
          value = var.latency_interval_duration
        }
        env {
          name  = "LATENCY_MONITOR_RPC_SERVER_ONE_HOP_ADDRESS"
          value = "http://localhost:8081/?provider=https://gnosis-provider.rpch.tech"
        }
        env {
          name  = "LATENCY_MONITOR_RPC_SERVER_ZERO_HOP_START"
          value = var.latency_start_at
        }
        env {
          name  = "LATENCY_MONITOR_RPC_SERVER_ZERO_HOP_ADDRESS"
          value = "http://localhost:8080/?provider=https://gnosis-provider.rpch.tech"
        }
        env {
          name  = "LATENCY_MONITOR_RPC_SERVER_ONE_HOP_START"
          value = var.latency_start_at + 1
        }
        env {
          name  = "LATENCY_MONITOR_LOCATION_REGION"
          value = var.google_region
        }
        env {
          name  = "LATENCY_MONITOR_LOCATION_NAME"
          value = var.location_name
        }
        env {
          name  = "LATENCY_MONITOR_LOCATION_LATITUDE"
          value = var.location_latitude
        }
        env {
          name  = "LATENCY_MONITOR_LOCATION_LONGITUDE"
          value = var.location_longitude
        }
        resources {
          requests = {
            cpu    = "1"
            memory = "512Mi"
          }
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }
        liveness_probe {
          timeout_seconds   = 10
          period_seconds    = 120
          failure_threshold = 3
          http_get {
            path = "/metrics"
            port = 80
          }
        }

        startup_probe {
          initial_delay_seconds = 5
          period_seconds        = 10
          failure_threshold     = 10
          timeout_seconds       = 10
          http_get {
            path = "/metrics"
            port = 80
          }
        }
      }
      containers {
        image   = length(var.rpc_server_container_tag) == 64 ? format("europe-west6-docker.pkg.dev/rpch-375921/rpch/rpc-server@sha256:%s", var.rpc_server_container_tag) : format("europe-west6-docker.pkg.dev/rpch-375921/rpch/rpc-server:%s", var.rpc_server_container_tag)
        name    = "rpc-server-zero-hop"
        args    = []
        command = []
        env {
          name  = "DEBUG"
          value = "rpch*"
        }
        env {
          name  = "PORT"
          value = "8080"
        }
        env {
          name  = "FRONTEND_HTTP_PORT"
          value = "45750"
        }
        env {
          name  = "FRONTEND_HTTPS_PORT"
          value = "45751"
        }
        env {
          name  = "RESPONSE_TIMEOUT"
          value = "60000"
        }
        env {
          name  = "DISCOVERY_PLATFORM_API_ENDPOINT"
          value = "https://discovery.${var.environment}.rpch.tech"
        }
        env {
          name  = "FORCE_ZERO_HOP"
          value = "true"
        }
        env {
          name = "CLIENT"
          value = var.client_token
        }
        resources {
          requests = {
            cpu    = "1"
            memory = "512Mi"
          }
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }
        startup_probe {
          initial_delay_seconds = 15
          period_seconds        = 15
          failure_threshold     = 8
          timeout_seconds       = 10
          tcp_socket {
            port = 45751
          }
        }
      }
      containers {
        image   = length(var.rpc_server_container_tag) == 64 ? format("europe-west6-docker.pkg.dev/rpch-375921/rpch/rpc-server@sha256:%s", var.rpc_server_container_tag) : format("europe-west6-docker.pkg.dev/rpch-375921/rpch/rpc-server:%s", var.rpc_server_container_tag)
        name    = "rpc-server-one-hop"
        args    = []
        command = []
        env {
          name  = "DEBUG"
          value = "rpch*"
        }
        env {
          name  = "PORT"
          value = "8081"
        }
        env {
          name  = "FRONTEND_HTTP_PORT"
          value = "45752"
        }
        env {
          name  = "FRONTEND_HTTPS_PORT"
          value = "45753"
        }
        env {
          name  = "RESPONSE_TIMEOUT"
          value = "60000"
        }
        env {
          name  = "DISCOVERY_PLATFORM_API_ENDPOINT"
          value = "https://discovery.staging.rpch.tech"
        }
        env {
          name = "CLIENT"
          value = var.client_token
        }
        resources {
          requests = {
            cpu    = "1"
            memory = "512Mi"
          }
          limits = {
            cpu    = "1"
            memory = "512Mi"
          }
        }
        startup_probe {
          initial_delay_seconds = 15
          period_seconds        = 15
          failure_threshold     = 8
          timeout_seconds       = 10
          tcp_socket {
            port = 45753
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

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.latency_monitor.location
  project  = google_cloud_run_service.latency_monitor.project
  service  = google_cloud_run_service.latency_monitor.name

  policy_data = data.google_iam_policy.noauth.policy_data
  depends_on  = [google_cloud_run_service.latency_monitor]
}
