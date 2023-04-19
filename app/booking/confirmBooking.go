package booking

import "fmt"

func (state *bookingState) confirmAction(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		switch state.booking.flow {
		case booking:
			state.confirmBooking()
		case canceling:
			state.confirmCancelation()
		}
		state.booking.substate = argument
		state.booking.flow = undefined
	case "n":
		fallthrough
	case "N":
		state.buffer.WriteStoreString("aborting")
		state.nextState = state.rootState
	default:
		state.buffer.WriteStoreString(fmt.Sprintf(
			"input of %v not support. Confirm or deny (Y/n)?",
			string(character),
		))
	}
}

func (state *bookingState) confirmBooking() {
	state.buffer.WriteStoreString("submitting booking")
	err := state.fetcher.BookPatient(
		state.booking.targetPatient,
		state.booking.targetTreatment,
		state.booking.appointmentDate,
	)
	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to make appointment: %v", err))
	} else {
		state.buffer.WriteStoreString("Successfully created booking")
		state.ClearBuffer()
	}
}

func (state *bookingState) confirmCancelation() {
	state.buffer.WriteStoreString("canceling appointment")
	err := state.fetcher.CancelAppointment(
		state.booking.targetAppointment.ID,
		state.booking.cancelMessage,
	)
	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to cancel appointment: %v", err))
	} else {
		state.buffer.WriteStoreString("Successfully cancelled appointment")
		state.ClearBuffer()
	}
}
