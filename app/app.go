package app

import (
	"github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type Application struct {
	logger    logger.Logger
	writer    screenWriter
	state     state
	allStates []state
}

type tempPatientFetcher struct{}

func (tempPatientFetcher) FindPatient(patientName string) ([]domain.Patient, error) {
	return []domain.Patient{
		{
			ID:                 1,
			FirstName:          "Billy",
			LastName:           "Bob",
			PreferredFirstName: "Will",
		},
		{
			ID:                 2,
			FirstName:          "Mark",
			LastName:           "Walberg",
			PreferredFirstName: "Will",
		},
		{
			ID:                 3,
			FirstName:          "Jimmy",
			LastName:           "Neutron",
			PreferredFirstName: "Will",
		},
	}, nil
}

func NewApplication(
	logger logger.Logger,
	user *domain.User,
	client client.Client,
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
	auth := authState{
		logger:        logger.AddContext("state", "auth"),
		writer:        screenWriter{"auth>"},
		authenticator: client,
	}
	schedule := scheduleState{
		logger:  logger.AddContext("state", "schedule"),
		writer:  screenWriter{"schedule>"},
		fetcher: client,
	}
	booking := bookingState{
		logger:         logger.AddContext("state", "booking"),
		writer:         screenWriter{"booking>"},
		patientFetcher: tempPatientFetcher{},
		// TODO: Add interface for clients
	}

	root.states = map[string]state{
		auth.name():     &auth,
		schedule.name(): &schedule,
		init.name():     &init,
		booking.name():  &booking,
	}

	auth.rootState = &root
	schedule.rootState = &root
	init.rootState = &root
	booking.rootState = &root

	root.initialize()

	return Application{
		logger:    logger,
		writer:    screenWriter{},
		state:     &root,
		allStates: []state{&root, &schedule, &init, &auth},
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
		app.writer.newLine()
		app.writer.writeString("Shutting down Jane CLI")
		app.writer.newLine()
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
