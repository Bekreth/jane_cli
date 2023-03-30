package app

import (
	"fmt"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type scheduleState struct {
	logger    logger.Logger
	rootState state
}

func (scheduleState) name() string {
	return "schedule"
}

func (schedule scheduleState) initialize() {
	schedule.logger.Debugf(
		"entering schedule. available states %v",
		schedule.rootState.name(),
	)
}

func (schedule scheduleState) handleKeyinput(character rune, key keyboard.Key) state {
	if string(character) == "2" {
		schedule.logger.Debugf("leaving schedule")
		return schedule.rootState
	}
	fmt.Print(character)
	return schedule
}

func (scheduleState) submit() {}

func (scheduleState) shutdown() {
}
