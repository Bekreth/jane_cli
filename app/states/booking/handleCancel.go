package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state *bookingState) handleCancel(flags map[string]string) {
	missingFlags := map[string]string{
		bookingDateFlag: bookingDateFlag,
	}

	for key := range missingFlags {
		delete(missingFlags, key)
	}
	if len(missingFlags) != 0 {
		joined := strings.Join(util.MapKeysString(missingFlags), ", ")
		notifcation := fmt.Sprintf("missing arguments %v", joined)
		state.buffer.AddString(notifcation)
		return
	}
	builder := bookingBuilder{
		substate: unknown,
		flow:     canceling,
	}

	// Setup Date
	dateValue, err := util.ParseDate(
		util.DateFormat,
		util.YearDateFormat,
		flags[bookingDateFlag],
	)
	if err != nil {
		state.buffer.AddString(err.Error())
		return
	}
	builder.appointmentDate = dateValue

	appointment, appointments, err := util.ParseAppointmentFlag(
		state.fetcher,
		dateValue,
		flags[patientFlag],
	)
	if err != nil {
		state.buffer.AddString(err.Error())
		return
	}
	builder.appointmentSelector = interactive.NewAppointmentSelector(
		appointment,
		appointments,
	)

	state.builder = builder
}
