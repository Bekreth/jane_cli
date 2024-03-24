package booking

import (
	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/util"
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
) (states.State, bool) {
	addNewLine := false
	if key == keyboard.KeyEsc && state.isInteractive() {
		state.builder = newBookingBuilder()
		state.buffer.AddString("Backing out of interactive action")
		return state.nextState, true
	}

	// Hand key push
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
		util.KeyHandler(key, state.buffer, state.triggerAutocomplete)
		if character != 0 {
			state.buffer.AddCharacter(character)
		}
	}

	// Determine and switch to next state
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

	// Print state of flow
	switch state.builder.substate {
	case treatmentSelector:
		state.buffer.AddString(
			interactive.PrintSelectorList(state.builder.treatmentSelector) + "\n",
		)

	case patientSelector:
		state.buffer.AddString(
			interactive.PrintSelectorList(state.builder.patientSelector) + "\n",
		)

	case appointmentSelector:
		state.buffer.AddString(
			interactive.PrintSelectorList(state.builder.appointmentSelector) + "\n",
		)

	case actionConfirmation:
		state.buffer.AddString(
			state.builder.confirmationMessage() + "\n",
		)

	default:
	}

	return state.nextState, addNewLine
}
