package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gochain/web3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var minLatency int64
var maxLatency int64

var minLatencyMetric = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "rpch", Subsystem: "latency", Name: "min"})
var maxLatencyMetric = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "rpch", Subsystem: "latency", Name: "max"})
var requestSuccessMetric = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "rpch", Subsystem: "requests", Name: "200"})
var requestFailure400Metric = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "rpch", Subsystem: "requests", Name: "400"})
var requestFailure408Metric = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "rpch", Subsystem: "requests", Name: "408"})

type RPCH struct {
	Client  web3.Client
	Channel chan bool
}

func setupEnv() {
	viper.SetEnvPrefix("RPCH_MONITOR")

	var err error

	err = viper.BindEnv("RPC_SERVER_ADDRESS")
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("RPC_SERVER_ADDRESS", "http://localhost:8080/?exit-provider=https://primary.gnosis-chain.rpc.hoprtech.net")

	err = viper.BindEnv("METRICS_ADDRESS")
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("METRICS_ADDRESS", "localhost:1234")

	err = viper.BindEnv("METRICS_PATH")
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("METRICS_PATH", "/metrics")

	_ = viper.BindEnv("METRICS_REQUEST_TIMEOUT")
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("METRICS_REQUEST_TIMEOUT", 5)

	err = viper.BindEnv("METRICS_RESET_TIMEOUT")
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("METRICS_RESET_TIMEOUT", 30)
}

func main() {
	setupEnv()

	fmt.Println(viper.GetString("RPC_SERVER_ADDRESS"))
	fmt.Println(viper.GetString("METRICS_ADDRESS"))
	fmt.Println(viper.GetString("METRICS_PATH"))
	fmt.Println(viper.GetString("METRICS_REQUEST_TIMEOUT"))
	fmt.Println(viper.GetString("METRICS_RESET_TIMEOUT"))
	fmt.Println(viper.GetString("RPC_SERVER_ADDRESS"))

	client, err := web3.Dial(viper.GetString("RPC_SERVER_ADDRESS"))
	if err != nil {
		log.Panic(err)
	}

	rpch := &RPCH{
		Client:  client,
		Channel: make(chan bool),
	}

	go func() {
		for {
			select {
			case <-time.Tick(time.Duration(viper.GetInt("METRICS_RESET_TIMEOUT"))):
				atomic.StoreInt64(&minLatency, 0)
				atomic.StoreInt64(&maxLatency, 0)

				minLatencyMetric.Set(0)
				maxLatencyMetric.Set(0)
				requestSuccessMetric.Set(0)
				requestFailure400Metric.Set(0)
				requestFailure408Metric.Set(0)
			}
		}
	}()

	go func() {
		for {
			go func() {
				latency, err := rpch.getRawLatency()
				if err != nil {
					requestFailure400Metric.Inc()

					return
				}

				requestSuccessMetric.Inc()

				if latency < minLatency || minLatency == 0 {
					minLatency = latency
				}

				if latency > maxLatency {
					maxLatency = latency
				}

				minLatencyMetric.Set(float64(minLatency))
				maxLatencyMetric.Set(float64(maxLatency))

				rpch.Channel <- true
			}()

			select {
			case <-time.Tick(time.Duration(viper.GetInt("METRICS_REQUEST_TIMEOUT"))):
				requestFailure408Metric.Inc()
			case <-rpch.Channel:
			}
		}
	}()

	http.Handle(viper.GetString("METRICS_PATH"), promhttp.Handler())

	err = http.ListenAndServe(viper.GetString("METRICS_ADDRESS"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (rpch *RPCH) getRawLatency() (int64, error) {
	now := time.Now()

	_, err := rpch.Client.GetBlockByNumber(context.Background(), nil, false)
	if err != nil {
		return 0, err
	}

	return time.Since(now).Milliseconds(), nil
}
