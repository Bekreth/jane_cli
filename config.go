package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/logger"
	"gopkg.in/yaml.v3"
)

const defaultConfigLocation = "etc/config.yaml"

type Config struct {
	Logger logger.Config `yaml:"logger"`
	Client client.Config `yaml:"client"`
}

func parseConfig() (Config, error) {
	output := Config{}

	configLocation := ""
	if len(os.Args) == 3 && os.Args[1] == "-c" {
		configLocation = os.Args[2]
	}

	if configLocation == "" {
		configLocation = defaultConfigLocation
	}

	path := filepath.FromSlash(configLocation)
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return output, fmt.Errorf("Failed to read config file at %v: %v", path, err)
	}
	err = yaml.Unmarshal(fileBytes, &output)
	if err != nil {
		return output, fmt.Errorf("Failed to parse config from %v: %v", path, err)
	}

	if output.Client == client.DefaultConfig {
		return output, fmt.Errorf("Configuration at %v missing Jane Client details", path)
	}
	return output, nil
}
