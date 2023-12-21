package charting

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/eiannone/keyboard"
)

func (state *chartingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	switch state.builder.substate {
	case patientSelector:
		state.builder.targetPatient = util.ElementSelector(
			character,
			state.builder.patients,
			state.buffer,
		)
	case chartSelector:
		state.builder.targetChart = util.ElementSelector(
			character,
			state.builder.charts,
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

	if state.builder.substate != argument {
		switch state.builder.flow {
		case read:
			if state.builder.targetPatient == domain.DefaultPatient {
				state.builder.substate = patientSelector
			} else if state.builder.targetChart.ID == 0 {
				state.fetchCharts()
				state.builder.substate = chartSelector
			} else {
				state.buffer.WriteStoreString(state.builder.targetChart.Snippet)
				state.builder.substate = unknown
				state.builder.flow = undefined
			}
		case create:
			panic("unimplemented!")
		}
	}

	switch state.builder.substate {
	case patientSelector:
		patientList := []string{"Select intended patient"}
		for i, patient := range state.builder.patients {
			patientList = append(
				patientList,
				fmt.Sprintf("%v: %v %v", i+1, patient.FirstName, patient.LastName),
			)
		}
		state.buffer.WriteStoreString(strings.Join(patientList, "\n"))
	case chartSelector:
		chartList := []string{fmt.Sprintf(
			"Select desired chart for %v",
			state.builder.targetPatient.PrintName(),
		)}
		state.logger.Debugf("Total charts: %v", len(state.builder.charts))
		for i, chart := range state.builder.charts {
			chartList = append(
				chartList,
				fmt.Sprintf("%v: %v", i+1, chart.EnteredOn.HumanDate()),
			)
		}
		state.buffer.WriteStoreString(strings.Join(chartList, "\n"))
	default:
	}

	return state.nextState
}
