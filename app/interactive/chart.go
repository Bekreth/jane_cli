package interactive

import (
	"github.com/Bekreth/jane_cli/domain/charts"
)

type SelectedChart struct {
	charts.ChartEntry
}

func (selection SelectedChart) GetID() int {
	return selection.ID
}

func (SelectedChart) PrintHeader() string {
	return "Selected desired chart"
}

func (selection SelectedChart) PrintSelector() string {
	return selection.ChartEntry.PrintSelector()
}

func (selection SelectedChart) Deref() charts.ChartEntry {
	return selection.ChartEntry
}

func (selection SelectedChart) hasSelection() bool {
	return selection.ID != 0
}

func EmptyChartSelector() Interactive[charts.ChartEntry] {
	var output *selector[charts.ChartEntry]
	return output
}

func NewChartSelector(
	selected charts.ChartEntry,
	possible []charts.ChartEntry,
) Interactive[charts.ChartEntry] {
	var selectedChart SelectedChart
	if selected.ID != 0 {
		selectedChart = SelectedChart{selected}
	}
	possibleCharts := make([]Selection[charts.ChartEntry], len(possible))
	for i, selection := range possible {
		possibleCharts[i] = SelectedChart{selection}
	}
	return &selector[charts.ChartEntry]{
		page:              0,
		possibleSelection: possibleCharts,
		selected:          selectedChart,
	}
}
