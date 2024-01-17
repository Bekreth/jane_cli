package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/Bekreth/jane_cli/app"
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/cache"
	"github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	tsize "github.com/kopoli/go-terminal-size"

	"github.com/eiannone/keyboard"
)

var version string

func main() {
	fmt.Printf("Starting Jane CLI version %v\n", version)
	config, err := parseConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	// TODO: this needs to be configurable v.v.v.v.v.v.v
	loc, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		fmt.Printf("failed to set timezone: %v\n", err)
		os.Exit(1)
	}
	time.Local = loc

	logger, err := logger.NewLogrusLogger(config.Logger)
	if err != nil {
		fmt.Printf("failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	terminalSize, err := tsize.GetSize()
	if err != nil {
		logger.Infof("failed to get terminal size: %v", err)
		fmt.Printf("failed to get terminal size: %v", err)
		os.Exit(1)
	}
	terminal := terminal.NewTerminal(terminalSize, terminal.NewWindow())
	logger.Debugf("Terminal dimensions %v", terminalSize)

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
		thisUser,
		thisUser.SaveUserFile,
	)
	if err != nil {
		logger.Infof("failed to build Jane client: %v", err)
		fmt.Printf("failed to build Jane client: %v", err)
		os.Exit(1)
	}

	cache, err := cache.NewCache(
		logger.AddContext("service", "cache"),
		client,
	)
	if err != nil {
		logger.Infof("failed to build patient cache: %v", err)
		fmt.Printf("failed to build patient cache: %v", err)
		os.Exit(1)
	}

	application := app.NewApplication(logger, terminal, thisUser, client, cache)

	logger.Infoln("Application initialized, starting run loop")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Jane CLI has crashed ungracefully.  Notify author of this issue")
			logger.Infof("Crash: %v", r)
			for i, s := range strings.Split(string(debug.Stack()), "\n") {
				logger.Infof("%v: %v", i, s)
			}
		}
	}()
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
