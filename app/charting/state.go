package charting

import (
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/charts"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
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
	rootState terminal.State

	builder   chartingBuilder
	nextState terminal.State
	buffer    *terminal.Buffer
}

func NewState(
	logger logger.Logger,
	fetcher chartingDataFetcher,
	writer terminal.ScreenWriter,
	rootState terminal.State,
) terminal.State {
	buffer := terminal.NewBuffer(writer)
	return &chartingState{
		logger:    logger,
		fetcher:   fetcher,
		rootState: rootState,
		buffer:    &buffer,
	}
}

func (chartingState) Name() string {
	return "charting"
}

func (state *chartingState) Initialize() {
	state.logger.Debugf(
		"entering charting. available states %v",
		state.rootState.Name(),
	)
	state.builder = chartingBuilder{
		substate: argument,
		flow:     undefined,
	}
	state.nextState = state
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

var autocompletes = map[string]string{
	helpCommand:   "",
	readCommand:   "",
	createCommand: "",
}

func (state *chartingState) triggerAutocomplete() {
	words := strings.Split(state.buffer.Read(), " ")
	for key := range autocompletes {
		if strings.HasPrefix(key, words[len(words)-1]) {
			arguments := append(words[0:len(words)-1], key)
			state.buffer.WriteString(strings.Join(arguments, " ") + " ")
		}
	}
}

func (state *chartingState) ClearBuffer() {
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *chartingState) RepeatLastOutput() {
	state.buffer.WritePrevious()
}
