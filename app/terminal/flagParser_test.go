package terminal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	trials := []struct {
		description    string
		input          string
		expectedOutput map[string]string
	}{
		{
			description: "Empty input, empty output",
			input:       "",
			expectedOutput: map[string]string{
				"": "",
			},
		},
		{
			description: "No flags",
			input:       "these aren't flags",
			expectedOutput: map[string]string{
				"these":  "",
				"aren't": "",
				"flags":  "",
			},
		},
		{
			description: "only flags",
			input:       "-flag1 value1 -flag2 value2",
			expectedOutput: map[string]string{
				"-flag1": "value1",
				"-flag2": "value2",
			},
		},
		{
			description: "flags and values",
			input:       "other1 -flag1 value1 -flag2 value2",
			expectedOutput: map[string]string{
				"other1": "",
				"-flag1": "value1",
				"-flag2": "value2",
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOuptut := ParseFlags(trial.input)
			fmt.Println("actual", actualOuptut)
			assert.Equal(tt, trial.expectedOutput, actualOuptut)
		})
	}
}
