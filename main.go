package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teymurgahramanov/from-node-exporter/modules"
	"gopkg.in/yaml.v3"
)

type Target struct {
	Type string `yaml:"type"`
	Address string `yaml:"address"`
	Interval int `yaml:"interval"`
}

var (
	probeStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "probe_status",
			Help: "Current status of the probe (1 for success, 0 for failure)",
		},
		[]string{"target", "module"},
	)
)

func worker(target string, module string, address string, interval int) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logMessageProbeSuccess := fmt.Sprintf("%v is UP",target)
	logMessageProbeFailure := fmt.Sprintf("%v is DOWN",target)
	logMessageWrongModule := fmt.Sprintf("Wrong module for %v",target)
	
	switch module {
	case "tcp":
		for {
			result,error := modules.ProbeTCP(address)
			if result {
				probeStatus.WithLabelValues(target, module).Set(1)
				logger.Info(logMessageProbeSuccess)
			} else {
				probeStatus.WithLabelValues(target, module).Set(0)
				if error != nil {
					logger.Error(fmt.Sprintf(error.Error()))
				} else {
					logger.Info(logMessageProbeFailure)
				}
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	case "http":
		for {
			result,error := modules.ProbeHTTP(address)
			if result {
				probeStatus.WithLabelValues(target, module).Set(1)
				logger.Info(logMessageProbeSuccess)
			} else {
				probeStatus.WithLabelValues(target, module).Set(0)
				if error != nil {
					logger.Error(fmt.Sprintf(error.Error()))
					return false
				} else {
					logger.Info(logMessageProbeFailure)
					return false
				}
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	default:
		logger.Error(logMessageWrongModule)
		return false
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	metricsListenPort := os.Getenv("METRICS_LISTEN_PORT")
	if metricsListenPort == "" {
		metricsListenPort = "8080"
	}
	metricsListenPath := os.Getenv("METRICS_LISTEN_PATH")
	if metricsListenPath == "" {
		metricsListenPath = "/metrics"
	}
	targetsFileName := "targets.yaml"

	data, err := os.ReadFile(targetsFileName)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	var targetEntries []map[string]Target
	err = yaml.Unmarshal(data, &targetEntries)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	promRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = promRegistry
	prometheus.DefaultGatherer = promRegistry
	prometheus.MustRegister(probeStatus)
	
	go func() {
		http.Handle(metricsListenPath, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		http.ListenAndServe(":"+metricsListenPort, nil)
	}()
	
	var wg sync.WaitGroup

	for _, entry := range targetEntries {
		for key, value := range entry {
			wg.Add(1)
			go func(target string, module string, address string, interval int) {
				defer wg.Done()
				logger.Info(fmt.Sprintf("Starting probe %v on %v using module %v for every %v seconds",target,address,strings.ToUpper(module),interval))
				if interval == 0 {
					interval = 22
				}
				worker(target, module, address, interval)
			}(key, value.Type, value.Address, value.Interval)
		}
	}
	
	wg.Wait()
}