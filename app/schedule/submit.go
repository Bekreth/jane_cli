package schedule

import (
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

const scheduleTimeYearFormat = "2006.01.02"
const scheduleTimeFormat = "01.02"
const dateFlag = "-d"

func (state *scheduleState) submit() {
	flags := terminal.ParseFlags(state.currentBuffer)
	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
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
			state.currentBuffer = ""
			state.writer.WriteString(err.Error())
			state.writer.NewLine()
			return
		}
		startAt = parsedTime
		endAt = parsedTime
		timeIsSet = true
	}

	if timeIsSet {
		fetchedSchedule, err := state.fetcher.FetchSchedule(startAt, endAt)
		if err != nil {
			state.writer.WriteStringf("failed to get schedule: %v", err)
			state.writer.NewLine()
		}
		if len(fetchedSchedule.Appointments) == 0 {
			state.writer.WriteStringf(
				"no shift between %v and %v",
				startAt.Format(scheduleTimeFormat),
				endAt.Format(scheduleTimeFormat),
			)
			state.writer.NewLine()
		} else {
			state.writer.WriteString("\n" + fetchedSchedule.ToString())
		}
		state.currentBuffer = ""
		return
	}

	state.nextState = state.rootState
}
