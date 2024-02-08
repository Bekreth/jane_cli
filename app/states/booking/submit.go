package booking

const bookingDateFlag = "-d"
const treatmentFlag = "-t"
const patientFlag = "-p"
const cancelFlag = "-c"

const cancelCommand = "cancel"
const bookCommand = "book"
const backCommand = ".."

func (state *bookingState) Submit(flags map[string]string) bool {
	state.logger.Debugf("submitting query flags: %v", flags)
	if _, exists := flags[backCommand]; exists {
		state.nextState = state.rootState
		return true
	} else if _, exists := flags[cancelCommand]; exists {
		state.handleCancel(flags)
	} else {
		state.handleBooking(flags)
	}
	return true
}
