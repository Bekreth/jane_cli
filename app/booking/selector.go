package booking

import (
	"strconv"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
)

func elementSelector[R domain.Treatment | domain.Patient](
	character rune,
	input []R,
	writer terminal.ScreenWriter,
) R {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		writer.WriteStringf(
			"selector value of %v unacceptable. select a value between 1 and %v",
			string(character),
			len(input),
		)
		writer.NewLine()
	}
	if index > len(input) {
		writer.WriteStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index,
			len(input),
		)
		writer.NewLine()
	}
	return input[index-1]
}
