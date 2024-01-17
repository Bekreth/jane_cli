package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

func (state chartingState) fetchAppointments() (
	interactive.Interactive[schedule.Appointment],
	error,
) {
	appointments, err := state.fetcher.FindAppointments(
		state.builder.date,
		state.builder.date.NextDay(),
		state.builder.patientSelector.TargetSelection().PrintSelector(),
	)
	if err != nil {
		return nil, err
	}

	if len(appointments) == 0 {
		return nil, fmt.Errorf(
			"no appointments found for patient %v",
			state.builder.patientSelector.TargetSelection().PrintSelector(),
		)
	} else if len(appointments) == 1 {
		return interactive.NewAppointmentSelector(
			appointments[0],
			appointments,
		), nil
	} else {
		return interactive.NewAppointmentSelector(
			schedule.Appointment{},
			appointments,
		), nil
	}
}
