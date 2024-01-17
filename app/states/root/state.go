package root

import (
	"strings"

	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type rootState struct {
	logger logger.Logger

	states    map[string]states.State
	buffer    *terminal.Buffer
	nextState states.State
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
) *rootState {
	buffer := terminal.NewBuffer(writer, "jane")
	return &rootState{
		logger: logger,
		buffer: &buffer,
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
	state.nextState = state
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *rootState) HandleKeyinput(character rune, key keyboard.Key) states.State {
	terminal.KeyHandler(key, state.buffer, state.triggerAutocomplete, state.submit)
	state.buffer.AddCharacter(character)
	state.buffer.Write()
	return state.nextState
}

func (state *rootState) triggerAutocomplete() {
	for _, stateName := range terminal.MapKeys(state.states) {
		if strings.HasPrefix(stateName, state.buffer.Read()) {
			state.buffer.WriteString(stateName)
		}
	}
}

func (state *rootState) submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.buffer.Clear()

	if _, exists := flags["help"]; exists {
		state.printHelp()
		return
	}

	for _, stateName := range terminal.MapKeys(state.states) {
		if _, exists := flags[stateName]; exists {
			state.nextState = state.states[stateName]
			return
		}
	}
	state.buffer.WriteStoreString("invalid command")
}

func (state *rootState) RegisterStates(states map[string]states.State) {
	state.states = states
}

func (state *rootState) ClearBuffer() {
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *rootState) RepeatLastOutput() {
	state.buffer.WritePrevious()
}

func (state *rootState) printHelp() {
	// TODO: automate this list of elements
	output := strings.Join([]string{
		"available commands: auth, init, schedule, booking, charting",
		"available autocommands:",
		"\tctrl+c\tclose Jane App",
		"\tctrl+u\tclear the current line",
		"\tctrl+r\trepeat last message",
	}, "\n")
	state.buffer.WriteStoreString(output)
}
