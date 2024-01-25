# https://cloud.google.com/compute/docs/regions-zones
# https://www.coordenadas-gps.com/
google_regions = {
  "europe-west2" = {
    name      = "London, England"
    latitude  = 51.5074456
    longitude = -0.1277653
    start_at  = 0
    interval_duration = 30
  }

  "europe-west6" = {
    name      = "Zurich, Switzerland"
    latitude  = 47.3668389
    longitude = 8.5339821
    start_at  = 4
    interval_duration = 30
  }

  "us-central1" = {
    name      = "Council Bluffs, Iowa"
    latitude  = 41.258841
    longitude = -95.8519484
    start_at  = 6
    interval_duration = 30
  }

  "southamerica-east1" = {
    name      = "Osasco, São Paulo, Brazil"
    latitude  = -23.5324859
    longitude = -46.7916801
    start_at  = 8
    interval_duration = 30
  }

  "australia-southeast1" = {
    name      = "Sydney, Australia"
    latitude  = -33.8698439
    longitude = 151.2082848
    start_at  = 10
    interval_duration = 30
  }

  "asia-east2" = {
    name      = "Hong Kong"
    latitude  = 22.2793278
    longitude = 114.1628131
    start_at  = 12
    interval_duration = 30
  }

  "me-central1" = {
    name      = "Doha, Qatar"
    latitude  = 25.2856329
    longitude = 51.5264162
    start_at  = 14
    interval_duration = 30
  }
}

environment = "production"
rpc_server_container_tag = "1.2.0"
