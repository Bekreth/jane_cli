package booking

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/eiannone/keyboard"
)

func (booking *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	switch booking.booking.substate {
	case bookingConfirmation:
		booking.confirmBooking(character)
	case treatmentSelector:
		booking.booking.targetTreatment = elementSelector[domain.Treatment](
			character,
			booking.booking.treatments,
			booking.writer,
		)
	case patientSelector:
		booking.booking.targetPatient = elementSelector[domain.Patient](
			character,
			booking.booking.patients,
			booking.writer,
		)
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
		booking.writer.WriteStringf(
			"Book %v for a %v at %v? (Y/n)",
			booking.booking.targetPatient.PreferredFirstName,
			booking.booking.targetTreatment.Name,
			booking.booking.appointmentDate.Format(bookingTimeFormat),
		)
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
