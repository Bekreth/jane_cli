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
		state.buffer.AddString("aborting")
		state.builder = newBookingBuilder()
		state.nextState = state.rootState

	default:
		state.buffer.AddString(fmt.Sprintf(
			"input of %v not support. Confirm or deny (Y/n)?",
			string(character),
		))
	}
}

func (state *bookingState) confirmBooking() {
	state.buffer.AddString("submitting booking\n")
	err := state.fetcher.BookPatient(
		state.builder.patientSelector.TargetSelection().Deref(),
		state.builder.treatmentSelector.TargetSelection().Deref(),
		state.builder.appointmentDate,
	)
	if err != nil {
		state.buffer.AddString(fmt.Sprintf("failed to make appointment: %v", err))
	} else {
		state.buffer.AddString("Successfully created booking")
	}
}

func (state *bookingState) confirmCancelation() {
	state.buffer.AddString("canceling appointment")
	err := state.fetcher.CancelAppointment(
		state.builder.appointmentSelector.TargetSelection().GetID(),
		state.builder.cancelMessage,
	)
	if err != nil {
		state.buffer.AddString(fmt.Sprintf("failed to cancel appointment: %v", err))
	} else {
		state.buffer.AddString("Successfully cancelled appointment")
	}
}
