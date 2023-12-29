package booking

import (
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
)

type bookingDataFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
	FindTreatment(treatmentName string) ([]domain.Treatment, error)
	FindAppointments(
		startDate schedule.JaneTime,
		endDate schedule.JaneTime,
		patientName string,
	) ([]schedule.Appointment, error)
	BookPatient(
		patient domain.Patient,
		treatment domain.Treatment,
		startTime schedule.JaneTime,
	) error
	CancelAppointment(appointmentID int, cancelMessage string) error
}

type bookingState struct {
	logger    logger.Logger
	fetcher   bookingDataFetcher
	builder   bookingBuilder
	rootState terminal.State

	nextState terminal.State
	buffer    *terminal.Buffer
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	fetcher bookingDataFetcher,
	rootState terminal.State,
) terminal.State {
	buffer := terminal.NewBuffer(writer)
	return &bookingState{
		logger:    logger,
		fetcher:   fetcher,
		rootState: rootState,
		buffer:    &buffer,
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
	state.builder = newBookingBuilder()
	state.nextState = state
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

var autocompletes = map[string]string{
	helpCommand:   "",
	cancelCommand: "",
	bookCommand:   "",
}

func (state *bookingState) triggerAutocomplete() {
	words := strings.Split(state.buffer.Read(), " ")
	for key := range autocompletes {
		if strings.HasPrefix(key, words[len(words)-1]) {
			arguments := append(words[0:len(words)-1], key)
			state.buffer.WriteString(strings.Join(arguments, " ") + " ")
		}
	}
}

func (state *bookingState) ClearBuffer() {
	state.buffer.Clear()
	state.buffer.PrintHeader()
}

func (state *bookingState) RepeatLastOutput() {
	state.buffer.WritePrevious()
}
