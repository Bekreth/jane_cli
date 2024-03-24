package states

import (
	"github.com/bekreth/screen_reader_terminal/buffer"
	"github.com/eiannone/keyboard"
)

type State interface {
	Name() string
	Initialize() *buffer.Buffer
	HandleKeyinput(character rune, key keyboard.Key) (State, bool)
	Submit(map[string]string) bool
	HelpString() string
}
