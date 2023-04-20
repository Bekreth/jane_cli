package booking

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

type substate = string

type processFlow = string

const (
	unknown substate = "unknown"

	argument            substate = "arguemnt"
	actionConfirmation  substate = "confirm"
	patientSelector     substate = "patientSelector"
	treatmentSelector   substate = "treatmentSelector"
	appointmentSelector substate = "appointmentSelector"
)

const (
	undefined processFlow = "undefined"

	booking   processFlow = "booking"
	canceling processFlow = "canceling"
)

type bookingBuilder struct {
	substate substate
	flow     processFlow

	patients          []domain.Patient
	targetPatient     domain.Patient
	treatments        []domain.Treatment
	targetTreatment   domain.Treatment
	appointments      []schedule.Appointment
	targetAppointment schedule.Appointment

	appointmentDate schedule.JaneTime
	cancelMessage   string
}

func (builder bookingBuilder) confirmationMessage() string {
	switch builder.flow {
	case booking:
		return fmt.Sprintf(
			"Book %v %v for a %v at %v? (Y/n)",
			builder.targetPatient.PreferredFirstName,
			builder.targetPatient.LastName,
			builder.targetTreatment.Name,
			builder.appointmentDate.HumanDateTime(),
		)
	case canceling:
		return fmt.Sprintf(
			"Cancel appointment with %v %v at %v? (Y/n)",
			builder.targetAppointment.Patient.PreferredFirstName,
			builder.targetAppointment.Patient.LastName,
			builder.appointmentDate.HumanDateTime(),
		)
	}
	return ""
}
