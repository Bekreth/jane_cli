package app

import (
	"testing"

	"github.com/Bekreth/jane_cli/logger"
	"github.com/stretchr/testify/assert"
)

func TestScheduleState_TriggerAutocomplete(t *testing.T) {
	trials := []struct {
		description    string
		input          string
		expectedOutput string
	}{
		{
			description:    "Expect no autocomplete for uncompletable word",
			input:          "nope",
			expectedOutput: "nope",
		},
		{
			description:    "Autocomplete 'opening'",
			input:          "open",
			expectedOutput: "openings ",
		},
		{
			description:    "Autocomplete 'opening' when last word",
			input:          "pickle open",
			expectedOutput: "pickle openings ",
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			stateUnderTest := &scheduleState{
				logger:        logger.NewTestLogger(tt),
				currentBuffer: trial.input,
			}

			stateUnderTest.triggerAutocomplete()
			assert.Equal(tt, stateUnderTest.currentBuffer, trial.expectedOutput)

		})
	}
}
