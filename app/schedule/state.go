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
	FetchSchedule(
		startDate schedule.JaneTime,
		endDate schedule.JaneTime,
	) (schedule.Schedule, error)
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

func (state *scheduleState) Initialize() {
	state.logger.Debugf(
		"entering schedule. available states %v",
		state.rootState.Name(),
	)
	state.currentBuffer = ""
	state.nextState = state
	state.writer.NewLine()
	state.writer.WriteString("")
}

func (state *scheduleState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	terminal.KeyHandler(
		key,
		&state.currentBuffer,
		state.triggerAutocomplete,
		state.submit,
	)

	if character != 0 {
		state.currentBuffer += string(character)
	}

	state.writer.WriteString(state.currentBuffer)
	return state.nextState
}

func (state *scheduleState) shutdown() {
	state.currentBuffer = ""
}

func (state *scheduleState) triggerAutocomplete() {
	words := strings.Split(state.currentBuffer, " ")

	for key := range autocompletes {
		if strings.HasPrefix(key, words[len(words)-1]) {
			arguments := append(words[0:len(words)-1], key)
			state.currentBuffer = strings.Join(arguments, " ")
		}
	}
}

func (state *scheduleState) ClearBuffer() {
	state.currentBuffer = ""
	state.writer.NewLine()
	state.writer.WriteString("")
}
