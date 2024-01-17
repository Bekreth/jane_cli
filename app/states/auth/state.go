package auth

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/logger"
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
	writer terminal.ScreenWriter,
	authenticator authenticator,
	rootState states.State,
) states.State {
	buffer := terminal.NewBuffer(writer, "auth")
	return &authState{
		logger:        logger,
		authenticator: authenticator,
		rootState:     rootState,
		buffer:        &buffer,
	}
}

func (authState) Name() string {
	return "auth"
}

func (state *authState) Initialize() {
	state.logger.Debugf(
		"entering authenticator. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.buffer.Clear()
	state.buffer.WriteNewLine()
}

func (state *authState) HandleKeyinput(character rune, key keyboard.Key) states.State {
	terminal.KeyHandler(key, state.buffer, state.triggerAutocomplete, state.submit)
	state.buffer.AddCharacter(character)
	state.buffer.Write()
	return state.nextState
}

func (state *authState) triggerAutocomplete() {
}

func (state *authState) submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.buffer.Clear()
	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
		return
	} else if _, exists := flags["help"]; exists {
		state.printHelp()
		return
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
	state.buffer.WriteStoreString(output)
}

func (state *authState) ClearBuffer() {
	state.buffer.Clear()
	state.buffer.WriteNewLine()
}

func (state *authState) RepeatLastOutput() {
	state.buffer.WritePrevious()
}

func (state *authState) printHelp() {
	// TODO: automate this list of elements
	state.buffer.WriteStoreString(fmt.Sprintf(
		"auth is used to login to your account:\n%v",
		"\t-p\tYour password",
	))
}
