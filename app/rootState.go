package app

import (
	"strings"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type rootState struct {
	logger        logger.Logger
	writer        screenWriter
	states        map[string]state
	currentBuffer string
}

func (rootState) name() string {
	return "root"
}

func (root rootState) initialize() {
	stateNames := []string{}
	for key := range root.states {
		stateNames = append(stateNames, key)
	}
	root.logger.Debugf(
		"entering root. available states: %v",
		stateNames,
	)
	root.writer.writeString("")
}

func (root *rootState) handleKeyinput(character rune, key keyboard.Key) state {
	var output state
	output = root
	switch key {
	case keyboard.KeySpace:
		root.currentBuffer += string(" ")

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
	for _, stateName := range mapKeys(root.states) {
		if strings.HasPrefix(stateName, root.currentBuffer) {
			root.currentBuffer = stateName
		}
	}
}

func (root *rootState) submit() state {
	for _, stateName := range mapKeys(root.states) {
		if root.currentBuffer == stateName {
			return root.states[stateName]
		}
	}
	root.writer.newLine()
	return root
}
