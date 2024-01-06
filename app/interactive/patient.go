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

func NewPatientSelector(
	selected domain.Patient,
	possible []domain.Patient,
) Interactive {
	var selectedPatient SelectedPatient
	if selected != domain.DefaultPatient {
		selectedPatient = SelectedPatient{selected}
	}
	possiblePatients := make([]Selection, len(possible))
	for i, selection := range possible {
		possiblePatients[i] = SelectedPatient{selection}
	}
	return &selector{
		page:              0,
		possibleSelection: possiblePatients,
		selected:          selectedPatient,
	}
}
