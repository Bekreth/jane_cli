package booking

import (
	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/terminal"
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
) states.State {
	if key == keyboard.KeyEsc && state.isInteractive() {
		state.builder = newBookingBuilder()
		state.buffer.PrintHeader()
		return state.nextState
	}

	switch state.builder.substate {
	case actionConfirmation:
		state.confirmAction(character)

	case treatmentSelector:
		state.builder.treatmentSelector.SelectElement(character)

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

	if state.builder.substate != argument {
		switch state.builder.flow {
		case booking:
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if !state.builder.treatmentSelector.HasSelection() {
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
		state.buffer.WriteStoreString(
			interactive.PrintSelectorList(state.builder.treatmentSelector),
		)

	case patientSelector:
		state.buffer.WriteStoreString(
			interactive.PrintSelectorList(state.builder.patientSelector),
		)

	case appointmentSelector:
		state.builder.appointmentSelector.SelectElement(character)

	default:
	}

	return state.nextState
}
