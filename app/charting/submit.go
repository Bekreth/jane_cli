package charting

import (
	"github.com/Bekreth/jane_cli/app/terminal"
)

const dateFlag = "-d"
const patientFlag = "-p"
const noteFlag = "-n"

const readCommand = "read"
const createCommand = "create"
const helpCommand = "help"
const backCommand = ".."

func (state *chartingState) Submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.logger.Debugf("submitting query flags: %v", flags)
	state.buffer.Clear()
	var err error
	if _, exists := flags[backCommand]; exists {
		state.nextState = state.rootState
	} else if _, exists := flags[helpCommand]; exists {
		state.printHelp()
	} else if _, exists := flags[createCommand]; exists {
		state.builder, err = state.handleCreate(flags)
	} else if _, exists := flags[readCommand]; exists {
		state.builder, err = state.handleRead(flags)
	} else {
		state.buffer.WriteStoreString("No subcommand specified. Please specify 'read' or 'create'")
	}
	if err != nil {
		state.builder = newChartingBuilder()
		state.buffer.WriteStoreString(err.Error())
	}
}
