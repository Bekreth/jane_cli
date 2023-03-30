package app

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type rootState struct {
	logger        logger.Logger
	scheduleState state
	currentBuffer string
}

func (rootState) name() string {
	return "root"
}

func (root rootState) initialize() {
	root.logger.Debugf("entering root. available states %v", root.scheduleState.name())
}

func (root *rootState) handleKeyinput(character rune, key keyboard.Key) state {
	switch key {
	case keyboard.KeyTab:
		root.triggerAutocomplete()
		// TODO: replace
		fmt.Println(root.currentBuffer)
		return root
	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		if len(root.currentBuffer) != 0 {
			root.currentBuffer = root.currentBuffer[0 : len(root.currentBuffer)-1]
		}
		// TODO: replace
		fmt.Println(root.currentBuffer)
	case keyboard.KeyEnter:
		return root.submit()
	}

	if character != 0 {
		root.currentBuffer += string(character)
		// TODO: replace
		fmt.Println(root.currentBuffer)
	}
	return root
}

func (root *rootState) shutdown() {
	root.currentBuffer = ""
}

func (root *rootState) triggerAutocomplete() {
	root.logger.Debugf("triggering autocomplete")

	if strings.HasPrefix(root.scheduleState.name(), root.currentBuffer) {
		root.currentBuffer = root.scheduleState.name()
	}
}

func (root *rootState) submit() state {
	if root.currentBuffer == root.scheduleState.name() {
		return root.scheduleState
	}
	return root
}
