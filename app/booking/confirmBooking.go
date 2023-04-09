package booking

func (state *bookingState) confirmBooking(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		state.writer.WriteString("submitting booking")
		state.writer.NewLine()
		err := state.fetcher.BookPatient(
			state.booking.targetPatient,
			state.booking.targetTreatment,
			state.booking.appointmentDate,
		)
		if err != nil {
			state.writer.WriteStringf("failed to make appointment: %v", err)
			state.writer.NewLine()
		} else {
			state.writer.NewLine()
			state.nextState = state.rootState
		}
	case "n":
		fallthrough
	case "N":
		state.writer.WriteString("aborting booking")
		state.writer.NewLine()
		state.nextState = state.rootState
	default:
		state.writer.WriteStringf(
			"input of %v not support.  Confirm booking (Y/n)?",
			string(character),
		)
		state.writer.NewLine()
	}
}
