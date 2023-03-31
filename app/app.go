package app

import (
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type Application struct {
	logger    logger.Logger
	state     state
	allStates []state
}

func NewApplication(logger logger.Logger) Application {
	tempState := noState{}

	root := rootState{
		logger:        logger.AddContext("state", "root"),
		writer:        screenWriter{"jane>"},
		scheduleState: tempState,
	}
	schedule := scheduleState{
		logger:    logger.AddContext("state", "schedule"),
		writer:    screenWriter{"schedule>"},
		rootState: tempState,
		subState:  none,
	}

	root.scheduleState = &schedule
	schedule.rootState = &root

	root.initialize()

	return Application{
		logger:    logger,
		state:     &root,
		allStates: []state{&root, &schedule},
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
