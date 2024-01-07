package util

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain/schedule"
)

type AppointmentFetcher interface {
	FindAppointments(
		startDate schedule.JaneTime,
		endDate schedule.JaneTime,
		patientName string,
	) ([]schedule.Appointment, error)
}

func ParseAppointmentFlag(
	fetcher AppointmentFetcher,
	dateTime schedule.JaneTime,
	patientName string,
) (schedule.Appointment, []schedule.Appointment, error) {
	appointment := schedule.Appointment{}
	appointments, err := fetcher.FindAppointments(
		dateTime.ThisDay(),
		dateTime.NextDay(),
		patientName,
	)
	if err != nil {
		return appointment, appointments, err
	}

	if len(appointments) == 0 {
		return appointment, appointments, fmt.Errorf("no appointments found")
	} else if len(appointments) == 1 {
		appointment = appointments[0]
	}
	return appointment, appointments, nil
}
