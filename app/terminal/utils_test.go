package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndiciesOfChar(t *testing.T) {
	trials := []struct {
		description    string
		input          string
		expectedOutput []int
	}{
		{
			description:    "proper count",
			input:          "this is a happy test string",
			expectedOutput: []int{4, 7, 9, 15, 20},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			assert.Equal(tt, trial.expectedOutput, indiciesOfChar(trial.input, ' '))
		})
	}
}
