package booking

import (
	"time"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
)

type bookingDataFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
	FindTreatment(treatmentName string) ([]domain.Treatment, error)
	BookPatient(
		patient domain.Patient,
		treatment domain.Treatment,
		startTime time.Time,
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

func (booking *bookingState) Initialize() {
	booking.logger.Debugf(
		"entering booking. available states %v",
		booking.rootState.Name(),
	)
	booking.nextState = booking
	booking.booking = bookingBuilder{
		substate: argument,
	}
	booking.writer.NewLine()
	booking.writer.WriteString("")
}

func (booking *bookingState) triggerAutocomplete() {
}
