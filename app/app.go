package app

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type Application struct {
	logger    logger.Logger
	state     state
	allStates []state
}

func NewApplication(
	logger logger.Logger,
	user domain.User,
) Application {

	root := rootState{
		logger: logger.AddContext("state", "root"),
		writer: screenWriter{"jane>"},
	}
	init := initState{
		logger: logger.AddContext("state", "init"),
		writer: screenWriter{"init>"},
		user:   user,
	}
	/*
		auth := authState{
			logger:    logger.AddContext("state", "auth"),
			writer:    screenWriter{"auth>"},
			rootState: tempState,
		}
	*/
	schedule := scheduleState{
		logger:   logger.AddContext("state", "schedule"),
		writer:   screenWriter{"schedule>"},
		subState: none,
	}

	root.states = map[string]state{
		schedule.name(): &schedule,
		init.name():     &init,
	}

	//auth.rootState = &root
	schedule.rootState = &root
	init.rootState = &root

	root.initialize()

	return Application{
		logger:    logger,
		state:     &root,
		allStates: []state{&root, &schedule, &init},
	}
}

const ctrlC = rune(0x03)

func (app *Application) HandleKeyinput(character rune, key keyboard.Key) bool {

	switch key {
	case keyboard.KeyCtrlC:
		app.logger.Infoln("exiting the application")
		for _, state := range app.allStates {
			state.shutdown()
		}
		return false
	}

	nextState := app.state.handleKeyinput(character, key)
	if nextState.name() != app.state.name() {
		app.logger.Debugf(
			"transitioning from %v state to %v state",
			app.state.name(),
			nextState.name(),
		)
		nextState.initialize()
		app.state.shutdown()
		app.state = nextState
	}
	return true
}

func (app Application) CurrentState() string {
	return app.state.name()
}
