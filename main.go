package main

import (
	"fmt"
	"os"

	"github.com/Bekreth/jane_cli/app"
	"github.com/Bekreth/jane_cli/cache"
	"github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"

	"github.com/eiannone/keyboard"
)

func main() {
	fmt.Println("Starting Jane CLI")
	config, err := parseConfig()

	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logger.NewLogrusLogger(config.Logger)
	if err != nil {
		fmt.Printf("failed to setup logger: %v\n", err)
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
	//TODO: This sucks, depointer it
	thisUser := &user

	client, err := client.NewClient(
		logger.AddContext("service", "httpClient"),
		config.Client,
		&thisUser.Auth,
		thisUser.SaveUserFile,
	)

	if err != nil {
		logger.Infof("failed to build Jane client: %v", err)
		fmt.Printf("failed to build Jane client: %v", err)
		os.Exit(1)
	}

	cache, err := cache.NewCache(logger.AddContext("service", "cache"), client, client)

	if err != nil {
		logger.Infof("failed to build patient cache: %v", err)
		fmt.Printf("failed to build patient cache: %v", err)
		os.Exit(1)
	}

	application := app.NewApplication(logger, thisUser, client, cache)

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
