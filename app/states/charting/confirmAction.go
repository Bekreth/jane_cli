package charting

import "fmt"

func (state *chartingState) confirmAction(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		state.confirmSign()
		state.builder = newChartingBuilder()
		state.nextState = state.rootState

	case "n":
		fallthrough
	case "N":
		state.buffer.WriteStoreString("aborting")
		state.builder = newChartingBuilder()
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
		state.builder.patientSelector.TargetSelection().GetID(),
		state.builder.appointmentSelector.TargetSelection().GetID(),
	)

	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to create chart: %v", err))
		return
	}
	state.logger.Debugf("created chart: %v", chart)

	err = state.fetcher.UpdatePatientChart(
		chart.ChartParts[0].ID,
		state.builder.note,
	)

	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to update chart: %v", err))
		return
	}

	err = state.fetcher.SignChart(
		chart,
		state.builder.patientSelector.TargetSelection().GetID(),
	)
	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to sign chart: %v", err))
		return
	}

	state.buffer.WriteStoreString("Successfully created chart!")
	state.ClearBuffer()
}
