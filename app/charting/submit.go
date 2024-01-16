package charting

import (
	"github.com/Bekreth/jane_cli/app/interactive"
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
	//TODO: this is brittle
	flags := terminal.ParseFlags(state.buffer.Read())
	state.logger.Debugf("submitting query flags: %v", flags)
	state.buffer.Clear()
	var err error
	if _, exists := flags[backCommand]; exists {
		state.nextState = state.rootState
	} else if _, exists := flags[helpCommand]; exists {
		state.printHelp()
	} else if _, exists := flags[createCommand]; exists {
		state.builder.flow = create
		state.builder.substate = unknown
		state.handleCreateNote(flags)
	} else if _, exists := flags[readCommand]; exists {
		state.builder.patientSelector, err = state.fetchPatients(flags)
		if err != nil {
			state.buffer.WriteStoreString(err.Error())
		} else {
			state.builder.flow = read
			state.builder.substate = unknown
		}
	} else {
		state.buffer.WriteStoreString("No subcommand specified. Please specify 'read' or 'create'")
	}
}

func (state *chartingState) handleCreateNote(flags map[string]string) {
	// Setup Patient
	patientName, exists := flags[patientFlag]
	if !exists {
		state.buffer.WriteStoreString("missing argument -p")
		return
	}

	var err error
	targetPatient, patients, err := util.ParsePatientValue(
		state.fetcher,
		patientName,
	)
	state.builder.patientSelector = interactive.NewPatientSelector(targetPatient, patients)
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
		state.builder.date = parsedDate
	}

	// Setup Note
	state.builder.note, exists = flags[noteFlag]
}
