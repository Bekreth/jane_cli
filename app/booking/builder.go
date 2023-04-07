package booking

import (
	"time"

	"github.com/Bekreth/jane_cli/domain"
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
	appointmentDate time.Time
}
