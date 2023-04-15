package schedule

import (
	"testing"
	"time"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	startAt, err := time.Parse(hourMinuteFormat, "12:12")
	if err != nil {
		t.Fatal(err)
	}
	endAt, err := time.Parse(hourMinuteFormat, "15:15")
	if err != nil {
		t.Fatal(err)
	}

	trials := []struct {
		description     string
		testAppointment Appointment
		expectedOutput  string
	}{
		{
			description: "No appointment",
			testAppointment: Appointment{
				StartAt: JaneTime{startAt},
				EndAt:   JaneTime{endAt},
				State:   "unscheduled",
			},
			expectedOutput: "12:12 to 15:15: unscheduled",
		},
		{
			description: "Appointment with person",
			testAppointment: Appointment{
				StartAt: JaneTime{startAt},
				EndAt:   JaneTime{endAt},
				State:   "booked",
				Patient: domain.Patient{
					PreferredFirstName: "Billy",
				},
			},
			expectedOutput: "12:12 to 15:15: booked with Billy",
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := trial.testAppointment.ToString()
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
