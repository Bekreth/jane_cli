package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state *bookingState) handleCancel(flags map[string]string) {
	missingFlags := map[string]string{
		bookingDateFlag: "",
		patientFlag:     "",
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
		cancelTimeFormatDay,
		cancelTimeFormatYear,
		flags[bookingDateFlag],
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}
	builder.appointmentDate = dateValue

	builder, err = state.parseAppointmentValue(flags[patientFlag], builder)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	state.booking = builder
}
