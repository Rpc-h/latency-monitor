# latency-monitor

This is a monitoring tool for collecting latency metrics for RPCh. It works by opening a new socket on `${MONITOR_METRICS_ADDRESS}${MONITOR_METRICS_PATH}` for Prometheus to scrape metrics from. Then it makes calls to `${MONITOR_RPC_SERVER_ADDRESS}` every `${MONITOR_METRICS_REQUEST_INTERVAL}` seconds. The default logging level is `${MONITOR_LOG_LEVEL}` and can range from -1 (trace) up to 5 (panic), logging can be disabled with 7 (https://github.com/rs/zerolog#leveled-logging).

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
MONITOR_METRICS_REQUEST_INTERVAL=3
MONITOR_LOG_LEVEL=1
```

## Metrics

Besides the default metrics, the custom exposed metrics are:

```
# HELP rpch_latencies_success
# TYPE rpch_latencies_success summary
rpch_latencies_success{quantile="0.5"} 0
rpch_latencies_success{quantile="0.7"} 0
rpch_latencies_success{quantile="0.9"} 0
rpch_latencies_success{quantile="0.99"} 0
rpch_latencies_success_sum 0
rpch_latencies_success_count 0

# HELP rpch_latencies_failure
# TYPE rpch_latencies_failure summary
rpch_latencies_failure{quantile="0.5"} 0
rpch_latencies_failure{quantile="0.7"} 0
rpch_latencies_failure{quantile="0.9"} 0
rpch_latencies_failure{quantile="0.99"} 0
rpch_latencies_failure_sum 0
rpch_latencies_failure_count 0
```