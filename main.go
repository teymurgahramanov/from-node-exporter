package main

import (
	"fmt"
	"os"

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

func worker(address string, probeType string) {

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
		}
	}
}
