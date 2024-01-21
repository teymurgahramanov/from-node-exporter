package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

// Target represents the structure of the YAML data
type Target struct {
	Address string `yaml:"address"`
	Type    string `yaml:"type"`
}

// config represents the dynamic expansion of targets
type Config struct {
	Targets []map[string]Target `yaml:"targets"`
}

var (
	probeSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "probe_success_total",
			Help: "Total number of successful probes",
		},
		[]string{"target", "probeType"},
	)

	probeFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "probe_failure_total",
			Help: "Total number of failed probes",
		},
		[]string{"target", "probeType"},
	)
)

func init() {
	prometheus.MustRegister(probeSuccess)
	prometheus.MustRegister(probeFailure)
}

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

func dispatch(target string, address string, probeType string) bool {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	/*
		1. Check probeType and use apr. function
			1.1 Write metric
			1.2 Repeat fumction
	*/
	switch probeType {
	case "tcp":
		success := probeTCP(address)
		if success {
			probeSuccess.WithLabelValues(target, probeType).Inc()
		} else {
			probeFailure.WithLabelValues(target, probeType).Inc()
		}
		return success
	case "http":
		success := probeHTTP(address)
		if success {
			probeSuccess.WithLabelValues(target, probeType).Inc()
		} else {
			probeFailure.WithLabelValues(target, probeType).Inc()
		}
		return success
	default:
		logger.Error("Unknown probe type")
		return false
	}
}

func main() {
	// Read YAML file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	// Unmarshal YAML data into config
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	// Access and print values
	for _, targetMap := range config.Targets {
		for key, value := range targetMap {
			fmt.Printf("Target: %s\n", key)
			fmt.Printf("Address: %s\n", value.Address)
			fmt.Printf("Type: %s\n", value.Type)
			fmt.Printf("Result: %t\n", dispatch(key, value.Address, value.Type))
			http.Handle("/metrics", promhttp.Handler())
			go func() {
				logger.Error(http.ListenAndServe(":8080", nil))
			}()
			dispatch(key, value.Address, value.Type)
		}
	}
	func main() {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			logger.Error(http.ListenAndServe(":8080", nil))
		}()
	
		// Example usage
		dispatch("example", "localhost:8080", "http")
	}
}
