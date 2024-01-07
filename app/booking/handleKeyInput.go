package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/eiannone/keyboard"
)

func (state *bookingState) isInteractive() bool {
	substate := state.builder.substate

	return substate == actionConfirmation ||
		substate == patientSelector ||
		substate == treatmentSelector ||
		substate == appointmentSelector
}

func (state *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	if key == keyboard.KeyEsc && state.isInteractive() {
		state.builder = newBookingBuilder()
		state.buffer.PrintHeader()
		return state.nextState
	}

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
		state.builder.patientSelector.SelectElement(character)

	case appointmentSelector:
		state.builder.appointmentSelector.SelectElement(character)

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
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if state.builder.targetTreatment == domain.DefaultTreatment {
				state.builder.substate = treatmentSelector
			} else {
				state.builder.substate = actionConfirmation
			}
		case canceling:
			if !state.builder.appointmentSelector.HasSelection() {
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
		treatmentList := []string{"Select intended treatment (or ESC to back out)"}
		for i, treatment := range state.builder.treatments {
			treatmentList = append(
				treatmentList,
				fmt.Sprintf("%v: %v", i+1, treatment.Name),
			)
		}
		state.buffer.WriteStoreString(strings.Join(treatmentList, "\n"))

	case patientSelector:
		state.buffer.WriteStoreString(
			interactive.PrintSelector(state.builder.patientSelector),
		)

	case appointmentSelector:
		state.builder.appointmentSelector.SelectElement(character)

	default:
	}

	return state.nextState
}
