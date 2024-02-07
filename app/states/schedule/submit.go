package schedule

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

const dateFlag = "-d"
const appointmentFlag = "-a"
const breakFlag = "-b"
const openFlag = "-o"
const showAllFlag = "-s"

func (state *scheduleState) Submit(flags map[string]string) bool {
	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
		return true
	}

	setIncludeFlags(flags)

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
			util.DateFormat,
			util.YearDateFormat,
			flags[dateFlag],
		)
		if err != nil {
			//TODO: Test
			state.buffer.AddString(err.Error())
			return true
		}
		startAt = parsedTime
		endAt = parsedTime
		timeIsSet = true
	}

	if timeIsSet {
		fetchedSchedule, err := state.fetcher.FetchSchedule(startAt, endAt)
		if err != nil {
			//TODO: Test
			state.buffer.AddString(fmt.Sprintf("failed to get schedule: %v", err))
		}
		if len(fetchedSchedule.Appointments) == 0 {
			//TODO: Test
			state.buffer.AddString(fmt.Sprintf(
				"no shift between %v and %v",
				startAt.Format(util.DateFormat),
				endAt.Format(util.DateFormat),
			))
		} else {
			if _, exists := flags[showAllFlag]; exists {
				state.logger.Debugf("show all appointments")
				fetchedSchedule = fetchedSchedule.ShowAll()
			}
			//TODO: Test
			output := fetchedSchedule.OnlyInclude(flagsToAppointmentFilters(flags)).ToString()
			state.buffer.AddString(output)
		}
		return true
	}
	return true
}

func setIncludeFlags(input map[string]string) {
	flagSet := false
	if _, exists := input[openFlag]; exists {
		flagSet = true
	}
	if _, exists := input[breakFlag]; exists {
		flagSet = true
	}
	if _, exists := input[appointmentFlag]; exists {
		flagSet = true
	}

	if !flagSet {
		input[openFlag] = string(schedule.Unscheduled)
		input[breakFlag] = string(schedule.Break)
		input[appointmentFlag] = string(schedule.Booked)
	}
}

func flagsToAppointmentFilters(flags map[string]string) []schedule.AppointmentType {
	output := []schedule.AppointmentType{}
	if _, exists := flags[openFlag]; exists {
		output = append(output, schedule.Unscheduled)
	}
	if _, exists := flags[breakFlag]; exists {
		output = append(output, schedule.Break)
	}
	if _, exists := flags[appointmentFlag]; exists {
		output = append(output, schedule.Booked)
		output = append(output, schedule.Arrived)
	}
	return output
}

func (state *scheduleState) HelpString() string {
	return strings.Join([]string{
		"",
		"available commands for schedule:",
		"\ttoday\t\tthe schedule for today which excluding appointments that have happened",
		"\ttomorrow\tthe schedulee for tomorrow",
		"available flags for schedule:",
		"\t-d\tspecific day to lookup the schedule for. Takes date as 'MM.DD'",
		"\t-o\tflag to include only the openings in the schedule",
		"\t-b\tflag to include only the breaks in the schedule",
		"\t-a\tflag to include only the appointments in the schedule",
		"\t-s\tshow all appointments, not just those that are upcomming",
	}, "\n")
}
