package charting

import (
	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/eiannone/keyboard"
)

func (state *chartingState) isInteractive() bool {
	substate := state.builder.substate

	return substate == actionConfirmation ||
		substate == patientSelector ||
		substate == chartSelector ||
		substate == noteEditor
}

func (state *chartingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	if key == keyboard.KeyEsc && state.isInteractive() {
		state.builder = newChartingBuilder()
		state.buffer.PrintHeader()
		return state.nextState
	}

	switch state.builder.substate {
	case actionConfirmation:
		state.confirmAction(character)

	case patientSelector:
		state.builder.patientSelector.SelectElement(character)

	case chartSelector:
		state.builder.chartSelector.SelectElement(character)

	case appointmentSelector:
		state.builder.appointmentSelector.SelectElement(character)

	case noteEditor:
		//TODO: Limit to standard keys
		switch key {
		case keyboard.KeyDelete:
			fallthrough
		case keyboard.KeyBackspace2:
			fallthrough
		case keyboard.KeyBackspace:
			currentNote := state.builder.noteUnderEdit
			state.builder.noteUnderEdit = currentNote[:len(currentNote)-1]
		case keyboard.KeyEnter:
			state.builder.note = state.builder.noteUnderEdit
		case keyboard.KeySpace:
			state.builder.noteUnderEdit = state.builder.noteUnderEdit + " "
		default:
			state.builder.noteUnderEdit = state.builder.noteUnderEdit + string(character)
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

	if state.builder.substate != argument {
		switch state.builder.flow {
		case read:
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if !state.builder.chartSelector.HasSelection() {
				state.fetchCharts()
				state.builder.substate = chartSelector
			} else {
				state.buffer.WriteStoreString(state.builder.chartSelector.TargetSelection().Deref().Snippet)
				state.buffer.PrintHeader()
				state.builder = newChartingBuilder()
			}

		case create:
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if !state.builder.appointmentSelector.HasSelection() {
				state.fetchAppointments()
				state.builder.substate = appointmentSelector
			} else if state.builder.note == "" {
				state.builder.substate = noteEditor
			} else {
				state.builder.substate = actionConfirmation
			}
		}
	}

	switch state.builder.substate {
	case actionConfirmation:
		state.buffer.WriteStoreString(state.builder.confirmationMessage())

	case patientSelector:
		state.buffer.WriteStoreString(
			interactive.PrintSelectorList(state.builder.patientSelector),
		)

	case chartSelector:
		state.buffer.WriteStoreString(
			interactive.PrintSelectorList(state.builder.chartSelector),
		)

	case appointmentSelector:
		state.buffer.WriteStoreString(
			interactive.PrintSelectorList(state.builder.appointmentSelector),
		)

	case noteEditor:
		if state.builder.noteUnderEdit == "" {
			state.buffer.WriteStoreString("Write chart notes (or ESC to back out): ")
		} else {
			state.buffer.WriteString(state.builder.noteUnderEdit)
		}

	default:
	}

	return state.nextState
}
