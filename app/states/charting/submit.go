package charting

const dateFlag = "-d"
const patientFlag = "-p"
const noteFlag = "-n"

const readCommand = "read"
const createCommand = "create"

func (state *chartingState) Submit(flags map[string]string) bool {
	if state.builder.substate == noteEditor {
		return false
	}
	state.logger.Debugf("submitting query flags: %v", flags)
	state.buffer.Clear()
	var err error
	if _, exists := flags[createCommand]; exists {
		state.builder, err = state.handleCreate(flags)
	} else if _, exists := flags[readCommand]; exists {
		state.builder, err = state.handleRead(flags)
	} else {
		state.buffer.AddString("No subcommand specified. Please specify 'read' or 'create'")
	}
	if err != nil {
		state.builder = newChartingBuilder()
		state.buffer.AddString(err.Error())
	}
	return true
}
