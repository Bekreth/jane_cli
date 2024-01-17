package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state *bookingState) handleCancel(flags map[string]string) {
	missingFlags := map[string]string{
		bookingDateFlag: "",
	}

	for key := range missingFlags {
		delete(missingFlags, key)
	}
	if len(missingFlags) != 0 {
		joined := strings.Join(terminal.MapKeysString(missingFlags), ", ")
		notifcation := fmt.Sprintf("missing arguments %v", joined)
		state.buffer.WriteStoreString(notifcation)
		return
	}
	builder := bookingBuilder{
		substate: unknown,
		flow:     canceling,
	}

	dateValue, err := util.ParseDate(
		util.DateFormat,
		util.YearDateFormat,
		flags[bookingDateFlag],
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}
	builder.appointmentDate = dateValue

	appointment, appointments, err := util.ParseAppointmentFlag(
		state.fetcher,
		dateValue,
		flags[patientFlag],
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}
	builder.appointmentSelector = interactive.NewAppointmentSelector(
		appointment,
		appointments,
	)

	state.builder = builder
}
