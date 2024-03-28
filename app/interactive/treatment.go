package interactive

import (
	"sort"

	"github.com/Bekreth/jane_cli/app/util"
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

func (selection SelectedTreatment) hasSelection() bool {
	return selection.Treatment != domain.DefaultTreatment
}

func EmptyTreatmentSelector() Interactive[domain.Treatment] {
	var output *selector[domain.Treatment]
	return output
}

func NewTreatmentSelector(
	selected domain.Treatment,
	possible []domain.Treatment,
) Interactive[domain.Treatment] {
	possibleTreatments := make([]Selection[domain.Treatment], len(possible))
	sort.Sort(util.Treatments(possible))
	for i, selection := range possible {
		possibleTreatments[i] = SelectedTreatment{selection}
	}
	return &selector[domain.Treatment]{
		page:              0,
		possibleSelection: possibleTreatments,
		selected:          SelectedTreatment{selected},
	}
}
