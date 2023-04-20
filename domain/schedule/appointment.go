package schedule

import (
	"fmt"
	"time"

	"github.com/Bekreth/jane_cli/domain"
)

type AppointmentType string

const (
	Booked      AppointmentType = "booked"
	Arrived     AppointmentType = "arrived"
	Break       AppointmentType = "break"
	Unscheduled AppointmentType = "unscheduled"
)

type Appointment struct {
	ID        int             `json:"id"`
	PatientID int             `json:"patient_id"`
	StartAt   JaneTime        `json:"start_at"`
	EndAt     JaneTime        `json:"end_at"`
	State     AppointmentType `json:"state"`
	Patient   domain.Patient  `json:"patient"`
}

func (appointment Appointment) HasPassed() bool {
	return time.Now().After(appointment.StartAt.Time)
}

func (appointment Appointment) ToString() string {
	output := fmt.Sprintf(
		"%v to %v: %v",
		appointment.StartAt.Format(hourMinuteFormat),
		appointment.EndAt.Format(hourMinuteFormat),
		appointment.State,
	)
	if appointment.State == "booked" || appointment.State == "arrived" {
		output += fmt.Sprintf(" with %v", appointment.Patient.PreferredFirstName)
	}
	return output
}

var DefaultAppointment = Appointment{}
