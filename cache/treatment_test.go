package cache

import (
	"testing"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/stretchr/testify/assert"
)

func TestMatchingTreatment(t *testing.T) {
	trials := []struct {
		description    string
		inputTreatment domain.Treatment
		inputName      string
		expectedOutput bool
	}{
		{
			description: "expecting to match",
			inputTreatment: domain.Treatment{
				Name: "60 minutes",
			},
			inputName:      "60",
			expectedOutput: true,
		},
		{
			description: "fail to match",
			inputTreatment: domain.Treatment{
				Name: "60 minutes",
			},
			inputName:      "30",
			expectedOutput: false,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := matchingTreatment(trial.inputTreatment, trial.inputName)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
