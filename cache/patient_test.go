package cache

import (
	"testing"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/stretchr/testify/assert"
)

func TestMatchingPatient(t *testing.T) {
	test_patient := domain.Patient{
		FirstName:          "bob",
		LastName:           "ross",
		PreferredFirstName: "painter",
	}
	trials := []struct {
		description    string
		inputName      string
		expectedOutput bool
	}{
		{
			description:    "Should match by first name",
			inputName:      "Bob",
			expectedOutput: true,
		},
		{
			description:    "Should match by last name",
			inputName:      "Ross",
			expectedOutput: true,
		},
		{
			description:    "Should match by preferred name",
			inputName:      "Painter",
			expectedOutput: true,
		},
		{
			description:    "Should match on multiple names",
			inputName:      "Bob Ross",
			expectedOutput: true,
		},
		{
			description:    "Should fail to match if one name of multiples is bad",
			inputName:      "Bob Marley",
			expectedOutput: false,
		},
		{
			description:    "Should fail to match if one name of multiples is bad",
			inputName:      "Marley",
			expectedOutput: false,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := matchingPatient(test_patient, trial.inputName)
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
