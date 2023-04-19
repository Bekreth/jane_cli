package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/eiannone/keyboard"
)

func (state *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	switch state.booking.substate {
	case actionConfirmation:
		state.confirmAction(character)
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
	case appointmentSelector:
		state.booking.targetAppointment = elementSelector(
			character,
			state.booking.appointments,
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
		switch state.booking.flow {
		case booking:
			if state.booking.targetPatient == domain.DefaultPatient {
				state.booking.substate = patientSelector
			} else if state.booking.targetTreatment == domain.DefaultTreatment {
				state.booking.substate = treatmentSelector
			} else {
				state.booking.substate = actionConfirmation
			}
		case canceling:
			if state.booking.targetAppointment == schedule.DefaultAppointment {
				state.booking.substate = appointmentSelector
			} else {
				state.booking.substate = actionConfirmation
			}
		}
	}

	switch state.booking.substate {
	case actionConfirmation:
		state.buffer.WriteStoreString(state.booking.confirmationMessage())
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
	case appointmentSelector:
		appointmentList := []string{"Select intended appointment"}
		for i, appointment := range state.booking.appointments {
			appointmentList = append(
				appointmentList,
				fmt.Sprintf(
					"%v: %v with %v %v",
					i+1,
					appointment.StartAt.HumanDateTime(),
					appointment.Patient.PreferredFirstName,
					appointment.Patient.LastName,
				),
			)
		}
		state.buffer.WriteStoreString(strings.Join(appointmentList, "\n"))
	default:
	}

	return state.nextState
}
