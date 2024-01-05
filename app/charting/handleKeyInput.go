package charting

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/eiannone/keyboard"
)

func (state *chartingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	var selectorErr error
	switch state.builder.substate {
	case actionConfirmation:
		state.confirmAction(character)

	case patientSelector:
		if key == keyboard.KeyEsc {
			state.builder = newChartingBuilder()
			state.buffer.PrintHeader()
			return state.nextState
		}
		possiblePatient, err := util.ElementSelector(
			character,
			state.builder.patients,
		)
		selectorErr = err
		if err == nil {
			state.builder.targetPatient = *possiblePatient
		}

	case chartSelector:
		if key == keyboard.KeyEsc {
			state.builder = newChartingBuilder()
			state.buffer.PrintHeader()
			return state.nextState
		}
		possibleChart, err := util.ElementSelector(
			character,
			state.builder.charts,
		)
		selectorErr = err
		if err == nil {
			state.builder.targetChart = *possibleChart
		}

	case appointmentSelector:
		if key == keyboard.KeyEsc {
			state.builder = newChartingBuilder()
			state.buffer.PrintHeader()
			return state.nextState
		}
		possibleAppointment, err := util.ElementSelector(
			character,
			state.builder.appointments,
		)
		selectorErr = err
		if err == nil {
			state.builder.targetAppointment = *possibleAppointment
		}

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

	if selectorErr != nil {
		state.buffer.WriteStoreString(selectorErr.Error())
		state.builder = newChartingBuilder()
		return state.nextState
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
				state.buffer.PrintHeader()
				state.builder.substate = unknown
				state.builder.flow = undefined
			}
		case create:
			if state.builder.targetPatient == domain.DefaultPatient {
				state.builder.substate = patientSelector
			} else if state.builder.targetAppointment == schedule.DefaultAppointment {
				//TODO: Handle Error
				var err error
				state.builder.appointments, err = state.fetcher.FindAppointments(
					state.builder.date,
					state.builder.date.NextDay(),
					state.builder.targetPatient.FirstName,
				)
				if err != nil {
					state.buffer.WriteStoreString(err.Error())
					state.builder.substate = argument
				} else {
					state.builder.substate = appointmentSelector
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
		patientList := []string{"Select intended patient (or ESC to back out)"}
		for i, patient := range state.builder.patients {
			patientList = append(
				patientList,
				fmt.Sprintf("%v: %v %v", i+1, patient.FirstName, patient.LastName),
			)
		}
		state.buffer.WriteStoreString(strings.Join(patientList, "\n"))

	case chartSelector:
		chartList := []string{fmt.Sprintf(
			"Select desired chart for %v (or ESC to back out)",
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

	case appointmentSelector:
		appointmentList := []string{"Select intended appointment (or ESC to back out)"}
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
	case noteEditor:
		if state.builder.noteUnderEdit == "" {
			state.buffer.WriteStoreString("Write chart notes: ")
		} else {
			state.buffer.WriteString(state.builder.noteUnderEdit)
		}
	default:
	}

	return state.nextState
}