package schedule

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
)

type Appointment struct {
	StartAt JaneTime       `json:"start_at"`
	EndAt   JaneTime       `json:"end_at"`
	State   string         `json:"state"`
	Patient domain.Patient `json:"patient"`
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
