package app

import (
	"strings"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type rootState struct {
	logger        logger.Logger
	writer        screenWriter
	scheduleState state
	currentBuffer string
}

func (rootState) name() string {
	return "root"
}

func (root rootState) initialize() {
	root.logger.Debugf("entering root. available states %v", root.scheduleState.name())
	root.writer.writeString("")
}

func (root *rootState) handleKeyinput(character rune, key keyboard.Key) state {
	var output state
	output = root
	switch key {
	case keyboard.KeyTab:
		root.triggerAutocomplete()

	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		if len(root.currentBuffer) != 0 {
			root.currentBuffer = root.currentBuffer[0 : len(root.currentBuffer)-1]
		}

	case keyboard.KeyEnter:
		output = root.submit()
	}

	if character != 0 {
		root.currentBuffer += string(character)
	}
	root.writer.writeString(root.currentBuffer)
	return output
}

func (root *rootState) shutdown() {
	root.currentBuffer = ""
}

func (root *rootState) triggerAutocomplete() {
	if strings.HasPrefix(root.scheduleState.name(), root.currentBuffer) {
		root.currentBuffer = root.scheduleState.name()
	}
}

func (root *rootState) submit() state {
	if root.currentBuffer == root.scheduleState.name() {
		return root.scheduleState
	}
	root.writer.newLine()
	return root
}
