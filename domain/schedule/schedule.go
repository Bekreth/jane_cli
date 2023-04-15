package schedule

import (
	"strings"
)

type Schedule struct {
	Appointments []Appointment `json:"appointments"`
	Shifts       []Shift       `json:"shifts"`
	include      map[AppointmentType]interface{}
}

func New() Schedule {
	return Schedule{
		include: make(map[AppointmentType]interface{}),
	}
}

func (schedule Schedule) OnlyInclude(include []AppointmentType) Schedule {
	output := schedule
	for _, appointmentType := range include {
		output.include[appointmentType] = struct{}{}
	}
	return output
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
			include:     schedule.include,
		}
	}

	pairingStrings := []string{}
	for _, pairing := range pairings {
		pairingStrings = append(pairingStrings, pairing.ToString())
	}

	return strings.Join(pairingStrings, "\n\n")
}

// TODO: Extract into its own file and have it deserialized
type Shift struct {
	StartAt JaneTime `json:"start_at"`
	EndAt   JaneTime `json:"end_at"`
}
