package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state *bookingState) handleBooking(flags map[string]string) {
	missingFlags := map[string]string{
		bookingDateFlag: bookingDateFlag,
		treatmentFlag:   treatmentFlag,
		patientFlag:     patientFlag,
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
		flow:     booking,
	}

	// Setup Patient
	var err error
	patient, patients, err := util.ParsePatientValue(
		state.fetcher,
		flags[patientFlag],
	)
	if err != nil {
		state.buffer.AddString(err.Error())
		return
	}
	builder.patientSelector = interactive.NewPatientSelector(patient, patients)

	// Setup Treatment
	treatment, treatments, err := util.ParseTreatmentFlag(
		state.fetcher,
		flags[treatmentFlag],
	)
	if err != nil {
		state.buffer.AddString(err.Error())
		return
	}
	builder.treatmentSelector = interactive.NewTreatmentSelector(treatment, treatments)

	// Setup Date
	builder.appointmentDate, err = util.ParseDate(
		util.DateTimeFormat,
		util.YearDateTimeFormat,
		flags[bookingDateFlag],
	)
	if err != nil {
		state.buffer.AddString(err.Error())
		return
	}

	state.builder = builder
}
