package app

import (
	"github.com/Bekreth/jane_cli/app/auth"
	"github.com/Bekreth/jane_cli/app/booking"
	"github.com/Bekreth/jane_cli/app/initialize"
	"github.com/Bekreth/jane_cli/app/root"
	"github.com/Bekreth/jane_cli/app/schedule"
	"github.com/Bekreth/jane_cli/app/terminal"
	Cache "github.com/Bekreth/jane_cli/cache"
	Client "github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type Application struct {
	logger    logger.Logger
	writer    terminal.ScreenWriter
	state     terminal.State
	allStates []terminal.State
}

func NewApplication(
	logger logger.Logger,
	user *domain.User,
	client Client.Client,
	cache Cache.Cache,
) Application {

	rootState := root.NewState(
		logger.AddContext("state", "root"),
		terminal.NewScreenWriter("jane:"),
	)
	initState := initialize.NewState(
		logger.AddContext("state", "init"),
		terminal.NewScreenWriter("init:"),
		user,
		rootState,
	)
	authState := auth.NewState(
		logger.AddContext("state", "auth"),
		terminal.NewScreenWriter("auth:"),
		client,
		rootState,
	)

	scheduleState := schedule.NewState(
		logger.AddContext("state", "schedule"),
		terminal.NewScreenWriter("schedule:"),
		client,
		rootState,
	)

	fetcher := struct {
		Cache.Cache
		Client.Client
	}{
		Cache:  cache,
		Client: client,
	}
	bookingState := booking.NewState(
		logger.AddContext("state", "booking"),
		terminal.NewScreenWriter("booking:"),
		fetcher,
		rootState,
	)

	rootState.RegisterStates(map[string]terminal.State{
		initState.Name():     initState,
		authState.Name():     authState,
		scheduleState.Name(): scheduleState,
		bookingState.Name():  bookingState,
	})
	rootState.Initialize()

	return Application{
		logger: logger,
		writer: terminal.NewScreenWriter(""),
		state:  rootState,
		allStates: []terminal.State{
			rootState,
			initState,
			authState,
			scheduleState,
			bookingState,
		},
	}
}

func (app *Application) HandleKeyinput(character rune, key keyboard.Key) bool {

	switch key {
	case keyboard.KeyCtrlC:
		app.logger.Infoln("exiting the application")
		app.writer.NewLine()
		app.writer.WriteString("Shutting down Jane CLI")
		app.writer.NewLine()
		return false
	case keyboard.KeyCtrlU:
		app.state.ClearBuffer()
		return true
	}

	nextState := app.state.HandleKeyinput(character, key)
	if nextState.Name() != app.state.Name() {
		app.logger.Debugf(
			"transitioning from %v state to %v state",
			app.state.Name(),
			nextState.Name(),
		)
		nextState.Initialize()
		app.state = nextState
	}
	return true
}

func (app Application) CurrentState() string {
	return app.state.Name()
}
