package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/domain/charts"
)

func (state *chartingState) fetchCharts() (
	interactive.Interactive[charts.ChartEntry],
	error,
) {
	chartEntries, err := state.fetcher.FetchPatientCharts(
		state.builder.patientSelector.TargetSelection().GetID(),
	)
	if err != nil {
		return nil, err
	}

	if len(chartEntries) == 0 {
		return nil, fmt.Errorf(
			"no charts found for patient %v",
			state.builder.patientSelector.TargetSelection().PrintSelector(),
		)
	} else if len(chartEntries) == 1 {
		return interactive.NewChartSelector(
			chartEntries[0],
			chartEntries,
		), nil
	} else {
		return interactive.NewChartSelector(
			charts.ChartEntry{},
			chartEntries,
		), nil
	}
}
