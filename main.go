package main

import (
	"bytes"
	"encoding/json"
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

var latenciesSuccessZeroHop = promauto.NewSummary(prometheus.SummaryOpts{
	Namespace: "rpch",
	Subsystem: "latencies",
	Name:      "success-zerohop",
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
	Name:      "failure-zerohop",
	Objectives: map[float64]float64{
		0.5:  0,
		0.7:  0,
		0.9:  0,
		0.99: 0,
	},
})

func setup() error {
	var err error

	viper.SetEnvPrefix("MONITOR")

	err = viper.BindEnv("LOG_LEVEL")
	if err != nil {
		return err
	}

	//Set time format and global log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.Level(viper.GetInt("LOG_LEVEL")))

	err = viper.BindEnv("RPC_SERVER_ADDRESS")
	if err != nil {
		return err
	}

	err = viper.BindEnv("RPC_SERVER_ADDRESS_ZERO_HOP")
	if err != nil {
		return err
	}

	err = viper.BindEnv("METRICS_ADDRESS")
	if err != nil {
		return err
	}

	err = viper.BindEnv("METRICS_PATH")
	return err
}

func main() {
	var err error
	err = setup()
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	started := time.Now().Unix()
	client := http.Client{}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(1))
		for t := range ticker.C {
			diff := t.Unix() - started
			rem := diff % 2

			// alternative rpc server calls every other second
			if rem == 0 {

				// run default rpc
				go func() {
					server := viper.GetString("RPC_SERVER_ADDRESS")
					latency, err := getRawLatency(&client, server, diff)
					if err != nil {
						log.Err(err).Send()
						if latency != 0 {
							latenciesFailure.Observe(latency)
						}
						return
					}

					log.Debug().Msgf("Successful observation, latency: %v", latency)
					latenciesSuccess.Observe(latency)
				}()

			} else {

				// run zero hop rpc
				go func() {
					server := viper.GetString("RPC_SERVER_ADDRESS_ZERO_HOP")
					latency, err := getRawLatency(&client, server, diff)
					if err != nil {
						log.Err(err).Send()
						if latency != 0 {
							latenciesFailureZeroHop.Observe(latency)
						}
						return
					}

					log.Debug().Msgf("Successful observation, latency: %v", latency)
					latenciesSuccessZeroHop.Observe(latency)
				}()

			}
		}
		ticker.Stop()
	}()

	http.Handle(viper.GetString("METRICS_PATH"), promhttp.Handler())

	log.Info().Msgf("Webserver listening on %s", viper.GetString("METRICS_ADDRESS"))

	err = http.ListenAndServe(viper.GetString("METRICS_ADDRESS"), nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func getRawLatency(client *http.Client, server string, id int64) (float64, error) {
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

	latency := time.Since(now).Seconds()

	if response.StatusCode != 200 {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return 0, err
		}

		return latency, fmt.Errorf("%s", b)
	}

	return latency, nil
}
