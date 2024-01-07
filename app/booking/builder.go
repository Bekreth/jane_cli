package booking

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
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

	patientSelector     interactive.Interactive[domain.Patient]
	treatmentSelector   interactive.Interactive[domain.Treatment]
	appointmentSelector interactive.Interactive[schedule.Appointment]

	appointmentDate schedule.JaneTime
	cancelMessage   string
}

func newBookingBuilder() bookingBuilder {
	return bookingBuilder{
		substate: argument,
		flow:     undefined,

		patientSelector:     interactive.EmptyPatientSelector(),
		treatmentSelector:   interactive.EmptyTreatmentSelector(),
		appointmentSelector: interactive.EmptyAppointmentSelector(),
	}
}

func (builder bookingBuilder) confirmationMessage() string {
	switch builder.flow {
	case booking:
		return fmt.Sprintf(
			"Book %v for a %v at %v? (Y/n)",
			builder.patientSelector.TargetSelection().PrintSelector(),
			builder.treatmentSelector.TargetSelection().PrintHeader(),
			builder.appointmentDate.HumanDateTime(),
		)
	case canceling:
		return fmt.Sprintf(
			"Cancel appointment with %v at %v? (Y/n)",
			builder.appointmentSelector.TargetSelection().Deref().Patient.PrintName(),
			builder.appointmentDate.HumanDateTime(),
		)
	}
	return ""
}
