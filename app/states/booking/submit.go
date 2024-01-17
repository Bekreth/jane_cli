package booking

import (
	"github.com/Bekreth/jane_cli/app/terminal"
)

const bookingDateFlag = "-d"
const treatmentFlag = "-t"
const patientFlag = "-p"
const cancelFlag = "-c"

const cancelCommand = "cancel"
const bookCommand = "book"
const helpCommand = "help"
const backCommand = ".."

func (state *bookingState) Submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.logger.Debugf("submitting query flags: %v", flags)
	state.buffer.Clear()
	if _, exists := flags[backCommand]; exists {
		state.nextState = state.rootState
		return
	} else if _, exists := flags[helpCommand]; exists {
		state.printHelp()
		return
	} else if _, exists := flags[cancelCommand]; exists {
		state.handleCancel(flags)
	} else {
		state.handleBooking(flags)
	}
}
