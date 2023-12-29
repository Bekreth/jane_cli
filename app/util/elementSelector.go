package util

import (
	"fmt"
	"strconv"
)

func ElementSelector[R interface{}](
	character rune,
	input []R,
) (*R, error) {
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
	return &input[index-1], nil
}
