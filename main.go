package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Latencies struct {
	Segments int `json:"segDur"`
	RpcCall  int `json:"rpcDur"`
	ExitNode int `json:"exitNodeDur"`
	Hopr     int `json:"hoprDur"`
}

type ResponseBody struct {
    Lats Latencies `json:"stats"`
}

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

	viper.BindEnv("LOCATION_LATITUDE")
	if !viper.IsSet("LOCATION_LATITUDE") {
		err = errors.New("Environment variable \"LATENCY_MONITOR_LOCATION_LATITUDE\" is not set ")
	}

	viper.BindEnv("LOCATION_LONGITUDE")
	if !viper.IsSet("LOCATION_LONGITUDE") {
		err = errors.New("Environment variable \"LATENCY_MONITOR_LOCATION_LONGITUDE\" is not set ")
	}

	viper.BindEnv("LOCATION_NAME")
	if !viper.IsSet("LOCATION_NAME") {
		err = errors.New("Environment variable \"LATENCY_MONITOR_LOCATION_NAME\" is not set ")
	}

	viper.BindEnv("LOCATION_REGION")
	if !viper.IsSet("LOCATION_REGION") {
		err = errors.New("Environment variable \"LATENCY_MONITOR_LOCATION_REGION\" is not set ")
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

	log.Info().Msgf("Metrics listening on %s%s", viper.GetString("METRICS_ADDRESS"), viper.GetString("METRICS_PATH"))
	go func() {
		var interval = viper.GetInt32("REQUEST_INTERVAL_DURATION")
		go startLatencyMonitor(viper.GetString("RPC_SERVER_ZERO_HOP_ADDRESS"), interval, viper.GetInt64("RPC_SERVER_ZERO_HOP_START"), "0")
		go startLatencyMonitor(viper.GetString("RPC_SERVER_ONE_HOP_ADDRESS"), interval, viper.GetInt64("RPC_SERVER_ONE_HOP_START"), "1")
	}()

	http.Handle(viper.GetString("METRICS_PATH"), promhttp.Handler())
	err = http.ListenAndServe(viper.GetString("METRICS_ADDRESS"), nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

}

func startLatencyMonitor(server string, interval int32, start_at int64, hops string) {
	var successMetric = promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "rpch",
		Subsystem: "latencies",
		Name:      "success",
		ConstLabels: map[string]string{
			"location":  viper.GetString("LOCATION_NAME"),
			"region":    viper.GetString("LOCATION_REGION"),
			"latitude":  viper.GetString("LOCATION_LATITUDE"),
			"longitude": viper.GetString("LOCATION_LONGITUDE"),
			"hops":      hops,
		},
		Help: "Successfull latency of 1 hop in milliseconds",
		Objectives: map[float64]float64{
			0.5:  0,
			0.7:  0,
			0.9:  0,
			0.99: 0,
		},
	})

	var failureMetric = promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "rpch",
		Subsystem: "latencies",
		Name:      "failure",
		ConstLabels: map[string]string{
			"location":  viper.GetString("LOCATION_NAME"),
			"region":    viper.GetString("LOCATION_REGION"),
			"latitude":  viper.GetString("LOCATION_LATITUDE"),
			"longitude": viper.GetString("LOCATION_LONGITUDE"),
			"hops":      hops,
		},
		Help: "Failure latency of 1 hop in milliseconds",
		Objectives: map[float64]float64{
			0.5:  0,
			0.7:  0,
			0.9:  0,
			0.99: 0,
		},
	})

	var errorMetric = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rpch",
		Subsystem: "latencies",
		Name:      "error",
		ConstLabels: map[string]string{
			"location":  viper.GetString("LOCATION_NAME"),
			"region":    viper.GetString("LOCATION_REGION"),
			"latitude":  viper.GetString("LOCATION_LATITUDE"),
			"longitude": viper.GetString("LOCATION_LONGITUDE"),
			"hops":      hops,
		},
		Help: "Error latency of 1 hop in milliseconds",
	})

	sleep_time := start_at + 60 - (time.Now().Unix() % 60)
	log.Info().Msgf("Starting monitor for %s hops in %d seconds", hops, sleep_time)
	time.Sleep(time.Duration(sleep_time) * time.Second)
	started := time.Now()
	log.Info().Msgf("Started at: %s", started)

	ticker := time.NewTicker(time.Second * time.Duration(interval))
	for t := range ticker.C {
		go func(tCopy time.Time) {
			diff := tCopy.Unix() - started.Unix()
			latency, lats, err := getRawLatency(server, diff)
			fmt.Println("lats", lats);
			if err != nil {
				log.Err(err).Send()
				if latency == 0 { // Assign largest response time
					log.Warn().Msgf("Latency not reported, a major error exist with rpc-server")
					errorMetric.Add(float64(latency))
				} else {
					log.Err(fmt.Errorf("Failed to send %s hop message in %v", hops, time.Duration(latency)*time.Millisecond))
					failureMetric.Observe(float64(latency))
				}
			} else {
				log.Debug().Msgf("Successfully send %s hop message in %v", hops, time.Duration(latency)*time.Millisecond)
				successMetric.Observe(float64(latency))
			}
		}(t)
	}

	ticker.Stop()
	log.Err(fmt.Errorf("Latency monitor %s hop stopped working", hops))
	os.Exit(3)

}

func getRawLatency(server string, id int64) (int64, *Latencies, error) {
	client := http.Client{Timeout: 60 * time.Second}
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
		return 0, nil, err
	}

	request, err := http.NewRequest(http.MethodPost, server, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, nil, err
	}

	now := time.Now()

	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}

	latency := time.Since(now).Milliseconds()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	if response.StatusCode != 200 {
		return latency, nil, fmt.Errorf("%s", body)
	}

    var payload ResponseBody
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return 0, nil, err
	}
	return latency, &payload.Lats, nil
}
