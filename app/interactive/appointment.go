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

func (selection SelectedAppointment) Deref() schedule.Appointment {
	return selection.Appointment
}

func (selection SelectedAppointment) hasSelection() bool {
	return selection.Appointment != schedule.DefaultAppointment
}

func EmptyAppointmentSelector() Interactive[schedule.Appointment] {
	var output *selector[schedule.Appointment]
	return output
}

func NewAppointmentSelector(
	selected schedule.Appointment,
	possible []schedule.Appointment,
) Interactive[schedule.Appointment] {
	possiblePatients := make([]Selection[schedule.Appointment], len(possible))
	for i, selection := range possible {
		possiblePatients[i] = SelectedAppointment{selection}
	}
	return &selector[schedule.Appointment]{
		page:              0,
		possibleSelection: possiblePatients,
		selected:          SelectedAppointment{selected},
	}
}
