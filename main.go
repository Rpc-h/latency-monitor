package main

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gochain/web3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	Client web3.Client
}

func setup() {
	viper.SetEnvPrefix("MONITOR")

	var err error

	err = viper.BindEnv("LOG_LEVEL")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("LOG_LEVEL", zerolog.InfoLevel)

	//Set time format and global log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(viper.GetInt("LOG_LEVEL")))

	err = viper.BindEnv("RPC_SERVER_ADDRESS")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("RPC_SERVER_ADDRESS", "http://localhost:8080/?exit-provider=https://primary.gnosis-chain.rpc.hoprtech.net")

	err = viper.BindEnv("METRICS_ADDRESS")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("METRICS_ADDRESS", "0.0.0.0:1234")

	err = viper.BindEnv("METRICS_PATH")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("METRICS_PATH", "/metrics")

	err = viper.BindEnv("METRICS_REQUEST_INTERVAL")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("METRICS_REQUEST_INTERVAL", 1)

	err = viper.BindEnv("METRICS_REQUEST_TIMEOUT")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("METRICS_REQUEST_TIMEOUT", 3)

	err = viper.BindEnv("METRICS_RESET_TIMEOUT")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	viper.SetDefault("METRICS_RESET_TIMEOUT", 30)
}

func main() {
	setup()

	client, err := web3.Dial(viper.GetString("RPC_SERVER_ADDRESS"))
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	rpch := &RPCH{
		Client: client,
	}

	go func() {
		for {
			select {
			case <-time.Tick(time.Second * time.Duration(viper.GetInt("METRICS_RESET_TIMEOUT"))):
				log.Debug().Msg("reset")

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

	timeout := time.Second * time.Duration(viper.GetInt("METRICS_REQUEST_TIMEOUT"))
	tickerTimeout := time.NewTicker(timeout)

	interval := time.Second * time.Duration(viper.GetInt("METRICS_REQUEST_INTERVAL"))
	tickerInterval := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-tickerTimeout.C:
				log.Error().Msg("timeout")

				requestFailure408Metric.Inc()

				//Reset the interval ticker
				tickerInterval.Reset(interval)
			case <-tickerInterval.C:
				latency, err := rpch.getRawLatency()
				if err != nil {
					log.Error().Msg(err.Error())

					requestFailure400Metric.Inc()

					//Reset the timout ticker
					tickerTimeout.Reset(timeout)

					continue
				}

				requestSuccessMetric.Inc()

				if latency < minLatency || minLatency == 0 {
					minLatency = latency
				}

				if latency > maxLatency {
					maxLatency = latency
				}

				//TODO - calculate average latency here

				minLatencyMetric.Set(float64(minLatency))
				maxLatencyMetric.Set(float64(maxLatency))

				log.Debug().Msg("success")

				//Reset the timout ticker
				tickerTimeout.Reset(timeout)
			}
		}
	}()

	http.Handle(viper.GetString("METRICS_PATH"), promhttp.Handler())

	log.Info().Msgf("Webserver listening on %s", viper.GetString("METRICS_ADDRESS"))

	err = http.ListenAndServe(viper.GetString("METRICS_ADDRESS"), nil)
	if err != nil {
		log.Fatal().Msg(err.Error())
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
