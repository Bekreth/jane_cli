package charting

import "fmt"

func (state *chartingState) fetchCharts() {
	chartEntries, err := state.fetcher.FetchPatientCharts(state.builder.targetPatient.ID)
	if err != nil {
		//TODO
	}

	if len(chartEntries) == 0 {
		state.buffer.WriteStoreString(fmt.Sprintf(
			"no charts found for patient %v",
			state.builder.targetPatient.PrintName(),
		))
	} else if len(chartEntries) == 1 {
		state.builder.targetChart = chartEntries[0]
		state.builder.charts = chartEntries
	} else if len(chartEntries) <= 9 {
		state.builder.charts = chartEntries
	} else if len(chartEntries) > 9 {
		state.builder.charts = chartEntries[0:9]
	}
}
