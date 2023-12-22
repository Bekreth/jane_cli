package util

import (
	"fmt"
	"strconv"

	"github.com/Bekreth/jane_cli/app/terminal"
)

func ElementSelector[R interface{}](
	character rune,
	input []R,
	buffer *terminal.Buffer,
) R {
	//TODO: Handle errors _correctly_
	index, err := strconv.Atoi(string(character))
	if err != nil {
		buffer.WriteStoreString(fmt.Sprintf(
			"selector value of %v unacceptable. select a value between 1 and %v",
			string(character),
			len(input),
		))
	}
	if index > len(input) {
		buffer.WriteStoreString(fmt.Sprintf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index,
			len(input),
		))
	}
	return input[index-1]
}
