package interactive

import (
	"github.com/Bekreth/jane_cli/domain"
)

type SelectedTreatment struct {
	domain.Treatment
}

func (selection SelectedTreatment) GetID() int {
	return selection.ID
}

func (SelectedTreatment) PrintHeader() string {
	return "Select intended treatment"
}

func (selection SelectedTreatment) PrintSelector() string {
	return selection.Name
}

func (selection SelectedTreatment) Deref() domain.Treatment {
	return selection.Treatment
}

func NewTreatmentSelector(
	selected domain.Treatment,
	possible []domain.Treatment,
) Interactive[domain.Treatment] {
	var selectedTreatment SelectedTreatment
	if selected != domain.DefaultTreatment {
		selectedTreatment = SelectedTreatment{selected}
	}
	possiblePatients := make([]Selection[domain.Treatment], len(possible))
	for i, selection := range possible {
		possiblePatients[i] = SelectedTreatment{selection}
	}
	return &selector[domain.Treatment]{
		page:              0,
		possibleSelection: possiblePatients,
		selected:          selectedTreatment,
	}
}
