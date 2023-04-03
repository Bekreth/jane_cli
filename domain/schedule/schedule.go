package schedule

import (
	"strings"
)

type Schedule struct {
	Appointments []Appointment `json:"appointments"`
	Shifts       []Shift       `json:"shifts"`
}

func (schedule Schedule) ToString() string {
	pairings := map[int]pairedShiftAppointment{}

	for _, shift := range schedule.Shifts {
		pairings[shift.StartAt.Day()] = pairedShiftAppointment{
			shift: shift,
		}
	}

	for _, appointment := range schedule.Appointments {
		pair := pairings[appointment.StartAt.Day()]
		updatedAppointment := append(pair.appointment, appointment)
		pairings[appointment.StartAt.Day()] = pairedShiftAppointment{
			shift:       pair.shift,
			appointment: updatedAppointment,
		}
	}

	pairingStrings := []string{}
	for _, pairing := range pairings {
		pairingStrings = append(pairingStrings, pairing.ToString())
	}

	return strings.Join(pairingStrings, "\n\n")
}

type Patient struct {
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PreferredFirstName string `json:"preferred_first_name"`
}

type Shift struct {
	StartAt JaneTime `json:"start_at"`
	EndAt   JaneTime `json:"end_at"`
}
