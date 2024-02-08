package auth

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/logger"
	terminal "github.com/bekreth/screen_reader_terminal"
	"github.com/eiannone/keyboard"
)

type authenticator interface {
	Login(password string) error
}

const passwordFlag = "-p"

type authState struct {
	logger        logger.Logger
	authenticator authenticator
	rootState     states.State

	nextState states.State
	buffer    *terminal.Buffer
}

func NewState(
	logger logger.Logger,
	authenticator authenticator,
	rootState states.State,
) states.State {
	buffer := terminal.NewBuffer()
	return &authState{
		logger:        logger,
		authenticator: authenticator,
		rootState:     rootState,
		buffer:        buffer.SetPrefix("auth: "),
	}
}

func (authState) Name() string {
	return "auth"
}

func (state *authState) Initialize() *terminal.Buffer {
	state.logger.Debugf(
		"entering authenticator. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.buffer.Clear()
	return state.buffer
}

func (state *authState) HandleKeyinput(character rune, key keyboard.Key) (states.State, bool) {
	util.KeyHandler(key, state.buffer, state.triggerAutocomplete)
	if character != 0 {
		state.buffer.AddCharacter(character)
	}
	return state.nextState, false
}

func (state *authState) triggerAutocomplete() {
}

func (state *authState) Submit(flags map[string]string) bool {
	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
		return true
	}

	var output string

	if password, ok := flags[passwordFlag]; ok {
		err := state.authenticator.Login(password)
		if err != nil {
			output = fmt.Sprintf("failed to login: %v", err)
		} else {
			output = "login successful"
			state.nextState = state.rootState
		}
	} else {
		output = "password not provided"
	}
	state.buffer.AddString(output)
	return true
}

func (state *authState) HelpString() string {
	return fmt.Sprintf(
		"auth is used to login to your account:\n%v",
		"\t-p\tYour password",
	)
}
