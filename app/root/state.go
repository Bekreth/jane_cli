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

func (state *rootState) Initialize() {
	stateNames := []string{}
	for key := range state.states {
		stateNames = append(stateNames, key)
	}
	state.logger.Debugf(
		"entering root. available states: %v",
		stateNames,
	)
	state.currentBuffer = ""
	state.nextState = state
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *rootState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &state.currentBuffer, state.triggerAutocomplete, state.Submit)

	if character != 0 {
		state.currentBuffer += string(character)
	}
	state.writer.WriteString(state.currentBuffer)
	return state.nextState
}

func (state *rootState) triggerAutocomplete() {
	for _, stateName := range terminal.MapKeys(state.states) {
		if strings.HasPrefix(stateName, state.currentBuffer) {
			state.currentBuffer = stateName
		}
	}
}

func (state *rootState) Submit() {
	words := strings.Split(state.currentBuffer, " ")
	if words[0] == "help" {
		state.printHelp()
		state.currentBuffer = ""
	} else {
		for _, stateName := range terminal.MapKeys(state.states) {
			if words[0] == stateName {
				// TODO: deal with passing arguments forward
				//root.arguments = words[1 : len(words)-1]
				state.nextState = state.states[stateName]
				return
			}
		}
		// TODO: Deal with failed command
		state.currentBuffer = ""
		state.writer.WriteString(fmt.Sprintf("'%v' is not a valid command", words[0]))
		state.writer.NewLine()
		state.nextState = state
	}
}

func (state *rootState) RegisterStates(states map[string]terminal.State) {
	state.states = states
}

func (state *rootState) printHelp() {
	// TODO: automate this list of elements
	state.writer.WriteStringf("available commands: auth, init, schedule, booking")
	state.writer.NewLine()
}
