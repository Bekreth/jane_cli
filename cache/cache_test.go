package cache

import (
	"testing"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/stretchr/testify/assert"
)

func TestMatchingPatient(t *testing.T) {
	trials := []struct {
		description string
		patient     domain.Patient
		nameToCheck string
		shouldMatch bool
	}{
		{
			description: "Should match by first name",
			patient: domain.Patient{
				FirstName: "Bobby",
			},
			nameToCheck: "bob",
			shouldMatch: true,
		},
		{
			description: "Should match by last name",
			patient: domain.Patient{
				LastName: "Bobby",
			},
			nameToCheck: "bob",
			shouldMatch: true,
		},
		{
			description: "Should match by preferred name",
			patient: domain.Patient{
				PreferredFirstName: "Bobby",
			},
			nameToCheck: "bob",
			shouldMatch: true,
		},
		{
			description: "Should fail to match",
			patient: domain.Patient{
				FirstName:          "Roy",
				LastName:           "Boston",
				PreferredFirstName: "Job",
			},
			nameToCheck: "bob",
			shouldMatch: false,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			didMatch := matchingPatient(trial.patient, trial.nameToCheck)
			assert.Equal(tt, trial.shouldMatch, didMatch)
		})
	}
}
