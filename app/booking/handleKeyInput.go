package booking

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/eiannone/keyboard"
)

func (state *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	switch state.booking.substate {
	case bookingConfirmation:
		state.confirmBooking(character)
	case treatmentSelector:
		state.booking.targetTreatment = elementSelector[domain.Treatment](
			character,
			state.booking.treatments,
			state.writer,
		)
	case patientSelector:
		state.booking.targetPatient = elementSelector[domain.Patient](
			character,
			state.booking.patients,
			state.writer,
		)
	default:
		terminal.KeyHandler(
			key,
			&state.currentBuffer,
			state.triggerAutocomplete,
			state.Submit,
		)
		if character != 0 {
			state.currentBuffer += string(character)
		}
		state.writer.WriteString(state.currentBuffer)
	}

	if state.booking.substate != argument {
		if state.booking.targetPatient == domain.DefaultPatient {
			state.booking.substate = patientSelector
		} else if state.booking.targetTreatment == domain.DefaultTreatment {
			state.booking.substate = treatmentSelector
		} else {
			state.booking.substate = bookingConfirmation
		}
	}

	switch state.booking.substate {
	case bookingConfirmation:
		state.writer.WriteStringf(
			"Book %v for a %v at %v? (Y/n)",
			state.booking.targetPatient.PreferredFirstName,
			state.booking.targetTreatment.Name,
			state.booking.appointmentDate.HumanDate(),
		)
	case treatmentSelector:
		treatmentList := "Select intended treatment\n"
		for i, treatment := range state.booking.treatments {
			treatmentList += fmt.Sprintf("%v: %v \n", i+1, treatment.Name)
		}
		state.writer.WriteString(treatmentList)
	case patientSelector:
		patientList := "Select intended patient\n"
		for i, patient := range state.booking.patients {
			patientList += fmt.Sprintf("%v: %v %v \n", i+1, patient.FirstName, patient.LastName)
		}
		state.writer.WriteString(patientList)
	default:
	}

	return state.nextState
}
