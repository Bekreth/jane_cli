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
		state.buffer.AddString("aborting")
		state.builder = newChartingBuilder()
		state.nextState = state.rootState

	case "e":
		fallthrough
	case "E":
		state.builder.note = ""
		state.builder.substate = noteEditor

	default:
		state.buffer.AddString(fmt.Sprintf(
			"input of %v not support. Confirm, deny, or edit (Y/n/E)?",
			string(character),
		))
	}
}

func (state *chartingState) confirmSign() {
	state.buffer.AddString("submitting chart")
	chart, err := state.fetcher.CreatePatientCharts(
		state.builder.patientSelector.TargetSelection().GetID(),
		state.builder.appointmentSelector.TargetSelection().GetID(),
	)

	if err != nil {
		state.buffer.AddString(fmt.Sprintf("failed to create chart: %v", err))
		return
	}
	state.logger.Debugf("created chart", chart)

	err = state.fetcher.UpdatePatientChart(
		chart.ChartParts[0].ID,
		state.builder.note,
	)

	if err != nil {
		state.buffer.AddString(fmt.Sprintf("failed to update chart: %v", err))
		return
	}

	err = state.fetcher.SignChart(
		chart,
		state.builder.patientSelector.TargetSelection().GetID(),
	)
	if err != nil {
		state.buffer.AddString(fmt.Sprintf("failed to sign chart: %v", err))
		return
	}

	state.buffer.AddString("Successfully created chart!")
	state.buffer.Clear()
}
