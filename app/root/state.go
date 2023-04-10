package root

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type rootState struct {
	logger logger.Logger
	writer terminal.ScreenWriter

	states        map[string]terminal.State
	currentBuffer string
	nextState     terminal.State
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
) *rootState {
	return &rootState{
		logger: logger,
		writer: writer,
	}
}

func (rootState) Name() string {
	return "root"
}

func (root *rootState) Initialize() {
	stateNames := []string{}
	for key := range root.states {
		stateNames = append(stateNames, key)
	}
	root.logger.Debugf(
		"entering root. available states: %v",
		stateNames,
	)
	root.currentBuffer = ""
	root.nextState = root
	root.writer.NewLine()
	root.writer.WriteString("")
}

func (root *rootState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &root.currentBuffer, root.triggerAutocomplete, root.Submit)

	if character != 0 {
		root.currentBuffer += string(character)
	}
	root.writer.WriteString(root.currentBuffer)
	return root.nextState
}

func (root *rootState) triggerAutocomplete() {
	for _, stateName := range terminal.MapKeys(root.states) {
		if strings.HasPrefix(stateName, root.currentBuffer) {
			root.currentBuffer = stateName
		}
	}
}

func (root *rootState) Submit() {
	words := strings.Split(root.currentBuffer, " ")
	for _, stateName := range terminal.MapKeys(root.states) {
		if words[0] == stateName {
			// TODO: deal with passing arguments forward
			//root.arguments = words[1 : len(words)-1]
			root.nextState = root.states[stateName]
			return
		}
	}
	// TODO: Deal with failed command
	root.currentBuffer = ""
	root.writer.WriteString(fmt.Sprintf("'%v' is not a valid command", words[0]))
	root.writer.NewLine()
	root.nextState = root
}

func (root *rootState) RegisterStates(states map[string]terminal.State) {
	root.states = states
}
