package schedule

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

const scheduleTimeYearFormat = "2006.01.02"
const scheduleTimeFormat = "01.02"
const dateFlag = "-d"
const appointmentFlag = "-a"
const breakFlag = "-b"
const openFlag = "-o"
const showAllFlag = "-s"

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
			if _, exists := flags[showAllFlag]; exists {
				state.logger.Debugf("show all appointments")
				fetchedSchedule = fetchedSchedule.ShowAll()
			}
			output := fetchedSchedule.OnlyInclude(flagsToAppointmentFilters(flags)).ToString()
			state.buffer.WriteStoreString("\n" + output)
		}
		return
	}
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

func (state *scheduleState) printHelp() {
	// TODO: automate this list of elements
	helpString := strings.Join([]string{
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
	state.buffer.WriteStoreString(helpString)
}
