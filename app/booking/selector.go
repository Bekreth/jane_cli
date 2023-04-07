package booking

import (
	"strconv"

	"github.com/Bekreth/jane_cli/domain"
)

// TODO: figure out generics
/*
func (state *bookingState) handleListSelector[T interface{}](
	character rune,
	logPrefix string,
	list []T,
) T {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		state.writer.WriteStringf(
			"selector value of %v unacceptable. select a value between 1 and %v",
			string(character),
			len(state.booking.treatments),
		)
		state.writer.NewLine()
	}
	if index > len(state.booking.treatments) {
		state.writer.WriteStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index,
			len(state.booking.treatments),
		)
		state.writer.NewLine()
	}
	state.logger.Debugf("selected treatment at index %v", index)
	return state.booking.treatments[index-1]
}
*/

func (state *bookingState) handleTreatmentSelector(character rune) domain.Treatment {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		state.writer.WriteStringf(
			"selector value of %v unacceptable. select a value between 1 and %v",
			string(character),
			len(state.booking.treatments),
		)
		state.writer.NewLine()
	}
	if index > len(state.booking.treatments) {
		state.writer.WriteStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index,
			len(state.booking.treatments),
		)
		state.writer.NewLine()
	}
	state.logger.Debugf("selected treatment at index %v", index)
	return state.booking.treatments[index-1]
}

func (state *bookingState) handlePatientSelector(character rune) domain.Patient {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		state.writer.WriteStringf(
			"selector value of %v unacceptable.  select a value between 1 and %v",
			string(character),
			len(state.booking.patients),
		)
		state.writer.NewLine()
	}
	if index > len(state.booking.patients) {
		state.writer.WriteStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index, len(state.booking.patients),
		)
		state.writer.NewLine()
	}
	state.logger.Debugf("selected patient at index %v", index)
	return state.booking.patients[index-1]
}
