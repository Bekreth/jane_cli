package states

import (
	terminal "github.com/bekreth/screen_reader_terminal"
	"github.com/eiannone/keyboard"
)

type State interface {
	Name() string
	Initialize() *terminal.Buffer
	HandleKeyinput(character rune, key keyboard.Key) (State, bool)
	Submit(map[string]string) bool
	HelpString() string
}
