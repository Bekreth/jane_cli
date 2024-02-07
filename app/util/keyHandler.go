package util

import (
	terminal "github.com/bekreth/screen_reader_terminal"
	"github.com/eiannone/keyboard"
)

func KeyHandler(
	key keyboard.Key,
	buffer *terminal.Buffer,
	triggerAutocomplete func(),
) {
	switch key {
	case keyboard.KeySpace:
		buffer.AddCharacter(' ')

	case keyboard.KeyTab:
		triggerAutocomplete()

	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		buffer.RemoveCharacter()

	default:
		//TODO: Do I need generalized catch here?
	}
}
