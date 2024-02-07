package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
			description: "flags without arguments",
			input:       "-flag1 -flag2",
			expectedOutput: map[string]string{
				"-flag1": "",
				"-flag2": "",
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
		{
			description: "flags and values inside quotes",
			input:       "other1 -flag1 \"value1\" -flag2 \"value2\"",
			expectedOutput: map[string]string{
				"other1": "",
				"-flag1": "value1",
				"-flag2": "value2",
			},
		},
		{
			description: "flags and values inside quotes with spaces",
			input:       "other1 -flag1 \"value 1\" -flag2 \"value 2\"",
			expectedOutput: map[string]string{
				"other1": "",
				"-flag1": "value 1",
				"-flag2": "value 2",
			},
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := Parse(trial.input)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
