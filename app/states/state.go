package states

import "github.com/eiannone/keyboard"

type State interface {
	Name() string
	Initialize()
	HandleKeyinput(character rune, key keyboard.Key) State
	ClearBuffer()
	RepeatLastOutput()
}
