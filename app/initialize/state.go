package initialize

import (
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

const username = "-u"
const clinicDomain = "-c"

type initState struct {
	logger    logger.Logger
	writer    terminal.ScreenWriter
	user      *domain.User
	rootState terminal.State

	currentBuffer string
	nextState     terminal.State
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	user *domain.User,
	rootState terminal.State,
) terminal.State {
	return &initState{
		logger:    logger,
		writer:    writer,
		user:      user,
		rootState: rootState,
	}
}

func (initState) Name() string {
	return "init"
}

func (state *initState) Initialize() {
	state.logger.Debugf(
		"entering init. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *initState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &state.currentBuffer, state.triggerAutocomplete, state.submit)

	if character != 0 {
		state.currentBuffer += string(character)
	}

	state.writer.WriteString(state.currentBuffer)
	return state.nextState
}

func (state *initState) triggerAutocomplete() {
}

func (state *initState) ClearBuffer() {
	state.currentBuffer = ""
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *initState) printHelp() {
	// TODO: automate this list of elements
	state.writer.WriteStringf(
		"init should only need to be run the first time your setup your client:\n%v\n%v\n",
		"\t-u\tusername used to log in to Jane",
		"\t-c\tthe name of the clinic",
	)
	state.writer.NewLine()
}
