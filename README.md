# latency-monitor

This is a monitoring tool for collecting latency metrics for RPCh. It works by opening a new socket on `${MONITOR_METRICS_ADDRESS}${MONITOR_METRICS_PATH}` for Prometheus to scrape metrics from. Then it makes calls to `${MONITOR_RPC_SERVER_ADDRESS}` or timeouts after `${MONITOR_METRICS_REQUEST_TIMEOUT}` seconds, whichever comes first. The Prometheus metrics are reset every `${MONITOR_METRICS_RESET_TIMEOUT}` seconds.

## Running

### RPC server

You have to provide address for RPC server in `${MONITOR_RPC_SERVER_ADDRESS}`. The easiest way to have a local RPC server is to follow the guide here: https://access.rpch.net/#section-2 , go to "RPCH DOCKER CONNECTOR", click on "Download", and then execute the `docker run` command provided.

### Latency monitor

To run the latency monitor, simply execute:

```shell
go run main.go
```

## Variables

These are the accepted environment variables and their default values:

```dotenv
MONITOR_RPC_SERVER_ADDRESS=http://localhost:8080/?exit-provider=https://primary.gnosis-chain.rpc.hoprtech.net
MONITOR_METRICS_ADDRESS=0.0.0.0:1234
MONITOR_METRICS_PATH=/metrics
MONITOR_METRICS_REQUEST_INTERVAL=5
MONITOR_METRICS_REQUEST_TIMEOUT=5
MONITOR_METRICS_RESET_TIMEOUT=30
```

## Metrics

Besides the default metrics, the custom exposed metrics are:

```
# HELP rpch_latency_max 
# TYPE rpch_latency_max gauge
rpch_latency_max 0
# HELP rpch_latency_min 
# TYPE rpch_latency_min gauge
rpch_latency_min 0
# HELP rpch_requests_200 
# TYPE rpch_requests_200 gauge
rpch_requests_200 0
# HELP rpch_requests_400 
# TYPE rpch_requests_400 gauge
rpch_requests_400 0
# HELP rpch_requests_408 
# TYPE rpch_requests_408 gauge
rpch_requests_408 0
```