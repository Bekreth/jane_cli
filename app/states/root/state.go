package root

import (
	"strings"

	"github.com/Bekreth/jane_cli/app/flag"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/logger"
	terminal "github.com/bekreth/screen_reader_terminal"
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
) *rootState {
	buffer := terminal.NewBuffer()
	return &rootState{
		logger: logger,
		buffer: buffer.SetPrefix("root: "),
	}
}

func (rootState) Name() string {
	return "root"
}

func (state *rootState) Initialize() *terminal.Buffer {
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
	return state.buffer
}

func (state *rootState) HandleKeyinput(character rune, key keyboard.Key) states.State {
	util.KeyHandler(key, state.buffer, state.triggerAutocomplete)
	if character != 0 {
		state.buffer.AddCharacter(character)
	}
	return state.nextState
}

func (state *rootState) triggerAutocomplete() {
	data, _ := state.buffer.Output()
	flags := flag.Parse(data)

	for stateName := range state.states {
		for flagKey := range flags {
			if strings.HasPrefix(stateName, flagKey) {
				state.buffer.AddString(strings.Replace(stateName, flagKey, "", 1))
			}
		}
	}
}

func (state *rootState) Submit(flags map[string]string) bool {
	for _, stateName := range util.MapKeys(state.states) {
		if _, exists := flags[stateName]; exists {
			state.nextState = state.states[stateName]
			return true
		}
	}
	state.buffer.AddString("invalid command")
	return true
}

func (state *rootState) RegisterStates(states map[string]states.State) {
	state.states = states
}

func (state *rootState) HelpString() string {
	return strings.Join([]string{
		"available commands: auth, init, schedule, booking, charting",
		"available autocommands:",
		"\tctrl+c\tclose Jane App",
		"\tctrl+u\tclear the current line",
		"\tctrl+r\trepeat last message",
	}, "\n")
}
