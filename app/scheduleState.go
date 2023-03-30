package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/Bekreth/jane_cli/schedule"
	"github.com/eiannone/keyboard"
)

type scheduleFetcher interface {
	fetchSchedule(startDate time.Time, endDate time.Time) ([]schedule.Day, error)
}

type subcommand string

const (
	none    subcommand = "none"
	opening subcommand = "opening"
	book    subcommand = "book"
)

type scheduleState struct {
	logger        logger.Logger
	rootState     state
	subState      subcommand
	fetcher       scheduleFetcher
	currentBuffer string
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

func (schedule *scheduleState) handleKeyinput(character rune, key keyboard.Key) state {
	switch key {
	case keyboard.KeyTab:
		schedule.triggerAutocomplete()
		// TODO: replace
		fmt.Println("schedule>", schedule.currentBuffer)
		return schedule
	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		if len(schedule.currentBuffer) != 0 {
			schedule.currentBuffer = schedule.currentBuffer[0 : len(schedule.currentBuffer)-1]
		}
		// TODO: replace
		fmt.Println("schedule> ", schedule.currentBuffer)

	case keyboard.KeyEnter:
		return schedule.submit()
	}

	if character != 0 {
		schedule.currentBuffer += string(character)
		//TODO: replace
		fmt.Println("schedule> ", schedule.currentBuffer)
	}
	return schedule
}

func (schedule *scheduleState) shutdown() {
	schedule.currentBuffer = ""
}

func (schedule *scheduleState) triggerAutocomplete() {
	schedule.logger.Debugln("triggering autocomplete")
	words := strings.Split(schedule.currentBuffer, " ")
	if completed := schedule.autocompleteWord(words[len(words)-1]); completed != "" {
		updatedBuffer := strings.Join(append(words[0:len(words)-1], completed), " ")
		schedule.currentBuffer = updatedBuffer
	}
}

func (schedule *scheduleState) autocompleteWord(word string) string {
	if strings.HasPrefix("openings", word) {
		return "openings "
	}
	return ""
}

func (schedule *scheduleState) submit() state {
	if schedule.currentBuffer == ".." {
		return schedule.rootState
	}
	return schedule
}
