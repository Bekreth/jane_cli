package initialize

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

const username = "-u"
const clinicDomain = "-c"

type initState struct {
	logger    logger.Logger
	user      *domain.User
	rootState terminal.State

	nextState terminal.State
	buffer    *terminal.Buffer
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	user *domain.User,
	rootState terminal.State,
) terminal.State {
	buffer := terminal.NewBuffer(writer)
	return &initState{
		logger:    logger,
		user:      user,
		rootState: rootState,
		buffer:    &buffer,
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
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *initState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, state.buffer, state.triggerAutocomplete, state.submit)
	state.buffer.AddCharacter(character)
	state.buffer.Write()
	return state.nextState
}

func (state *initState) triggerAutocomplete() {
}

func (state *initState) ClearBuffer() {
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *initState) RepeatLastOutput() {
	state.buffer.WritePrevious()
}

func (state *initState) printHelp() {
	// TODO: automate this list of elements
	state.buffer.WriteStoreString(fmt.Sprintf(
		"init should only need to be run the first time your setup your client:\n%v\n%v",
		"\t-u\tusername used to log in to Jane",
		"\t-c\tthe name of the clinic",
	))
}
