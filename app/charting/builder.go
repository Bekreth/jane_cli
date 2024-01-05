package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/charts"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

type substate = string
type processFlow = string

const (
	unknown substate = "unknown"

	argument            = "arguemnt"
	actionConfirmation  = "actionConfirmation"
	patientSelector     = "patientSelector "
	chartSelector       = "chartSelector"
	appointmentSelector = "appointmentSelector"
	noteEditor          = "noteEditor"
)

const (
	undefined processFlow = "undefined"

	read   = "read"
	create = "create"
)

type chartingBuilder struct {
	substate substate
	flow     processFlow

	date schedule.JaneTime
	note string

	noteUnderEdit     string
	patients          []domain.Patient
	targetPatient     domain.Patient
	charts            []charts.ChartEntry
	targetChart       charts.ChartEntry
	appointments      []schedule.Appointment
	targetAppointment schedule.Appointment
}

func newChartingBuilder() chartingBuilder {
	return chartingBuilder{
		substate: argument,
		flow:     undefined,
	}
}

func (builder chartingBuilder) confirmationMessage() string {
	return fmt.Sprintf(
		"Would you like to sign the chart for %v for appointment on %v with contents:\n%v\n(Y/n/E)",
		builder.targetPatient.PrintName(),
		builder.targetAppointment.StartAt.HumanDate(),
		builder.note,
	)
}
