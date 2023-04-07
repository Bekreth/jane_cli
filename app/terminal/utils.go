package terminal

import (
	"github.com/eiannone/keyboard"
)

func MapKeysString(input map[string]string) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}

func MapKeys(input map[string]State) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}

func KeyHandler(
	key keyboard.Key,
	buffer *string,
	triggerAutocomplete func(),
	submit func(),
) {
	switch key {
	case keyboard.KeySpace:
		*buffer += string(" ")

	case keyboard.KeyTab:
		triggerAutocomplete()

	case keyboard.KeyDelete:
		fallthrough
	case keyboard.KeyBackspace2:
		fallthrough
	case keyboard.KeyBackspace:
		if len(*buffer) != 0 {
			currentValue := *buffer
			*buffer = currentValue[0 : len(currentValue)-1]
		}

	case keyboard.KeyEnter:
		submit()

	default:
		//TODO: Do I need generalized catch here?
	}
}
