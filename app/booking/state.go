package booking

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type bookingDataFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
	FindTreatment(treatmentName string) ([]domain.Treatment, error)
	// BookPatient(
	//
	//	patient domain.Patient,
	//	treatment domain.Treatment,
	//	startTime time.Time,
	//
	// ) error
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

func (booking *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	switch booking.booking.substate {
	case bookingConfirmation:
		booking.confirmBooking(character)
	case treatmentSelector:
		booking.booking.targetTreatment = booking.handleTreatmentSelector(character)
	case patientSelector:
		booking.booking.targetPatient = booking.handlePatientSelector(character)
	default:
		terminal.KeyHandler(
			key,
			&booking.currentBuffer,
			booking.triggerAutocomplete,
			booking.Submit,
		)
		if character != 0 {
			booking.currentBuffer += string(character)
		}
		booking.writer.WriteString(booking.currentBuffer)
	}

	if booking.booking.substate != argument {
		if booking.booking.targetPatient == domain.DefaultPatient {
			booking.booking.substate = patientSelector
		} else if booking.booking.targetTreatment == domain.DefaultTreatment {
			booking.booking.substate = treatmentSelector
		} else {
			booking.booking.substate = bookingConfirmation
		}
	}

	switch booking.booking.substate {
	case bookingConfirmation:
		confirmationString := fmt.Sprintf(
			"Book %v for a %v at %v? (Y/n)",
			booking.booking.targetPatient.PreferredFirstName,
			booking.booking.targetTreatment.Name,
			booking.booking.appointmentDate.Format(bookingTimeFormat),
		)
		booking.writer.WriteString(confirmationString)
	case treatmentSelector:
		treatmentList := "Select intended treatment\n"
		for i, treatment := range booking.booking.treatments {
			treatmentList += fmt.Sprintf("%v: %v \n", i+1, treatment.Name)
		}
		booking.writer.WriteString(treatmentList)
	case patientSelector:
		patientList := "Select intended patient\n"
		for i, patient := range booking.booking.patients {
			patientList += fmt.Sprintf("%v: %v %v \n", i+1, patient.FirstName, patient.LastName)
		}
		booking.writer.WriteString(patientList)
	default:
	}

	return booking.nextState
}
func (booking *bookingState) triggerAutocomplete() {
}
