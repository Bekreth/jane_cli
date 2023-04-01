package main

import (
	"fmt"
	"os"

	"github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/logger"
	"gopkg.in/yaml.v3"
)

// TODO: Make this agnostic to OS
const defaultConfigLocation = "etc/config.yaml"

type Config struct {
	Logger logger.Config `yaml:"logger"`
	Client client.Config `yaml:"client"`
}

func parseConfig(configLocation string) (Config, error) {
	output := Config{}

	if configLocation == "" {
		configLocation = defaultConfigLocation
	}

	fileBytes, err := os.ReadFile(configLocation)
	if err != nil {
		return output, fmt.Errorf("Failed to read config file: %v", err)
	}
	return output, yaml.Unmarshal(fileBytes, &output)
}
