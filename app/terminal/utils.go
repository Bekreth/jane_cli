package terminal

import (
	"github.com/Bekreth/jane_cli/app/states"
	"github.com/eiannone/keyboard"
)

func indiciesOfChar(input string, char rune) []int {
	output := make([]int, 0)
	for i, c := range input {
		if c == char {
			output = append(output, i)
		}
	}
	return output
}

func MapKeysString(input map[string]string) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}

func MapKeys(input map[string]states.State) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}

func KeyHandler(
	key keyboard.Key,
	buffer *Buffer,
	triggerAutocomplete func(),
	submit func(),
) {
	switch key {
	case keyboard.KeySpace:
		buffer.AddCharacter(rune(' '))

	case keyboard.KeyTab:
		triggerAutocomplete()

	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		buffer.RemoveCharacter()

	case keyboard.KeyEnter:
		submit()

	case keyboard.KeyArrowUp:
	case keyboard.KeyArrowDown:
	case keyboard.KeyArrowLeft:
		buffer.MoveLeft()
	case keyboard.KeyArrowRight:
		buffer.MoveRight()

	default:
		//TODO: Do I need generalized catch here?
	}
}
