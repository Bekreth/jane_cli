package interactive

import (
	"github.com/Bekreth/jane_cli/domain"
)

type SelectedPatient struct {
	domain.Patient
}

func (selection SelectedPatient) GetID() int {
	return selection.ID
}

func (SelectedPatient) PrintHeader() string {
	return "Select intended patient"
}

func (selection SelectedPatient) PrintSelector() string {
	return selection.PrintName()
}

func (selection SelectedPatient) Deref() domain.Patient {
	return selection.Patient
}

func (selection SelectedPatient) hasSelection() bool {
	return selection.Patient != domain.DefaultPatient
}

func NewPatientSelector(
	selected domain.Patient,
	possible []domain.Patient,
) Interactive[domain.Patient] {
	possiblePatients := make([]Selection[domain.Patient], len(possible))
	for i, selection := range possible {
		possiblePatients[i] = SelectedPatient{selection}
	}
	return &selector[domain.Patient]{
		page:              0,
		possibleSelection: possiblePatients,
		selected:          SelectedPatient{selected},
	}
}
