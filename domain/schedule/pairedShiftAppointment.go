package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type pairedShiftAppointment struct {
	shift       Shift
	appointment []Appointment
}

func (pair pairedShiftAppointment) ToString() string {
	sort.Slice(
		pair.appointment,
		func(i int, j int) bool {
			return pair.appointment[i].StartAt.Before(pair.appointment[j].StartAt.Time)
		})

	updatedAppointments := []Appointment{}
	timePointer := pair.shift.StartAt.Time
	for _, appointment := range pair.appointment {
		if appointment.StartAt.Sub(timePointer).Abs() > time.Minute {
			updatedAppointments = append(updatedAppointments, Appointment{
				StartAt: JaneTime{timePointer},
				EndAt:   appointment.StartAt,
				State:   "unscheduled",
			})
		}
		updatedAppointments = append(updatedAppointments, appointment)
		timePointer = appointment.EndAt.Time
	}

	if timePointer.Sub(pair.shift.EndAt.Time).Abs() > time.Minute {
		updatedAppointments = append(updatedAppointments, Appointment{
			StartAt: JaneTime{timePointer},
			EndAt:   pair.shift.EndAt,
			State:   "unscheduled",
		})
	}

	appointmentString := []string{}
	for _, appointment := range updatedAppointments {
		appointmentString = append(appointmentString, appointment.ToString())
	}

	return fmt.Sprintf(
		"shift from %v to %v\n%v",
		pair.shift.StartAt.Format(hourMinuteFormat),
		pair.shift.EndAt.Format(hourMinuteFormat),
		strings.Join(appointmentString, "\n"),
	)
}
