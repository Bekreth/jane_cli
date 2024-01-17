package app

import (
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/states/auth"
	"github.com/Bekreth/jane_cli/app/states/booking"
	"github.com/Bekreth/jane_cli/app/states/charting"
	"github.com/Bekreth/jane_cli/app/states/initialize"
	"github.com/Bekreth/jane_cli/app/states/root"
	"github.com/Bekreth/jane_cli/app/states/schedule"
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
	state     states.State
	allStates []states.State
}

func NewApplication(
	logger logger.Logger,
	screenWriter terminal.ScreenWriter,
	user *domain.User,
	client Client.Client,
	cache Cache.Cache,
) Application {
	rootState := root.NewState(
		logger.AddContext("state", "root"),
		screenWriter,
	)
	initState := initialize.NewState(
		logger.AddContext("state", "init"),
		screenWriter,
		user,
		rootState,
	)
	authState := auth.NewState(
		logger.AddContext("state", "auth"),
		screenWriter,
		client,
		rootState,
	)

	scheduleState := schedule.NewState(
		logger.AddContext("state", "schedule"),
		screenWriter,
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
		screenWriter,
		fetcher,
		rootState,
	)

	chartingState := charting.NewState(
		logger.AddContext("state", "charting"),
		fetcher,
		screenWriter,
		rootState,
	)

	rootState.RegisterStates(map[string]states.State{
		initState.Name():     initState,
		authState.Name():     authState,
		scheduleState.Name(): scheduleState,
		bookingState.Name():  bookingState,
		chartingState.Name(): chartingState,
	})
	rootState.Initialize()

	return Application{
		logger: logger,
		writer: screenWriter,
		state:  rootState,
		allStates: []states.State{
			rootState,
			initState,
			authState,
			scheduleState,
			bookingState,
			chartingState,
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
	case keyboard.KeyCtrlR:
		app.state.RepeatLastOutput()
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
