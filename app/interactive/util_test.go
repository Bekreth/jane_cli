package interactive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSelector struct {
	ID int
}

func (selector testSelector) hasSelection() bool {
	return false
}

func (selector testSelector) GetID() int {
	return selector.ID
}

func (testSelector) PrintHeader() string {
	return "test selector header"
}

func (testSelector) PrintSelector() string {
	return "test selector selector"
}

func (selector testSelector) Deref() testSelector {
	return selector
}

func newTestInteractive(
	page int,
	elementCount int,
	selectedID int,
) Interactive[testSelector] {
	possibleSelection := make([]Selection[testSelector], elementCount)
	for i := 0; i < elementCount; i++ {
		possibleSelection[i] = testSelector{ID: i + 1}
	}
	return &selector[testSelector]{
		page:              page,
		possibleSelection: possibleSelection,
		selected:          nil,
	}
}

func TestPrintSelectorList(t *testing.T) {
	trials := []struct {
		description    string
		input          Interactive[testSelector]
		expectedOutput string
	}{
		{
			description:    "an empty interactive should return empty string",
			input:          newTestInteractive(0, 0, 0),
			expectedOutput: "",
		},
		{
			description: "0 page should only list header",
			input:       newTestInteractive(0, 5, 0),
			expectedOutput: fmt.Sprintf(
				"%v.  %v\n1: %v\n2: %v\n3: %v\n4: %v\n5: %v",
				testSelector{}.PrintHeader(),
				"(or ESC to back out)",
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
			),
		},
		{
			description: "0 page should only list header",
			input:       newTestInteractive(0, 10, 0),
			expectedOutput: fmt.Sprintf(
				"%v. %v %v\n1: %v\n2: %v\n3: %v\n4: %v\n5: %v\n6: %v\n7: %v\n8: %v\n9: %v",
				testSelector{}.PrintHeader(),
				"Showing page 1 of 2, 'f' to page forwards, 'b' to page backwards ",
				"(or ESC to back out)",
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
				testSelector{}.PrintSelector(),
			),
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := PrintSelectorList(trial.input)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}

}

func TestElementSelector(t *testing.T) {
	trials := []struct {
		description    string
		character      rune
		input          []Selection[testSelector]
		expectedOutput Selection[testSelector]
		expectedErr    error
	}{
		{
			description:    "No input, no input error",
			character:      rune('1'),
			input:          newTestInteractive(0, 0, 0).PossibleSelections(),
			expectedOutput: nil,
			expectedErr:    fmt.Errorf("Input has size 0"),
		},
		{
			description:    "Bad input",
			character:      rune('a'),
			input:          newTestInteractive(0, 3, 0).PossibleSelections(),
			expectedOutput: nil,
			expectedErr:    fmt.Errorf("selector value of 'a' unacceptable. select a value between 1 and 3"),
		},
		{
			description:    "Input is too large",
			character:      rune('5'),
			input:          newTestInteractive(0, 3, 0).PossibleSelections(),
			expectedOutput: nil,
			expectedErr:    fmt.Errorf("selector value of '5' is too large.  select a value between 1 and 3"),
		},
		{
			description:    "Get proper output",
			character:      rune('2'),
			input:          newTestInteractive(0, 3, 0).PossibleSelections(),
			expectedOutput: testSelector{ID: 2},
			expectedErr:    nil,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput, actualErr := ElementSelector(trial.character, trial.input)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
			assert.Equal(tt, trial.expectedErr, actualErr)
		})
	}
}
