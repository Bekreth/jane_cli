package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
)

const ESCAPE_STRING = "(or ESC to back out)"

func PrintSelector(interactive Interactive) string {
	output := []string{}
	for i, possible := range interactive.PossibleSelections() {
		if i == 0 {
			output = append(output, fmt.Sprintf("%v %v", possible.PrintHeader(), ESCAPE_STRING))
		}
		output = append(output, fmt.Sprintf("%v: %v", i+1, possible.PrintSelector()))
	}
	return strings.Join(output, "\n")
}

func ElementSelector(
	character rune,
	input []Selection,
) (Selection, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Input has size 0")
	}
	index, err := strconv.Atoi(string(character))
	if err != nil || index == 0 {
		return nil, fmt.Errorf(
			"selector value of '%v' unacceptable. select a value between 1 and %v",
			string(character),
			len(input),
		)
	}
	if index > len(input) {
		return nil, fmt.Errorf(
			"selector value of '%v' is too large.  select a value between 1 and %v",
			index,
			len(input),
		)
	}
	return input[index-1], nil
}

func HandleSelection(
	character rune,
	key keyboard.Key,
) (Selection, error) {
	return nil, nil
}

func mod(a, b int) int {
	return (a%b + b) % b
}
