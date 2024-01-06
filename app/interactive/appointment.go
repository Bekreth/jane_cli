package interactive

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain/schedule"
)

type SelectedAppointment struct {
	schedule.Appointment
}

func (selection SelectedAppointment) GetID() int {
	return selection.ID
}

func (SelectedAppointment) PrintHeader() string {
	return "Select intended appointment"
}

func (selection SelectedAppointment) PrintSelector() string {
	return fmt.Sprintf(
		"%v with %v",
		selection.StartAt.HumanDateTime(),
		selection.Patient.PrintName(),
	)
}

func NewAppointmentSelector(
	selected schedule.Appointment,
	possible []schedule.Appointment,
) Interactive {
	var selectedAppointment SelectedAppointment
	if selected != schedule.DefaultAppointment {
		selectedAppointment = SelectedAppointment{selected}
	}
	possiblePatients := make([]Selection, len(possible))
	for i, selection := range possible {
		possiblePatients[i] = SelectedAppointment{selection}
	}
	return &selector{
		page:              0,
		possibleSelection: possiblePatients,
		selected:          selectedAppointment,
	}
}
