package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teymurgahramanov/from-node-exporter/modules"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Targets []map[string]TargetConfig `yaml:"targets"`
	Exporter ExporterConfig `yaml:"exporter"`
}

type TargetConfig struct {
	Address string `yaml:"address"`
	Module string `yaml:"module"`
	Interval int `yaml:"interval"`
}

type ExporterConfig struct {
	MetricsListenPath string `yaml:"metricsListenPath"`
	MetricsListenPort int `yaml:"metricsListenPort"`
	DefaultProbeInterval int `yaml:"defaultProbeInterval"`
}

var (
	probeStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "probe_success",
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
				} else {
					logger.Info(logMessageProbeFailure)
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
	targetsFileName := "config.yaml"

	data, err := os.ReadFile(targetsFileName)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	if config.Exporter.MetricsListenPort == 0 {
		config.Exporter.MetricsListenPort = 8080
	}
	if config.Exporter.MetricsListenPath == "" {
		config.Exporter.MetricsListenPath = "/metrics"
	}
	if config.Exporter.DefaultProbeInterval == 0 {
		config.Exporter.DefaultProbeInterval = 22
	}
	
	promRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = promRegistry
	prometheus.DefaultGatherer = promRegistry
	prometheus.MustRegister(probeStatus)
	
	go func() {
		http.Handle(config.Exporter.MetricsListenPath, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		http.ListenAndServe(":"+fmt.Sprint(config.Exporter.MetricsListenPort), nil)
	}()
	
	var wg sync.WaitGroup

	for _, entry := range config.Targets {
		for key, value := range entry {
			wg.Add(1)
			go func(target string, module string, address string, interval int) {
				defer wg.Done()
				if interval == 0 {
					interval = config.Exporter.DefaultProbeInterval
				}
				logger.Info(fmt.Sprintf("Starting probe %v at address %v for every %v seconds using module %v",target,address,interval,module))
				worker(target, module, address, interval)
			}(key, value.Module, value.Address, value.Interval)
		}
	}
	
	wg.Wait()
}