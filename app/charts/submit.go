package charts

import (
	"github.com/Bekreth/jane_cli/app/terminal"
)

const dateFlag = "-d"
const patientFlag = "-p"

const readCommand = "read"
const createCommand = "create"
const helpCommand = "help"
const backCommand = ".."

func (state *chartingState) Submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.logger.Debugf("submitting query flags: %v", flags)
	state.buffer.Clear()
	if _, exists := flags[backCommand]; exists {
		state.nextState = state.rootState
		return
	} else if _, exists := flags[helpCommand]; exists {
		state.printHelp()
		return
	} else if _, exists := flags[readCommand]; exists {
		state.fetchPatients(flags)
	} else if _, exists := flags[createCommand]; exists {
		panic("unimplemented")
	} else {
		state.buffer.WriteStoreString("No subcommand specified. Please specify 'read' or 'create'")
	}
}
