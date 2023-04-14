package booking

import (
	"fmt"
	"strings"

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
		state.booking.targetTreatment = elementSelector(
			character,
			state.booking.treatments,
			state.buffer,
		)
	case patientSelector:
		state.booking.targetPatient = elementSelector(
			character,
			state.booking.patients,
			state.buffer,
		)
	default:
		terminal.KeyHandler(
			key,
			state.buffer,
			state.triggerAutocomplete,
			state.Submit,
		)
		state.buffer.AddCharacter(character)
		state.buffer.Write()
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
		state.buffer.WriteStoreString(fmt.Sprintf(
			"Book %v for a %v at %v? (Y/n)",
			state.booking.targetPatient.PreferredFirstName,
			state.booking.targetTreatment.Name,
			state.booking.appointmentDate.HumanDate(),
		))
	case treatmentSelector:
		treatmentList := []string{"Select intended treatment"}
		for i, treatment := range state.booking.treatments {
			treatmentList = append(
				treatmentList,
				fmt.Sprintf("%v: %v", i+1, treatment.Name),
			)
		}
		state.buffer.WriteStoreString(strings.Join(treatmentList, "\n"))
	case patientSelector:
		patientList := []string{"Select intended patient"}
		for i, patient := range state.booking.patients {
			patientList = append(
				patientList,
				fmt.Sprintf("%v: %v %v", i+1, patient.FirstName, patient.LastName),
			)
		}
		state.buffer.WriteStoreString(strings.Join(patientList, "\n"))
	default:
	}

	return state.nextState
}
