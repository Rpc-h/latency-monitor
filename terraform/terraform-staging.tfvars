# https://cloud.google.com/compute/docs/regions-zones
# https://www.coordenadas-gps.com/
google_regions = {

  "europe-west6" = {
    name      = "Zurich, Switzerland"
    latitude  = 47.3668389
    longitude = 8.5339821
    start_at  = 4
  }

  "us-central1" = {
    name      = "Council Bluffs, Iowa"
    latitude  = 41.258841
    longitude = -95.8519484
    start_at  = 6
  }
}

environment = "staging"
rpc_server_container_tag = "1.1.3"
