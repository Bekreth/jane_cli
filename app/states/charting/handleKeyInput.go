package charting

import (
	"fmt"

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
		var err error
		switch state.builder.flow {
		case read:
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if !state.builder.chartSelector.HasSelection() {
				if len(state.builder.chartSelector.PossibleSelections()) == 0 {
					state.builder.chartSelector, err = state.fetchCharts()
				}
				if err != nil {
					state.buffer.WriteStoreString(err.Error())
					state.builder = newChartingBuilder()
					state.buffer.PrintHeader()
				} else {
					if state.builder.chartSelector.HasSelection() {
						state.builder.substate = complete
					} else {
						state.builder.substate = chartSelector
					}
				}
			} else {
				state.builder.substate = complete
			}

		case create:
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if !state.builder.appointmentSelector.HasSelection() {
				if len(state.builder.appointmentSelector.PossibleSelections()) == 0 {
					state.builder.appointmentSelector, err = state.fetchAppointments()
				}
				if err != nil {
					state.buffer.WriteStoreString(err.Error())
					state.builder = newChartingBuilder()
					state.buffer.PrintHeader()
				} else {
					if state.builder.appointmentSelector.HasSelection() {
						state.builder.substate = noteEditor
					} else {
						state.builder.substate = appointmentSelector
					}
				}
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
			output := fmt.Sprintf(
				"Write chart notes for %v(or ESC to back out): ",
				state.builder.appointmentSelector.TargetSelection().PrintSelector(),
			)
			state.buffer.WriteStoreString(output)
		} else {
			state.buffer.WriteString(state.builder.noteUnderEdit)
		}

	case complete:
		//TODO: Fix this nonsense
		state.buffer.WriteStoreString(state.builder.chartSelector.TargetSelection().Deref().PrintText())
		state.buffer.PrintHeader()
		state.builder = newChartingBuilder()

	default:
	}

	return state.nextState
}
