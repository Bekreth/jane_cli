package schedule

import (
	"fmt"
)

// state: booked, break

type Appointment struct {
	StartAt JaneTime `json:"start_at"`
	EndAt   JaneTime `json:"end_at"`
	State   string   `json:"state"`
	Patient Patient  `json:"patient"`
}

func (appointment Appointment) ToString() string {
	output := fmt.Sprintf(
		"* %v from %v to %v",
		appointment.State,
		appointment.StartAt.Format(hourMinuteFormat),
		appointment.EndAt.Format(hourMinuteFormat),
	)
	if appointment.State == "booked" {
		output += fmt.Sprintf(" with %v", appointment.Patient.PreferredFirstName)
	}
	return output
}
