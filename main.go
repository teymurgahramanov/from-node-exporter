package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

// Target represents the structure of the YAML data
type Target struct {
	Type string `yaml:"type"`
	Address string `yaml:"address"`
	Interval int `yaml:"interval"`
}

// config represents the dynamic expansion of targets
type Config struct {
	Targets []map[string]Target `yaml:"targets"`
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

func probeTCP(address string) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {

		logger.Error(fmt.Sprint(err))
		return false
	}
	defer conn.Close()
	return true
}

func probeHTTP(address string) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// Make a GET request to the URL
	resp, err := http.Get(address)
	if err != nil {
		logger.Error(fmt.Sprint(err))
		return false
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}

func worker(target string, module string, address string, interval int) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logMessageProbeSuccess := fmt.Sprintf("%v is UP",target)
	logMessageProbeFailure := fmt.Sprintf("%v is DOWN",target)
	logMessageWrongModule := fmt.Sprintf("Wrong module for %v",target)
	
	switch module {
	case "tcp":
		for {
			result := probeTCP(address)
			if result {
				probeStatus.WithLabelValues(target, module).Set(1)
				logger.Info(logMessageProbeSuccess)
			} else {
				probeStatus.WithLabelValues(target, module).Set(0)
				logger.Info(logMessageProbeFailure)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	case "http":
		for {
			result := probeHTTP(address)
			if result {
				probeStatus.WithLabelValues(target, module).Set(1)
				logger.Info(logMessageProbeSuccess)
			} else {
				probeStatus.WithLabelValues(target, module).Set(0)
				logger.Info(logMessageProbeFailure)
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
	// Read YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	// Unmarshal YAML data into config
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logger.Error(fmt.Sprint(err))
	}

	promRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = promRegistry
	prometheus.DefaultGatherer = promRegistry
	prometheus.MustRegister(probeStatus)
	
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		http.ListenAndServe(":8080", nil)
	}()
	
	var wg sync.WaitGroup

	// Access and print values
	for _, targetMap := range config.Targets {
		for key, value := range targetMap {
			wg.Add(1) // Increment the wait group for each new goroutine
			go func(target string, module string, address string, interval int) {
				defer wg.Done()
				logger.Info(fmt.Sprintf("Starting probe %v on %v using module %v for every %v seconds",target,address,strings.ToUpper(module),interval))
				worker(target, module, address, interval)
			}(key, value.Type, value.Address, value.Interval)
		}
	}
	
	// Wait for all the goroutines to finish
	wg.Wait()
}