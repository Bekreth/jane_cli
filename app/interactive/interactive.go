package interactive

type Selection[R interface{}] interface {
	hasSelection() bool
	GetID() int
	PrintHeader() string
	PrintSelector() string
	Deref() R
}

type Interactive[R interface{}] interface {
	SelectElement(character rune) error
	PossibleSelections() []Selection[R]
	TargetSelection() Selection[R]
	HasSelection() bool

	// PagingInfo specifies what the current page is, how many pages there are, and
	// how many total elements there are
	PagingInfo() (int, int, int)
}

type selector[R interface{}] struct {
	page              int
	possibleSelection []Selection[R]
	selected          Selection[R]
}

func (selection *selector[R]) SelectElement(character rune) error {
	pages := (len(selection.possibleSelection) / 9) + 1
	var output error
	switch character {
	case 'f':
		fallthrough
	case 'F':
		selection.page = (selection.page + 1) % pages
	case 'b':
		fallthrough
	case 'B':
		selection.page = mod((selection.page - 1), pages)
	default:
		selection.selected, output = ElementSelector(
			character,
			selection.PossibleSelections(),
		)
	}
	return output
}

func (selection *selector[R]) PossibleSelections() []Selection[R] {
	if selection == nil {
		return []Selection[R]{}
	}
	pages := len(selection.possibleSelection) / 9
	pageStart := selection.page * 9
	if selection.page < pages {
		pageEnd := (selection.page + 1) * 9
		return selection.possibleSelection[pageStart:pageEnd]
	} else {
		return selection.possibleSelection[pageStart:]
	}

}

func (selection *selector[R]) TargetSelection() Selection[R] {
	return selection.selected
}

func (selection *selector[R]) HasSelection() bool {
	if selection == nil {
		return false
	}
	return selection.selected.hasSelection()
}

func (selection *selector[R]) PagingInfo() (int, int, int) {
	if selection == nil {
		return 0, 0, 0
	}
	return selection.page,
		len(selection.possibleSelection) / 9,
		len(selection.possibleSelection)
}
