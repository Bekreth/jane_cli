package main

import (
	"github.com/Bekreth/jane_cli/client"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Debugger bool          `yaml:"debugger"`
	Client   client.Config `yaml:"client"`
}

func parseConfig(fileBytes []byte) (Config, error) {
	output := Config{}
	err := yaml.Unmarshal(fileBytes, &output)
	return output, err
}
