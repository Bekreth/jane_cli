package booking

func (state *bookingState) confirmBooking(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		state.writer.WriteString("submitting booking")
		state.writer.NewLine()
		//TODO
		//err := booking.fetcher.BookPatient(
		//	booking.booking.targetPatient,
		//	booking.booking.targetTreatment,
		//	booking.booking.appointmentDate,
		//)
		//if err != nil {
		//	booking.writer.writeStringf("failed to make appointment: %v", err)
		//	booking.writer.newLine()
		//}
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
