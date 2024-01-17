package booking

import "fmt"

func (state *bookingState) confirmAction(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		switch state.builder.flow {
		case booking:
			state.confirmBooking()
		case canceling:
			state.confirmCancelation()
		}
		state.builder = newBookingBuilder()
		state.nextState = state.rootState

	case "n":
		fallthrough
	case "N":
		state.buffer.WriteStoreString("aborting")
		state.builder = newBookingBuilder()
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
		state.builder.patientSelector.TargetSelection().Deref(),
		state.builder.treatmentSelector.TargetSelection().Deref(),
		state.builder.appointmentDate,
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
		state.builder.appointmentSelector.TargetSelection().GetID(),
		state.builder.cancelMessage,
	)
	if err != nil {
		state.buffer.WriteStoreString(fmt.Sprintf("failed to cancel appointment: %v", err))
	} else {
		state.buffer.WriteStoreString("Successfully cancelled appointment")
		state.ClearBuffer()
	}
}