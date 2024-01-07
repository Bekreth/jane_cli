package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/domain/charts"
)

func (state *chartingState) fetchCharts() {
	chartEntries, err := state.fetcher.FetchPatientCharts(
		state.builder.patientSelector.TargetSelection().GetID(),
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	if len(chartEntries) == 0 {
		state.buffer.WriteStoreString(fmt.Sprintf(
			"no charts found for patient %v",
			state.builder.patientSelector.TargetSelection().PrintSelector(),
		))
	} else if len(chartEntries) == 1 {
		state.builder.chartSelector = interactive.NewChartSelector(
			chartEntries[0],
			chartEntries,
		)
	} else {
		state.builder.chartSelector = interactive.NewChartSelector(
			charts.ChartEntry{},
			chartEntries,
		)
	}
}
