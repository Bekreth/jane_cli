package interactive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectElement(t *testing.T) {
	trials := []struct {
		description      string
		input            rune
		inputSelector    *selector[testSelector]
		expectedErr      error
		expectedPage     int
		expectedSelected Selection[testSelector]
	}{
		{
			description:      "page increment on 'f'",
			input:            'f',
			inputSelector:    newSelector(0, 12, 0),
			expectedErr:      nil,
			expectedPage:     1,
			expectedSelected: nil,
		},
		{
			description:      "page increment on 'F'",
			input:            'F',
			inputSelector:    newSelector(1, 112, 0),
			expectedErr:      nil,
			expectedPage:     2,
			expectedSelected: nil,
		},
		{
			description:      "page decrement on 'b'",
			input:            'b',
			inputSelector:    newSelector(1, 12, 0),
			expectedErr:      nil,
			expectedPage:     0,
			expectedSelected: nil,
		},
		{
			description:      "page increment on 'B'",
			input:            'B',
			inputSelector:    newSelector(4, 112, 0),
			expectedErr:      nil,
			expectedPage:     3,
			expectedSelected: nil,
		},
		{
			description:      "Error from bad run",
			input:            'Z',
			inputSelector:    newSelector(4, 112, 0),
			expectedErr:      fmt.Errorf("selector value of 'Z' unacceptable. select a value between 1 and 9"),
			expectedPage:     4,
			expectedSelected: nil,
		},
		{
			description:      "successfully selected the correct value",
			input:            '3',
			inputSelector:    newSelector(4, 112, 0),
			expectedErr:      nil,
			expectedPage:     4,
			expectedSelected: testSelector{ID: 39},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualErr := trial.inputSelector.SelectElement(trial.input)

			assert.Equal(tt, trial.expectedErr, actualErr)
			assert.Equal(tt, trial.expectedPage, trial.inputSelector.page)
			assert.Equal(tt, trial.expectedSelected, trial.inputSelector.selected)
		})

	}
}
