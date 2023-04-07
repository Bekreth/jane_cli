package auth

import (
	"fmt"

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

func (auth *authState) Initialize() {
	auth.logger.Debugf(
		"entering authenticator. available states %v",
		auth.rootState.Name(),
	)
	auth.nextState = auth
	auth.writer.NewLine()
	auth.writer.WriteString("")
}

func (auth *authState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &auth.currentBuffer, auth.triggerAutocomplete, auth.submit)

	if character != 0 {
		auth.currentBuffer += string(character)
	}

	auth.writer.WriteString(auth.currentBuffer)
	return auth.nextState
}

func (auth *authState) triggerAutocomplete() {
}

func (auth *authState) submit() {
	flags := terminal.ParseFlags(auth.currentBuffer)
	var err error
	if password, ok := flags[passwordFlag]; ok {
		err = auth.authenticator.Login(password)
	} else {
		auth.writer.WriteString("password not provided")
	}

	auth.currentBuffer = ""

	if err != nil {
		auth.writer.WriteString(fmt.Sprintf("failed to login: %v", err))
		auth.writer.NewLine()
	} else {
		auth.writer.WriteString("login successful")
		auth.writer.NewLine()
	}

	auth.nextState = auth.rootState
}
