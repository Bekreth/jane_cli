package app

import (
	"fmt"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type authenticator interface {
	LoginRequired() bool
	Login(password string) error
}

const passwordFlag = "-p"

type authState struct {
	logger        logger.Logger
	writer        screenWriter
	rootState     state
	authenticator authenticator
	currentBuffer string
}

func (authState) name() string {
	return "auth"
}

func (auth authState) initialize() {
	auth.logger.Debugf(
		"entering authenticator. available states %v",
		auth.rootState.name(),
	)
	auth.writer.writeString("")
}

func (auth *authState) handleKeyinput(character rune, key keyboard.Key) state {
	var output state
	output = auth

	switch key {
	case keyboard.KeySpace:
		auth.currentBuffer += string(" ")

	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		if len(auth.currentBuffer) != 0 {
			auth.currentBuffer = auth.currentBuffer[0 : len(auth.currentBuffer)-1]
		}

	case keyboard.KeyEnter:
		output = auth.submit()
	}

	if character != 0 {
		auth.currentBuffer += string(character)
	}

	auth.writer.writeString(auth.currentBuffer)
	return output
}

func (auth *authState) submit() state {
	flags := parseFlags(auth.currentBuffer)
	var err error
	if password, ok := flags[passwordFlag]; ok {
		err = auth.authenticator.Login(password)
	} else {
		auth.writer.writeString("password not provided")
	}
	if err != nil {
		auth.writer.writeString(fmt.Sprintf("failed to login: %v", err))
	} else {
		auth.writer.writeString("login successful")
	}
	return auth.rootState
}

func (auth *authState) shutdown() {}
