package app

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type rootState struct {
	logger        logger.Logger
	writer        screenWriter
	states        map[string]state
	currentBuffer string
	nextState     state
}

func (rootState) name() string {
	return "root"
}

func (root *rootState) initialize() {
	stateNames := []string{}
	for key := range root.states {
		stateNames = append(stateNames, key)
	}
	root.logger.Debugf(
		"entering root. available states: %v",
		stateNames,
	)
	root.nextState = root
	root.writer.newLine()
	root.writer.writeString("")
}

func (root *rootState) handleKeyinput(character rune, key keyboard.Key) state {
	keyHandler(key, &root.currentBuffer, root.triggerAutocomplete, root.submit)

	if character != 0 {
		root.currentBuffer += string(character)
	}
	root.writer.writeString(root.currentBuffer)
	return root.nextState
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

func (root *rootState) submit() {
	words := strings.Split(root.currentBuffer, " ")
	for _, stateName := range mapKeys(root.states) {
		if words[0] == stateName {
			// TODO: deal with passing arguments forward
			//root.arguments = words[1 : len(words)-1]
			root.nextState = root.states[stateName]
			return
		}
	}
	// TODO: Deal with failed command
	root.currentBuffer = ""
	root.writer.writeString(fmt.Sprintf("'%v' is not a valid command", words[0]))
	root.writer.newLine()
	root.nextState = root
}
