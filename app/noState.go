package app

import "github.com/eiannone/keyboard"

type state interface {
	name() string
	initialize()
	handleKeyinput(character rune, key keyboard.Key) state
	shutdown()
}

type noState struct{}

func (noState) name() string {
	return "noState"
}

func (noState) initialize() {}

func (self noState) handleKeyinput(character rune, key keyboard.Key) state {
	return self
}

func (noState) shutdown() {}

func (noState) submit() {}
