package charting

import (
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
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
	if _, exists := flags[backCommand]; exists {
		state.nextState = state.rootState
		return
	} else if _, exists := flags[helpCommand]; exists {
		state.printHelp()
		return
	} else if _, exists := flags[readCommand]; exists {
		state.fetchPatients(flags)
	} else if _, exists := flags[createCommand]; exists {
		state.handleCreateNote(flags)
	} else {
		state.buffer.WriteStoreString("No subcommand specified. Please specify 'read' or 'create'")
	}
}

func (state *chartingState) handleCreateNote(flags map[string]string) {
	builder := chartingBuilder{
		substate: unknown,
		flow:     create,
	}

	// Setup Patient
	patientName, exists := flags[patientFlag]
	if !exists {
		state.buffer.WriteStoreString("missing argument -p")
		return
	}

	var err error
	builder.targetPatient, builder.patients, err = util.ParsePatientValue(
		state.fetcher,
		patientName,
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	// Setup Date of Appointment
	date, exists := flags[dateFlag]
	if exists {
		parsedDate, err := util.ParseDate(
			util.DateFormat,
			util.YearDateFormat,
			date,
		)
		if err != nil {
			state.buffer.WriteStoreString(err.Error())
		}
		builder.date = parsedDate
	}

	// Setup Note
	builder.note, exists = flags[noteFlag]

	state.builder = builder
}
