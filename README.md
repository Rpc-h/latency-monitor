# latency-monitor

This is a monitoring tool for collecting latency metrics for RPCh. It works by opening a new socket on `${MONITOR_METRICS_ADDRESS}${MONITOR_METRICS_PATH}` for Prometheus to scrape metrics from. Then it makes calls to `${MONITOR_RPC_SERVER_ADDRESS}` or timeouts after `${MONITOR_METRICS_REQUEST_TIMEOUT}` seconds, whichever comes first. The Prometheus metrics are reset every `${MONITOR_METRICS_RESET_TIMEOUT}` seconds.

## Running

To run, simply execute:

```shell
go run main.go
```

## RPCh server

You have to provide address for RPC server in `${MONITOR_RPC_SERVER_ADDRESS}`. The easiest way to have a local RPC server is to follow the guide here: https://access.rpch.net/#section-2 , go to "RPCH DOCKER CONNECTOR", click on "Download", and then execute the `docker run` command provided.

## Variables

These are the variables and their default values:

```env
MONITOR_RPC_SERVER_ADDRESS=http://localhost:8080/?exit-provider=https://primary.gnosis-chain.rpc.hoprtech.net
MONITOR_METRICS_ADDRESS=0.0.0.0:1234
MONITOR_METRICS_PATH=/metrics
MONITOR_METRICS_REQUEST_TIMEOUT=5
MONITOR_METRICS_RESET_TIMEOUT=30
```