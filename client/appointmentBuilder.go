package client

import (
	"time"

	"github.com/Bekreth/jane_cli/domain"
)

func (client Client) BookPatient(
	patient domain.Patient,
	treatment domain.Treatment,
	startTime time.Time,
) error {
	appointmentID, err := client.CreateAppointment(
		startTime,
		startTime.Add(treatment.ScheduledDuration.Duration),
		false,
	)
	if err != nil {
		return err
	}
	return client.BookAppointment(appointmentID)
}
