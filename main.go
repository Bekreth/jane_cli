package main

import (
	"fmt"
	"os"

	"github.com/Bekreth/jane_cli/app"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"

	"github.com/eiannone/keyboard"
)

// TODO: Make this agnostic to OS
const defaultConfigLocation = "etc/config.yaml"

func main() {
	logger := logger.NewLogrusLogger()
	logger.Infoln("Starting Jane CLI")
	configLocation := defaultConfigLocation
	logger.Infof("Loading config from %v", configLocation)
	fileBytes, err := os.ReadFile(configLocation)
	if err != nil {
		outputError := fmt.Errorf("Failed to read config file: %v", err)
		logger.Infoln(outputError)
		panic(err)
	}
	config, err := parseConfig(fileBytes)
	if err != nil {
		outputError := fmt.Errorf("failed to read config file: %v", err)
		logger.Infoln(outputError)
		panic(err)
	}

	if config.Debugger {
		logger.EnableDebugger()
	}

	if config.Client.UserFilePath == "" {
		panic("no path is provided for user file path")
	}

	user, err := domain.NewUser(logger, config.Client.UserFilePath)
	if err != nil {
		os.Exit(1)
	}

	user.PostCheck()

	application := app.NewApplication(logger, user)

	logger.Infoln("Application initialized, starting run loop")
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			logger.Infof("failed to read from console: %v", err)
		}
		if !application.HandleKeyinput(char, key) {
			break
		}
	}
}
