package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

func (state *chartingState) fetchAppointments() {
	appointments, err := state.fetcher.FindAppointments(
		state.builder.date,
		state.builder.date.NextDay(),
		state.builder.patientSelector.TargetSelection().PrintSelector(),
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	if len(appointments) == 0 {
		state.buffer.WriteStoreString(fmt.Sprintf(
			"no appointments found for patient %v",
			state.builder.patientSelector.TargetSelection().PrintSelector(),
		))
	} else if len(appointments) == 1 {
		state.builder.appointmentSelector = interactive.NewAppointmentSelector(
			appointments[0],
			appointments,
		)
	} else {
		state.builder.appointmentSelector = interactive.NewAppointmentSelector(
			schedule.Appointment{},
			appointments,
		)
	}
}
