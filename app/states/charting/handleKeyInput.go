package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/eiannone/keyboard"
)

func (state *chartingState) isInteractive() bool {
	substate := state.builder.substate
	return substate == actionConfirmation ||
		substate == patientSelector ||
		substate == chartSelector ||
		substate == noteEditor
}

func (state *chartingState) setNoteEditor() {
	if state.builder.note == "" {
		state.buffer.Clear()
		state.builder.noteUnderEdit = ""
		state.builder.substate = noteEditor
	} else {
		state.builder.substate = actionConfirmation
	}
}

func (state *chartingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) (states.State, bool) {
	addNewLine := false
	if key == keyboard.KeyEsc && state.isInteractive() {
		state.builder = newChartingBuilder()
		state.buffer.AddString("Backing out of interactive action")
		state.buffer.SetPrefix("charting: ")
		return state.nextState, true
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
			state.buffer.RemoveCharacter()
		case keyboard.KeyEnter:
			state.builder.note = state.builder.noteUnderEdit
		case keyboard.KeySpace:
			state.buffer.AddCharacter(' ')
			state.builder.noteUnderEdit, _ = state.buffer.OutputWithoutPrefix()
		default:
			state.buffer.AddCharacter(character)
			state.builder.noteUnderEdit, _ = state.buffer.OutputWithoutPrefix()
		}

	default:
		util.KeyHandler(
			key,
			state.buffer,
			state.triggerAutocomplete,
		)
		if character != 0 {
			state.buffer.AddCharacter(character)
		}
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
					state.buffer.AddString(err.Error())
					state.builder = newChartingBuilder()
				} else {
					if state.builder.chartSelector.HasSelection() {
						state.builder.substate = complete
					} else {
						state.builder.substate = chartSelector
					}
				}
			} else {
				state.builder.substate = complete
				addNewLine = true
			}

		case create:
			if !state.builder.patientSelector.HasSelection() {
				state.builder.substate = patientSelector
			} else if !state.builder.appointmentSelector.HasSelection() {
				if len(state.builder.appointmentSelector.PossibleSelections()) == 0 {
					state.builder.appointmentSelector, err = state.fetchAppointments()
				}
				if err != nil {
					state.buffer.AddString(err.Error())
					state.builder = newChartingBuilder()
				} else {
					if state.builder.appointmentSelector.HasSelection() {
						state.setNoteEditor()
					} else {
						state.builder.substate = appointmentSelector
					}
				}
			} else {
				state.setNoteEditor()
			}
		}
	}

	switch state.builder.substate {
	case patientSelector:
		state.buffer.AddString(
			interactive.PrintSelectorList(state.builder.patientSelector) + "\n",
		)

	case chartSelector:
		state.buffer.AddString(
			interactive.PrintSelectorList(state.builder.chartSelector) + "\n",
		)

	case appointmentSelector:
		state.buffer.AddString(
			interactive.PrintSelectorList(state.builder.appointmentSelector) + "\n",
		)

	case noteEditor:
		if state.builder.noteUnderEdit == "" {
			output := fmt.Sprintf(
				"Write chart notes for %v(or ESC to back out): ",
				state.builder.appointmentSelector.TargetSelection().PrintSelector(),
			)
			state.buffer.SetPrefix(output)
		}

	case complete:
		state.buffer.AddString(
			state.builder.chartSelector.TargetSelection().Deref().PrintText(),
		)
		state.builder = newChartingBuilder()

	case actionConfirmation:
		state.buffer.SetPrefix("charting: ")
		state.buffer.AddString(
			state.builder.confirmationMessage() + "\n",
		)

	default:
	}

	return state.nextState, addNewLine
}
