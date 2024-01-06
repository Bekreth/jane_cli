package interactive

type Selection interface {
	GetID() int
	PrintHeader() string
	PrintSelector() string
}

type Interactive interface {
	SelectElement(character rune) error
	PossibleSelections() []Selection
	TargetSelection() Selection
	HasSelection() bool
}

type selector struct {
	page              int
	possibleSelection []Selection
	selected          Selection
}

func (selection *selector) SelectElement(character rune) error {
	pages := len(selection.possibleSelection) / 9
	var output error
	switch string(character) {
	case "f":
		fallthrough
	case "F":
		selection.page = (selection.page + 1) % pages
	case "b":
		fallthrough
	case "B":
		selection.page = mod((selection.page - 1), pages)
	default:
		selection.selected, output = ElementSelector(
			character,
			selection.PossibleSelections(),
		)
	}
	return output
}

func (selection *selector) PossibleSelections() []Selection {
	pages := len(selection.possibleSelection) / 9
	pageStart := selection.page * 9
	if selection.page < pages {
		pageEnd := (selection.page + 1) * 9
		return selection.possibleSelection[pageStart:pageEnd]
	} else {
		return selection.possibleSelection[pageStart:]
	}

}

func (selection *selector) TargetSelection() Selection {
	return selection.selected
}

func (selection *selector) HasSelection() bool {
	return selection.selected != nil
}
