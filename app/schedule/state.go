package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

const timeFormat = "2006-01-02"
const dateFlag = "-d"

type scheduleFetcher interface {
	FetchSchedule(startDate time.Time, endDate time.Time) (schedule.Schedule, error)
}

var oneDay = 24 * time.Hour

// TODO: Fill out more
var autocompletes = map[string]func() (time.Time, time.Time){
	"today": func() (time.Time, time.Time) {
		return time.Now().Local(), time.Now().Local()
	},
	"tomorrow": func() (time.Time, time.Time) {
		return time.Now().Local().Add(oneDay), time.Now().Local().Add(oneDay)
	},
}

type scheduleState struct {
	logger    logger.Logger
	writer    terminal.ScreenWriter
	fetcher   scheduleFetcher
	rootState terminal.State

	currentBuffer string
	nextState     terminal.State
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	fetcher scheduleFetcher,
	rootState terminal.State,
) terminal.State {
	return &scheduleState{
		logger:    logger,
		writer:    writer,
		fetcher:   fetcher,
		rootState: rootState,
	}
}

func (scheduleState) Name() string {
	return "schedule"
}

func (schedule *scheduleState) Initialize() {
	schedule.logger.Debugf(
		"entering schedule. available states %v",
		schedule.rootState.Name(),
	)
	schedule.nextState = schedule
	schedule.writer.NewLine()
	schedule.writer.WriteString("")
}

func (schedule *scheduleState) HandleKeyinput(character rune, key keyboard.Key) terminal.State {
	terminal.KeyHandler(key, &schedule.currentBuffer, schedule.triggerAutocomplete, schedule.submit)

	if character != 0 {
		schedule.currentBuffer += string(character)
	}

	schedule.writer.WriteString(schedule.currentBuffer)
	return schedule.nextState
}

func (schedule *scheduleState) shutdown() {
	schedule.currentBuffer = ""
}

func (schedule *scheduleState) triggerAutocomplete() {
	words := strings.Split(schedule.currentBuffer, " ")

	for key := range autocompletes {
		if strings.HasPrefix(key, words[len(words)-1]) {
			arguments := append(words[0:len(words)-1], key)
			schedule.currentBuffer = strings.Join(arguments, " ")
		}
	}
}

func (schedule *scheduleState) submit() {
	flags := terminal.ParseFlags(schedule.currentBuffer)
	if _, exists := flags[".."]; exists {
		schedule.nextState = schedule.rootState
		return
	}

	var startAt time.Time
	var endAt time.Time
	timeIsSet := false

	for key, times := range autocompletes {
		if _, exists := flags[key]; exists {
			startAt, endAt = times()
			timeIsSet = true
		}
	}

	if !timeIsSet {
		if dateString, exists := flags[dateFlag]; exists {
			now := time.Now().Local()
			year := now.Year()
			date, err := time.Parse(timeFormat, fmt.Sprintf("%v-%v", year, dateString))
			if err != nil {
				schedule.writer.WriteString("Failed to parse time.  Use format MM-DD")
				schedule.writer.NewLine()
				return
			}
			startAt = date
			endAt = date

			timeIsSet = true
		}
	}

	if timeIsSet {
		fetchedSchedule, err := schedule.fetcher.FetchSchedule(startAt, endAt)
		if err != nil {
			schedule.writer.WriteString(fmt.Sprintf("failed to get schedule: %v", err))
		}
		if len(fetchedSchedule.Appointments) == 0 {
			schedule.writer.WriteString(fmt.Sprintf(
				"no shift between %v and %v",
				startAt.Format(timeFormat),
				endAt.Format(timeFormat),
			))
			schedule.writer.NewLine()
		} else {
			schedule.writer.WriteString("\n" + fetchedSchedule.ToString())
		}
		schedule.currentBuffer = ""
		return
	}

	schedule.nextState = schedule.rootState
}
