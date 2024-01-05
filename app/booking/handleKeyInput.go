package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/eiannone/keyboard"
)

func (state *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	var selectorErr error
	switch state.builder.substate {
	case actionConfirmation:
		state.confirmAction(character)

	case treatmentSelector:
		possibleTreatment, err := util.ElementSelector(
			character,
			state.builder.treatments,
		)
		selectorErr = err
		if err == nil {
			state.builder.targetTreatment = *possibleTreatment
		}

	case patientSelector:
		possiblePatient, err := util.ElementSelector(
			character,
			state.builder.patients,
		)
		selectorErr = err
		if err == nil {
			state.builder.targetPatient = *possiblePatient
		}

	case appointmentSelector:
		possibleAppointment, err := util.ElementSelector(
			character,
			state.builder.appointments,
		)
		selectorErr = err
		if err == nil {
			state.builder.targetAppointment = *possibleAppointment
		}

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

	if selectorErr != nil {
		state.buffer.WriteStoreString(selectorErr.Error())
		state.builder = newBookingBuilder()
		return state.nextState
	}

	if state.builder.substate != argument {
		switch state.builder.flow {
		case booking:
			if state.builder.targetPatient == domain.DefaultPatient {
				state.builder.substate = patientSelector
			} else if state.builder.targetTreatment == domain.DefaultTreatment {
				state.builder.substate = treatmentSelector
			} else {
				state.builder.substate = actionConfirmation
			}
		case canceling:
			if state.builder.targetAppointment == schedule.DefaultAppointment {
				state.builder.substate = appointmentSelector
			} else {
				state.builder.substate = actionConfirmation
			}
		}
	}

	switch state.builder.substate {
	case actionConfirmation:
		state.buffer.WriteStoreString(state.builder.confirmationMessage())
	case treatmentSelector:
		treatmentList := []string{"Select intended treatment"}
		for i, treatment := range state.builder.treatments {
			treatmentList = append(
				treatmentList,
				fmt.Sprintf("%v: %v", i+1, treatment.Name),
			)
		}
		state.buffer.WriteStoreString(strings.Join(treatmentList, "\n"))
	case patientSelector:
		patientList := []string{"Select intended patient"}
		for i, patient := range state.builder.patients {
			patientList = append(
				patientList,
				fmt.Sprintf("%v: %v %v", i+1, patient.FirstName, patient.LastName),
			)
		}
		state.buffer.WriteStoreString(strings.Join(patientList, "\n"))
	case appointmentSelector:
		appointmentList := []string{"Select intended appointment"}
		for i, appointment := range state.builder.appointments {
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
