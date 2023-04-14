package booking

import "fmt"

func (state *bookingState) confirmBooking(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		state.buffer.WriteStoreString("submitting booking")
		err := state.fetcher.BookPatient(
			state.booking.targetPatient,
			state.booking.targetTreatment,
			state.booking.appointmentDate,
		)
		state.booking.substate = argument
		if err != nil {
			state.buffer.WriteStoreString(fmt.Sprintf("failed to make appointment: %v", err))
		} else {
			state.ClearBuffer()
		}
	case "n":
		fallthrough
	case "N":
		state.buffer.WriteStoreString("aborting booking")
		state.nextState = state.rootState
	default:
		state.buffer.WriteStoreString(fmt.Sprintf(
			"input of %v not support.  Confirm booking (Y/n)?",
			string(character),
		))
	}
}
