package booking

import "fmt"

// TODO: Add the subcommands
func (state *bookingState) printHelp() {
	// TODO: automate this list of elements
	state.buffer.WriteStoreString(fmt.Sprintf(
		"booking command is used to create new bookings:\n%v\n%v\n%v",
		"\t-d\tWhen to creating the appointment in the format of MM.DDTHH.MM",
		"\t-t\tThe treatment to use",
		"\t-p\tThe name of the patient (First, last, or preffered)",
	))
}
