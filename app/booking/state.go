package booking

import (
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
)

type bookingDataFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
	FindTreatment(treatmentName string) ([]domain.Treatment, error)
	BookPatient(
		patient domain.Patient,
		treatment domain.Treatment,
		startTime schedule.JaneTime,
	) error
}

type bookingState struct {
	logger    logger.Logger
	writer    terminal.ScreenWriter
	fetcher   bookingDataFetcher
	rootState terminal.State

	currentBuffer string
	nextState     terminal.State
	booking       bookingBuilder
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	fetcher bookingDataFetcher,
	rootState terminal.State,
) terminal.State {
	return &bookingState{
		logger:    logger,
		writer:    writer,
		fetcher:   fetcher,
		rootState: rootState,
	}
}

func (bookingState) Name() string {
	return "booking"
}

func (state *bookingState) Initialize() {
	state.logger.Debugf(
		"entering booking. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.booking = bookingBuilder{
		substate: argument,
	}
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *bookingState) triggerAutocomplete() {
}

func (state *bookingState) ClearBuffer() {
	state.currentBuffer = ""
	state.writer.NewLine()
	state.writer.WriteString("")
}
