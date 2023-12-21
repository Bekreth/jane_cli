package charting

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/charts"
)

type substate = string
type processFlow = string

const (
	unknown substate = "unknown"

	argument           = "arguemnt"
	actionConfirmation = "actionConfirmation"
	patientSelector    = "patientSelector "
	chartFetching      = "chartFetching"
	chartSelector      = "chartSelector"
)

const (
	undefined processFlow = "undefined"

	read   = "read"
	create = "create"
)

type chartingBuilder struct {
	substate substate
	flow     processFlow

	patients      []domain.Patient
	targetPatient domain.Patient
	charts        []charts.ChartEntry
	targetChart   charts.ChartEntry
}
