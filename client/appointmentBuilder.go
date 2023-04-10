package client

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

func (client Client) BookPatient(
	patient domain.Patient,
	treatment domain.Treatment,
	startTime schedule.JaneTime,
) error {
	endTime := schedule.JaneTime{
		Time: startTime.Add(treatment.ScheduledDuration.Duration),
	}
	appointment, err := client.CreateAppointment(
		startTime,
		endTime,
		false,
	)
	if err != nil {
		return err
	}
	return client.BookAppointment(
		appointment,
		treatment,
		patient,
	)
}
