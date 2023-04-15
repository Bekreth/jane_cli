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
	include     map[AppointmentType]interface{}
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
				State:   Unscheduled,
			})
		}
		updatedAppointments = append(updatedAppointments, appointment)
		timePointer = appointment.EndAt.Time
	}

	if timePointer.Sub(pair.shift.EndAt.Time).Abs() > time.Minute {
		updatedAppointments = append(updatedAppointments, Appointment{
			StartAt: JaneTime{timePointer},
			EndAt:   pair.shift.EndAt,
			State:   Unscheduled,
		})
	}

	appointmentString := []string{}
	for _, appointment := range updatedAppointments {
		if _, exists := pair.include[appointment.State]; exists {
			appointmentString = append(appointmentString, appointment.ToString())
		}
	}

	return fmt.Sprintf(
		"shift on %v from %v to %v\n%v",
		pair.shift.StartAt.Format(humanDateFormat),
		pair.shift.StartAt.Format(hourMinuteFormat),
		pair.shift.EndAt.Format(hourMinuteFormat),
		strings.Join(appointmentString, "\n"),
	)
}
