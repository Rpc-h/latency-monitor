package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
	Objectives: map[float64]float64{
		0.5:  0,
		0.7:  0,
		0.9:  0,
		0.99: 0,
	},
})

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
	viper.SetDefault("METRICS_REQUEST_INTERVAL", 3)
}

func main() {
	setup()

	go func() {
		for {
			select {
			case <-time.Tick(time.Second * time.Duration(viper.GetInt("METRICS_REQUEST_INTERVAL"))):
				go func() {
					latency, err := getRawLatency()
					if err != nil {
						log.Error().Msg(err.Error())

						latenciesFailure.Observe(latency)

						return
					}

					log.Debug().Msg("success")

					latenciesSuccess.Observe(latency)
				}()
			}
		}
	}()

	http.Handle(viper.GetString("METRICS_PATH"), promhttp.Handler())

	log.Info().Msgf("Webserver listening on %s", viper.GetString("METRICS_ADDRESS"))

	err := http.ListenAndServe(viper.GetString("METRICS_ADDRESS"), nil)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func getRawLatency() (float64, error) {
	now := time.Now()

	body, err := json.Marshal(struct {
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
		Id: fmt.Sprintf("%v", rand.Int()),
	})
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://primary.gnosis-chain.rpc.hoprtech.net", bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Origin", "http://primary.goerli.rpc.hoprtech.net")

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	log.Debug().Msg(fmt.Sprintf("%s", body))

	return time.Since(now).Seconds(), err
}
