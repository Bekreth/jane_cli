package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state *bookingState) handleBooking(flags map[string]string) {
	missingFlags := map[string]string{
		bookingDateFlag: "",
		treatmentFlag:   "",
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
		flow:     booking,
	}

	var err error
	builder.targetPatient, builder.patients, err = util.ParsePatientValue(
		state.fetcher,
		flags[patientFlag],
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	builder, err = state.parseTreatmentValue(flags[treatmentFlag], builder)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	builder.appointmentDate, err = util.ParseDate(
		bookingTimeFormat,
		bookingTimeYearFormat,
		flags[bookingDateFlag],
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	state.builder = builder
}
