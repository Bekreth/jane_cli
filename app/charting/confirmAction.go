package charting

import "fmt"

func (state *chartingState) confirmAction(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		state.confirmSign()
		state.builder.substate = argument
		state.builder.flow = undefined
	case "n":
		fallthrough
	case "N":
		state.buffer.WriteStoreString("aborting")
		state.nextState = state.rootState
	case "e":
		fallthrough
	case "E":
		state.builder.note = ""
		state.builder.substate = noteEditor
	default:
		state.buffer.WriteStoreString(fmt.Sprintf(
			"input of %v not support. Confirm, deny, or edit (Y/n/E)?",
			string(character),
		))
	}
}

func (state *chartingState) confirmSign() {
	state.buffer.WriteStoreString("submitting chart")
	chart, err := state.fetcher.CreatePatientCharts(
		state.builder.targetPatient.ID,
		state.builder.targetAppointment.ID,
	)

	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to create chart: %v", err))
	}
	state.logger.Debugf("created chart: %v", chart)

	err = state.fetcher.UpdatePatientChart(
		chart.ChartParts[0].ID,
		state.builder.note,
	)

	//TODO: add siging
	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to update chart: %v", err))
	} else {
		state.buffer.WriteStoreString("Successfully created chart!")
		state.ClearBuffer()
	}
}
