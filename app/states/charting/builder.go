package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/charts"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

type substate = string
type processFlow = string

const (
	unknown substate = "unknown"

	argument            substate = "arguemnt"
	actionConfirmation  substate = "actionConfirmation"
	patientSelector     substate = "patientSelector "
	chartSelector       substate = "chartSelector"
	appointmentSelector substate = "appointmentSelector"
	noteEditor          substate = "noteEditor"
	complete            substate = "complete"
)

const (
	undefined processFlow = "undefined"

	read   processFlow = "read"
	create processFlow = "create"
)

type chartingBuilder struct {
	substate substate
	flow     processFlow

	date schedule.JaneTime
	note string

	noteUnderEdit       string
	patientSelector     interactive.Interactive[domain.Patient]
	chartSelector       interactive.Interactive[charts.ChartEntry]
	appointmentSelector interactive.Interactive[schedule.Appointment]
}

func newChartingBuilder() chartingBuilder {
	return chartingBuilder{
		substate: argument,
		flow:     undefined,

		patientSelector:     interactive.EmptyPatientSelector(),
		chartSelector:       interactive.EmptyChartSelector(),
		appointmentSelector: interactive.EmptyAppointmentSelector(),
	}
}

func (builder chartingBuilder) confirmationMessage() string {
	return fmt.Sprintf(
		"Would you like to sign the chart for appointment on %v with contents:\n%v\n(Y/n/E)",
		builder.appointmentSelector.TargetSelection().PrintSelector(),
		builder.note,
	)
}
