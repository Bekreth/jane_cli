package app

import (
	"github.com/Bekreth/jane_cli/app/flag"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/states/auth"
	"github.com/Bekreth/jane_cli/app/states/booking"
	"github.com/Bekreth/jane_cli/app/states/charting"
	"github.com/Bekreth/jane_cli/app/states/initialize"
	"github.com/Bekreth/jane_cli/app/states/root"
	"github.com/Bekreth/jane_cli/app/states/schedule"
	Cache "github.com/Bekreth/jane_cli/cache"
	Client "github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/bekreth/screen_reader_terminal/terminal"
	"github.com/eiannone/keyboard"
)

type Application struct {
	logger    logger.Logger
	writer    terminal.Terminal
	state     states.State
	allStates []states.State
}

func NewApplication(
	logger logger.Logger,
	screenWriter terminal.Terminal,
	user *domain.User,
	client Client.Client,
	cache Cache.Cache,
) Application {
	rootState := root.NewState(
		logger.AddContext("state", "root"),
	)
	initState := initialize.NewState(
		logger.AddContext("state", "init"),
		user,
		rootState,
	)
	authState := auth.NewState(
		logger.AddContext("state", "auth"),
		client,
		rootState,
	)

	scheduleState := schedule.NewState(
		logger.AddContext("state", "schedule"),
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
		fetcher,
		rootState,
	)

	chartingState := charting.NewState(
		logger.AddContext("state", "charting"),
		fetcher,
		rootState,
	)

	rootState.RegisterStates(map[string]states.State{
		initState.Name():     initState,
		authState.Name():     authState,
		scheduleState.Name(): scheduleState,
		bookingState.Name():  bookingState,
		chartingState.Name(): chartingState,
	})
	screenWriter.AddBuffer(rootState.Initialize())
	screenWriter.Draw()

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
		app.writer.CurrentBuffer().AddString("Shutting down Jane CLI")
		app.writer.Draw()
		app.writer.NewLine()
		return false

	case keyboard.KeyCtrlU:
		app.writer.CurrentBuffer().SetString("")
		app.writer.Draw()
		return true

	case keyboard.KeyCtrlR:
		//TODO
		//app.state.RepeatLastOutput()
		return true

	case keyboard.KeyArrowRight:
		app.writer.CurrentBuffer().AdvanceCursor(1)
		app.writer.Draw()
		return true

	case keyboard.KeyArrowLeft:
		app.writer.CurrentBuffer().RetreatCursor(1)
		app.writer.Draw()
		return true

	case keyboard.KeyEnter:
		data, _ := app.writer.CurrentBuffer().Output()
		flags := flag.Parse(data)
		app.writer.NewLine()
		if _, exists := flags["help"]; exists {
			app.writer.CurrentBuffer().AddString(app.state.HelpString())
			app.writer.Draw()
		} else {
			if app.state.Submit(flags) {
				app.writer.Draw()
			}
		}
		app.writer.NewLine()
	}

	nextState, addNewLine := app.state.HandleKeyinput(character, key)
	if addNewLine {
		app.writer.Draw()
		app.writer.NewLine()
	}
	if nextState.Name() != app.state.Name() {
		app.logger.Debugf(
			"transitioning from %v state to %v state",
			app.state.Name(),
			nextState.Name(),
		)
		app.writer.AddBuffer(nextState.Initialize())
		app.state = nextState
	}
	app.writer.Draw()
	return true
}

func (app Application) CurrentState() string {
	return app.state.Name()
}
