package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var latenciesSuccess = promauto.NewSummary(prometheus.SummaryOpts{
	Namespace: "rpch",
	Subsystem: "latencies",
	Name:      "success",
	Help:      "Successfull latency of 1 hop in milliseconds",
	Objectives: map[float64]float64{
		0.5:  0,
		0.7:  0,
		0.9:  0,
		0.99: 0,
	},
})

var latenciesFailure = promauto.NewSummary(prometheus.SummaryOpts{
	Namespace: "rpch",
	Subsystem: "latencies",
	Name:      "failure",
	Help:      "Failure latency of 1 hop in milliseconds",
	Objectives: map[float64]float64{
		0.5:  0,
		0.7:  0,
		0.9:  0,
		0.99: 0,
	},
})

var latenciesSuccessZeroHop = promauto.NewSummary(prometheus.SummaryOpts{
	Namespace: "rpch",
	Subsystem: "latencies",
	Name:      "success_zerohop",
	Help:      "Successfull latency of 0 hop in milliseconds",
	Objectives: map[float64]float64{
		0.5:  0,
		0.7:  0,
		0.9:  0,
		0.99: 0,
	},
})

var latenciesFailureZeroHop = promauto.NewSummary(prometheus.SummaryOpts{
	Namespace: "rpch",
	Subsystem: "latencies",
	Name:      "failure_zerohop",
	Help:      "Failure latency of 0 hop in milliseconds",
	Objectives: map[float64]float64{
		0.5:  0,
		0.7:  0,
		0.9:  0,
		0.99: 0,
	},
})

func setup() error {
	var err error
	viper.SetEnvPrefix("LATENCY_MONITOR")

	err = viper.BindEnv("LOG_LEVEL")
	if err != nil {
		viper.SetDefault("LOG_LEVEL", 0)
	}
	//Set time format and global log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(viper.GetInt("LATENCY_MONITOR_LOG_LEVEL")))

	viper.BindEnv("METRICS_ADDRESS")
	if !viper.IsSet("METRICS_ADDRESS") {
		viper.Set("METRICS_ADDRESS", "0.0.0.0:80")
	}

	viper.BindEnv("METRICS_PATH")
	if !viper.IsSet("METRICS_PATH") {
		viper.Set("METRICS_PATH", "/metrics")
	}

	viper.BindEnv("REQUEST_INTERVAL_DURATION")
	if !viper.IsSet("REQUEST_INTERVAL_DURATION") {
		viper.Set("REQUEST_INTERVAL_DURATION", 2)
	}

	viper.BindEnv("RPC_SERVER_ONE_HOP_START")
	if !viper.IsSet("RPC_SERVER_ONE_HOP_START") {
		viper.Set("RPC_SERVER_ONE_HOP_START", 0)
	}

	viper.BindEnv("RPC_SERVER_ONE_HOP_ADDRESS")
	if !viper.IsSet("RPC_SERVER_ONE_HOP_ADDRESS") {
		err = errors.New("Environment variable \"LATENCY_MONITOR_RPC_SERVER_ONE_HOP_ADDRESS\" is not set ")
	}

	viper.BindEnv("RPC_SERVER_ZERO_HOP_START")
	if !viper.IsSet("RPC_SERVER_ZERO_HOP_START") {
		viper.Set("RPC_SERVER_ZERO_HOP_START", 1)
	}

	viper.BindEnv("RPC_SERVER_ZERO_HOP_ADDRESS")
	if !viper.IsSet("RPC_SERVER_ZERO_HOP_ADDRESS") {
		err = errors.New("Environment variable \"LATENCY_MONITOR_RPC_SERVER_ZERO_HOP_ADDRESS\" is not set ")
	}

	return err
}

func main() {
	var err error
	err = setup()
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	// started := time.Now().Unix()
	// client := http.Client{}
	log.Info().Msgf("Metrics listening on %s%s", viper.GetString("METRICS_ADDRESS"), viper.GetString("METRICS_PATH"))
	go func() {
		var interval = viper.GetInt32("REQUEST_INTERVAL_DURATION")
		go startLatencyMonitor(viper.GetString("RPC_SERVER_ZERO_HOP_ADDRESS"), interval, viper.GetInt64("RPC_SERVER_ZERO_HOP_START"), 0, latenciesSuccessZeroHop, latenciesFailureZeroHop)
		go startLatencyMonitor(viper.GetString("RPC_SERVER_ONE_HOP_ADDRESS"), interval, viper.GetInt64("RPC_SERVER_ONE_HOP_START"), 1, latenciesSuccess, latenciesFailure)
	}()

	http.Handle(viper.GetString("METRICS_PATH"), promhttp.Handler())
	err = http.ListenAndServe(viper.GetString("METRICS_ADDRESS"), nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

}

func startLatencyMonitor(server string, interval int32, start_at int64, hops int, successMetric prometheus.Summary, failureMetric prometheus.Summary) {
	client := http.Client{}
	sleep_time := start_at + 60 - (time.Now().Unix() % 60)
	//sleep_time = 0
	log.Info().Msgf("Starting monitor for %d hops in %d seconds", hops, sleep_time)
	time.Sleep(time.Duration(sleep_time) * time.Second)
	started := time.Now()
	log.Info().Msgf("Started at: %s", started)

	ticker := time.NewTicker(time.Second * time.Duration(interval))
	for t := range ticker.C {
		diff := t.Unix() - started.Unix()
		latency, err := getRawLatency(&client, server, diff)
		if err != nil {
			log.Err(err).Send()
			if latency != 0 {
				failureMetric.Observe(float64(latency))
			}
			return
		}
		log.Debug().Msgf("Successfully send %d hop message in %v", hops, time.Duration(latency)*time.Millisecond)
		successMetric.Observe(float64(latency))
	}
	ticker.Stop()
	return
}

func getRawLatency(client *http.Client, server string, id int64) (int64, error) {
	requestBody, err := json.Marshal(struct {
		Jsonrpc string   `json:"jsonrpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
		Id      string   `json:"id"`
	}{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockTransactionCountByNumber",
		Params: []string{
			"latest",
		},
		Id: fmt.Sprintf("%v", id),
	})
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodPost, server, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, err
	}

	now := time.Now()

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	latency := time.Since(now).Milliseconds()
	if response.StatusCode != 200 {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return 0, err
		}

		return latency, fmt.Errorf("%s", b)
	}

	return latency, nil
}
