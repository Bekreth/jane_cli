package cache

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

func (cache Cache) FindAppointments(
	startDate schedule.JaneTime,
	endDate schedule.JaneTime,
	patientName string,
) ([]schedule.Appointment, error) {
	possibleMatches := cache.appointmentsFromCache(startDate, endDate, patientName)
	if len(possibleMatches) > 0 {
		return possibleMatches, nil
	}
	cache.logger.Debugf("no appointments found in cache, looking up via Jane client")

	fetchedSchedule, err := cache.scheduleFetcher.FetchSchedule(startDate, endDate)
	if err != nil {
		return possibleMatches, err
	}

	for _, appointment := range fetchedSchedule.Appointments {
		thisAppointment := appointment
		thisAppointment.Patient = domain.Patient{}
		cache.appointments[appointment.ID] = thisAppointment
		patient := appointment.Patient
		patient.ID = appointment.PatientID
		cache.patients[appointment.PatientID] = patient
	}
	cache.logger.Debugf("searching for appointments again")
	possibleMatches = cache.appointmentsFromCache(startDate, endDate, patientName)

	return possibleMatches, err
}

func (cache Cache) appointmentsFromCache(
	startDate schedule.JaneTime,
	endDate schedule.JaneTime,
	patientName string,
) []schedule.Appointment {
	possibleMatches := []schedule.Appointment{}
	for _, appointment := range cache.appointments {
		thisAppointment := appointment
		thisAppointment.Patient = cache.patients[appointment.PatientID]

		if matchingAppointment(startDate, endDate, patientName, thisAppointment) {
			possibleMatches = append(possibleMatches, thisAppointment)
		}
	}
	return possibleMatches
}

func matchingAppointment(
	startDate schedule.JaneTime,
	endDate schedule.JaneTime,
	patientName string,
	appointment schedule.Appointment,
) bool {
	if appointment.State != schedule.Booked {
		return false
	}
	inTimeWindow := appointment.StartAt.After(startDate.Time) &&
		appointment.EndAt.Before(endDate.Time)
	if patientName == "" {
		return inTimeWindow
	} else {
		return inTimeWindow && matchingPatient(appointment.Patient, patientName)
	}
}
