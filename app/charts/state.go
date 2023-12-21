package charts

import (
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type chartState struct {
	logger logger.Logger
	rootState terminal.State

	buffer    *terminal.Buffer
	nextState terminal.State
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	rootState terminal.State,
) terminal.State {
	buffer := terminal.NewBuffer(writer)
	return &chartState{
		logger:    logger,
		rootState: rootState,
		buffer:    &buffer,
	}
}

func (chartState) Name() string {
	return "charts"
}

func (state *chartState) Initialize() {
	state.logger.Debugf(
		"entering schedule. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (*chartState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	panic("unimplemented")
}

func (state *chartState) ClearBuffer() {
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *chartState) RepeatLastOutput() {
	state.buffer.WritePrevious()
}
