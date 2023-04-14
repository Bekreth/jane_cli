package auth

import (
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
	writer        terminal.ScreenWriter
	authenticator authenticator
	rootState     terminal.State

	nextState     terminal.State
	currentBuffer string
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	authenticator authenticator,
	rootState terminal.State,
) terminal.State {
	return &authState{
		logger:        logger,
		writer:        writer,
		authenticator: authenticator,
		rootState:     rootState,
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
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *authState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &state.currentBuffer, state.triggerAutocomplete, state.submit)

	if character != 0 {
		state.currentBuffer += string(character)
	}

	state.writer.WriteString(state.currentBuffer)
	return state.nextState
}

func (state *authState) triggerAutocomplete() {
}

func (state *authState) submit() {
	flags := terminal.ParseFlags(state.currentBuffer)
	if _, exists := flags["help"]; exists {
		state.printHelp()
		state.currentBuffer = ""
		return
	}
	var err error
	if password, ok := flags[passwordFlag]; ok {
		err = state.authenticator.Login(password)
	} else {
		state.writer.WriteString("password not provided")
	}

	state.currentBuffer = ""

	if err != nil {
		state.writer.WriteStringf("failed to login: %v", err)
		state.writer.NewLine()
	} else {
		state.writer.WriteString("login successful")
		state.writer.NewLine()
	}

	state.nextState = state.rootState
}

func (state *authState) ClearBuffer() {
	state.currentBuffer = ""
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *authState) printHelp() {
	// TODO: automate this list of elements
	state.writer.WriteStringf(
		"auth is used to login to your account:\n%v",
		"\t-p\tYour password",
	)
	state.writer.NewLine()
}
