package main

import "gopkg.in/yaml.v3"

type Config struct {
	Debugger bool `yaml:"debugger"`
}

func parseConfig(fileBytes []byte) (Config, error) {
	output := Config{}
	err := yaml.Unmarshal(fileBytes, &output)
	return output, err
}
