package schedule

import (
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type scheduleFetcher interface {
	FetchSchedule(startDate schedule.JaneTime, endDate schedule.JaneTime) (schedule.Schedule, error)
}

var oneDay = 24 * time.Hour

// TODO: Fill out more
var autocompletes = map[string]func() (schedule.JaneTime, schedule.JaneTime){
	"today": func() (schedule.JaneTime, schedule.JaneTime) {
		return schedule.JaneTime{Time: time.Now()}, schedule.JaneTime{Time: time.Now()}
	},
	"tomorrow": func() (schedule.JaneTime, schedule.JaneTime) {
		return schedule.JaneTime{Time: time.Now().AddDate(0, 0, 1)}, schedule.JaneTime{Time: time.Now().AddDate(0, 0, 1)}
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
	schedule.currentBuffer = ""
	schedule.nextState = schedule
	schedule.writer.NewLine()
	schedule.writer.WriteString("")
}

func (schedule *scheduleState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	terminal.KeyHandler(
		key,
		&schedule.currentBuffer,
		schedule.triggerAutocomplete,
		schedule.submit,
	)

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
