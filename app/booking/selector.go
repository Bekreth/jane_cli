package booking

import (
	"fmt"
	"strconv"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

func elementSelector[R domain.Treatment | domain.Patient | schedule.Appointment](
	character rune,
	input []R,
	buffer *terminal.Buffer,
) R {
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
