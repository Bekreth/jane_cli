package main

import (
	"fmt"
	"os"

	"github.com/Bekreth/jane_cli/app"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"

	"github.com/eiannone/keyboard"
)

func main() {
	fmt.Println("Starting Jane CLI")
	config, err := parseConfig("")
	logger, err := logger.NewLogrusLogger(config.Logger)

	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	if config.Client.UserFilePath == "" {
		fmt.Println("no path is provided for user file path")
		os.Exit(1)
	}

	user, err := domain.NewUser(
		logger.AddContext("service", "userReader"),
		config.Client.UserFilePath,
	)
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
