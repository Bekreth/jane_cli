package schedule

import (
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/app/flag"
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
	terminal "github.com/bekreth/screen_reader_terminal"
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
	fetcher   scheduleFetcher
	rootState states.State

	buffer    *terminal.Buffer
	nextState states.State
}

func NewState(
	logger logger.Logger,
	fetcher scheduleFetcher,
	rootState states.State,
) states.State {
	buffer := terminal.NewBuffer()
	return &scheduleState{
		logger:    logger,
		fetcher:   fetcher,
		rootState: rootState,
		buffer:    buffer.SetPrefix("schedule: "),
	}
}

func (scheduleState) Name() string {
	return "schedule"
}

func (state *scheduleState) Initialize() *terminal.Buffer {
	state.logger.Debugf(
		"entering schedule. available states %v",
		state.rootState.Name(),
	)
	state.nextState = state
	state.buffer.Clear()
	return state.buffer
}

func (state *scheduleState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) states.State {
	util.KeyHandler(key, state.buffer, state.triggerAutocomplete)
	if character != 0 {
		state.buffer.AddCharacter(character)
	}
	return state.nextState
}

func (state *scheduleState) triggerAutocomplete() {
	data, _ := state.buffer.Output()
	flags := flag.Parse(data)

	for key := range autocompletes {
		for flagKey := range flags {
			if strings.HasPrefix(key, flagKey) {
				state.buffer.AddString(strings.Replace(key, flagKey, "", 1))
			}
		}
	}
}
