package schedule

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

const scheduleTimeYearFormat = "2006.01.02"
const scheduleTimeFormat = "01.02"
const dateFlag = "-d"

func (state *scheduleState) submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.buffer.Clear()

	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
		return
	} else if _, exists := flags["help"]; exists {
		state.printHelp()
		return
	}

	var startAt schedule.JaneTime
	var endAt schedule.JaneTime
	timeIsSet := false

	for key, times := range autocompletes {
		if _, exists := flags[key]; exists {
			startAt, endAt = times()
			timeIsSet = true
		}
	}

	if !timeIsSet {
		parsedTime, err := util.ParseDate(
			scheduleTimeFormat,
			scheduleTimeYearFormat,
			flags[dateFlag],
		)
		if err != nil {
			state.buffer.WriteStoreString(err.Error())
			return
		}
		startAt = parsedTime
		endAt = parsedTime
		timeIsSet = true
	}

	if timeIsSet {
		fetchedSchedule, err := state.fetcher.FetchSchedule(startAt, endAt)
		if err != nil {
			state.buffer.WriteStoreString(fmt.Sprintf("failed to get schedule: %v", err))
		}
		if len(fetchedSchedule.Appointments) == 0 {
			state.buffer.WriteStoreString(fmt.Sprintf(
				"no shift between %v and %v",
				startAt.Format(scheduleTimeFormat),
				endAt.Format(scheduleTimeFormat),
			))
		} else {
			state.buffer.WriteStoreString("\n" + fetchedSchedule.ToString())
		}
		return
	}
}

func (state *scheduleState) printHelp() {
	// TODO: automate this list of elements
	state.buffer.WriteStoreString(fmt.Sprintf(
		"schedule command takes a date arguemnt:\n%v",
		"\t-d\tMM.DD",
	))
}
