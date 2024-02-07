package initialize

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	terminal "github.com/bekreth/screen_reader_terminal"
	"github.com/eiannone/keyboard"
)

const username = "-u"
const clinicDomain = "-c"

type initState struct {
	logger    logger.Logger
	user      *domain.User
	rootState states.State

	nextState states.State
	buffer    *terminal.Buffer
}

func NewState(
	logger logger.Logger,
	user *domain.User,
	rootState states.State,
) states.State {
	buffer := terminal.NewBuffer()
	return &initState{
		logger:    logger,
		user:      user,
		rootState: rootState,
		buffer:    buffer.SetPrefix("init: "),
	}
}

func (initState) Name() string {
	return "init"
}

func (state *initState) Initialize() *terminal.Buffer {
	state.logger.Debugf(
		"entering init. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.buffer.Clear()
	return state.buffer
}

func (state *initState) HandleKeyinput(character rune, key keyboard.Key) states.State {
	util.KeyHandler(key, state.buffer, state.triggerAutocomplete)
	if character != 0 {
		state.buffer.AddCharacter(character)
	}
	return state.nextState
}

func (state *initState) triggerAutocomplete() {
}

func (state *initState) HelpString() string {
	return fmt.Sprintf(
		"init should only need to be run the first time your setup your client:\n%v\n%v",
		"\t-c\tthe name of the clinic",
		"\t-u\tusername used to log in to Jane",
	)
}
