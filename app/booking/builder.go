package booking

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

type substate = string

const (
	argument            substate = "arguemnt"
	unknown                      = "unknown"
	bookingConfirmation          = "bookingConfirmation"
	patientSelector              = "patientSelector"
	treatmentSelector            = "treatmentSelector"
)

type bookingBuilder struct {
	substate        substate
	patients        []domain.Patient
	targetPatient   domain.Patient
	treatments      []domain.Treatment
	targetTreatment domain.Treatment
	appointmentDate schedule.JaneTime
}
