package charting

import (
	"strings"

	"github.com/Bekreth/jane_cli/app/flag"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/charts"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
	terminal "github.com/bekreth/screen_reader_terminal"
)

type chartingDataFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
	FindAppointments(
		startDate schedule.JaneTime,
		endDate schedule.JaneTime,
		patientName string,
	) ([]schedule.Appointment, error)

	FetchPatientCharts(patientID int) ([]charts.ChartEntry, error)
	CreatePatientCharts(patientID int, appointmentID int) (charts.ChartEntry, error)
	UpdatePatientChart(chartPartID int, chartText string) error
	SetChartingAppointment(chartID int, appointmentID int) error
	SignChart(chart charts.ChartEntry, patientID int) error
}

type chartingState struct {
	logger    logger.Logger
	fetcher   chartingDataFetcher
	rootState states.State

	builder   chartingBuilder
	nextState states.State
	buffer    *terminal.Buffer
}

func NewState(
	logger logger.Logger,
	fetcher chartingDataFetcher,
	rootState states.State,
) states.State {
	buffer := terminal.NewBuffer()
	return &chartingState{
		logger:    logger,
		fetcher:   fetcher,
		rootState: rootState,
		buffer:    buffer.SetPrefix("charting: "),
	}
}

func (chartingState) Name() string {
	return "charting"
}

func (state *chartingState) Initialize() *terminal.Buffer {
	state.logger.Debugf(
		"entering charting. available states %v",
		state.rootState.Name(),
	)
	state.builder = newChartingBuilder()
	state.nextState = state
	state.buffer.Clear()
	state.buffer.SetPrefix("charting: ")
	return state.buffer
}

var autocompletes = map[string]string{
	readCommand:   "",
	createCommand: "",
}

func (state *chartingState) triggerAutocomplete() {
	data, _ := state.buffer.Output()
	flags := flag.Parse(data)

	for key := range autocompletes {
		for flagKey := range flags {
			if strings.HasPrefix(key, flagKey) {
				state.buffer.AddString(strings.Replace(key, flagKey, "", 1))
			}
		}
	}
}
