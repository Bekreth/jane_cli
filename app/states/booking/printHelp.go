package booking

import (
	"strings"
)

func (state *bookingState) printHelp() {
	helpString := []string{
		"Booking command is used for managing your schedule",
		"Available subcommands",
		"\tbook\tCreate a new appointment.  Requires date, patient name, and treatment name.  Default if no subcommand provided",
		"\tcancel\t Cancels a plan, requires date, optionally accepts patient name",
		"Available flags",
		"\t-d\tDate.  For booking, format is MM.DDTHH:MM, for canceling format is MM.DD",
		"\t-t\tThe treatment to use",
		"\t-p\tThe name of the patient (First, last, or preffered)",
	}
	state.buffer.WriteStoreString(strings.Join(helpString, "\n"))
}
