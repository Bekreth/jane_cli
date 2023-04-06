package app

import (
	"fmt"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type authenticator interface {
	Login(password string) error
}

const passwordFlag = "-p"

type authState struct {
	logger        logger.Logger
	writer        screenWriter
	rootState     state
	nextState     state
	authenticator authenticator
	currentBuffer string
}

func (authState) name() string {
	return "auth"
}

func (auth *authState) initialize() {
	auth.logger.Debugf(
		"entering authenticator. available states %v",
		auth.rootState.name(),
	)
	auth.nextState = auth
	auth.writer.newLine()
	auth.writer.writeString("")
}

func (auth *authState) handleKeyinput(character rune, key keyboard.Key) state {
	keyHandler(key, &auth.currentBuffer, auth.triggerAutocomplete, auth.submit)

	if character != 0 {
		auth.currentBuffer += string(character)
	}

	auth.writer.writeString(auth.currentBuffer)
	return auth.nextState
}

func (auth *authState) triggerAutocomplete() {
}

func (auth *authState) submit() {
	flags := parseFlags(auth.currentBuffer)
	var err error
	if password, ok := flags[passwordFlag]; ok {
		err = auth.authenticator.Login(password)
	} else {
		auth.writer.writeString("password not provided")
	}

	if err != nil {
		auth.writer.writeString(fmt.Sprintf("failed to login: %v", err))
		auth.writer.newLine()
	} else {
		auth.writer.writeString("login successful")
		auth.writer.newLine()
	}

	auth.nextState = auth.rootState
}

func (auth *authState) shutdown() {
	auth.currentBuffer = ""
}
