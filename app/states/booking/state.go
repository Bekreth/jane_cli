package booking

import (
	"strings"

	"github.com/Bekreth/jane_cli/app/flag"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/bekreth/screen_reader_terminal/buffer"
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
	rootState states.State

	nextState states.State
	buffer    *buffer.Buffer
}

func NewState(
	logger logger.Logger,
	fetcher bookingDataFetcher,
	rootState states.State,
) states.State {
	buffer := buffer.NewBuffer()
	return &bookingState{
		logger:    logger,
		fetcher:   fetcher,
		rootState: rootState,
		buffer:    buffer.SetPrefix("booking: "),
	}
}

func (bookingState) Name() string {
	return "booking"
}

func (state *bookingState) Initialize() *buffer.Buffer {
	state.logger.Debugf(
		"entering booking. available states %v",
		state.rootState.Name(),
	)
	state.builder = newBookingBuilder()
	state.nextState = state
	state.buffer.Clear()
	return state.buffer
}

var autocompletes = map[string]string{
	cancelCommand: "",
	bookCommand:   "",
}

func (state *bookingState) triggerAutocomplete() {
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
