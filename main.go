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

type configuration struct {
	Targets []map[string]targetConfig `yaml:"targets"`
	Exporter exporterConfig `yaml:"exporter"`
}

type targetConfig struct {
	Address string `yaml:"address"`
	Module string `yaml:"module"`
	Interval int `yaml:"interval"`
	Timeout int `yaml:"timeout"`
}

type exporterConfig struct {
	MetricsListenPath string `yaml:"metricsListenPath"`
	MetricsListenPort int `yaml:"metricsListenPort"`
	DefaultProbeInterval int `yaml:"defaultProbeInterval"`
	DefaultProbeTimeout int `yaml:"defaultProbeTimeout"`
}

var (
	probeResult = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "probe_result",
			Help: "Current status of the probe (1 for success, 0 for failure)",
		},
		[]string{"target", "module","address"},
	)
)

func probe(target string, module string, address string, interval int, timeout int) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	
	switch module {
	case "tcp":
		for {
			result,error := modules.ProbeTCP(address,timeout)
			if result {
				probeResult.WithLabelValues(target, module, address).Set(1)
			} else {
				if error != nil {
					logger.Error(fmt.Sprintf(error.Error()),"target",target)
				}
				probeResult.WithLabelValues(target, module, address).Set(0)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	case "http":
		for {
			result,error := modules.ProbeHTTP(address,timeout)
			if result {
				probeResult.WithLabelValues(target, module, address).Set(1)
			} else {
				if error != nil {
					logger.Error(fmt.Sprintf(error.Error()),"target",target)
				}
				probeResult.WithLabelValues(target, module, address).Set(0)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	default:
		logger.Error("Wrong module","target",target)
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

	var config configuration
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
	if config.Exporter.DefaultProbeTimeout == 0 {
		config.Exporter.DefaultProbeTimeout = 5
	}

	promRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = promRegistry
	prometheus.DefaultGatherer = promRegistry
	prometheus.MustRegister(probeResult)
	
	go func() {
		http.Handle(config.Exporter.MetricsListenPath, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		http.ListenAndServe(":"+fmt.Sprint(config.Exporter.MetricsListenPort), nil)
	}()
	
	var wg sync.WaitGroup

	for _, entry := range config.Targets {
		for key, value := range entry {
			wg.Add(1)
			go func(target string, module string, address string, interval int, timeout int) {
				defer wg.Done()
				if interval == 0 {
					interval = config.Exporter.DefaultProbeInterval
				}
				if timeout == 0 {
					timeout = config.Exporter.DefaultProbeTimeout
				}
				logger.Info(fmt.Sprintf("Starting probe %v at address %v for every %v seconds using module %v with timeout of %v seconds.",target,address,interval,module,timeout))
				probe(target, module, address, interval, timeout)
			}(key, value.Module, value.Address, value.Interval, value.Timeout)
		}
	}
	
	wg.Wait()
}