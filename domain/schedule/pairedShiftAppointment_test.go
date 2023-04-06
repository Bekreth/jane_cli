package schedule

import (
	"testing"
	"time"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/stretchr/testify/assert"
)

func makeTime(t *testing.T, timeAsString string) JaneTime {
	output, err := time.Parse(hourMinuteFormat, timeAsString)
	if err != nil {
		t.Fatal(err)
	}
	return JaneTime{output}
}

func TestPairedShiftAppointment(t *testing.T) {

	trials := []struct {
		description    string
		pairUndertest  pairedShiftAppointment
		expectedOutput string
	}{
		{
			description: "Shift with no holes and ordered events",
			pairUndertest: pairedShiftAppointment{
				shift: Shift{
					StartAt: makeTime(t, "09:00"),
					EndAt:   makeTime(t, "17:00"),
				},
				appointment: []Appointment{
					{
						StartAt: makeTime(t, "09:00"),
						EndAt:   makeTime(t, "12:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Roy",
						},
					},
					{
						StartAt: makeTime(t, "12:00"),
						EndAt:   makeTime(t, "13:00"),
						State:   "break",
					},
					{
						StartAt: makeTime(t, "13:00"),
						EndAt:   makeTime(t, "17:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Mary",
						},
					},
				},
			},
			expectedOutput: `shift from 09:00 to 17:00
* booked from 09:00 to 12:00 with Roy
* break from 12:00 to 13:00
* booked from 13:00 to 17:00 with Mary`,
		},
		{
			description: "Shift with a hole in the middle of the schedule",
			pairUndertest: pairedShiftAppointment{
				shift: Shift{
					StartAt: makeTime(t, "09:00"),
					EndAt:   makeTime(t, "17:00"),
				},
				appointment: []Appointment{
					{
						StartAt: makeTime(t, "09:00"),
						EndAt:   makeTime(t, "12:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Roy",
						},
					},
					{
						StartAt: makeTime(t, "13:00"),
						EndAt:   makeTime(t, "17:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Mary",
						},
					},
				},
			},
			expectedOutput: `shift from 09:00 to 17:00
* booked from 09:00 to 12:00 with Roy
* unscheduled from 12:00 to 13:00
* booked from 13:00 to 17:00 with Mary`,
		},
		{
			description: "Shift with a hole at the end of the schedule",
			pairUndertest: pairedShiftAppointment{
				shift: Shift{
					StartAt: makeTime(t, "09:00"),
					EndAt:   makeTime(t, "17:00"),
				},
				appointment: []Appointment{
					{
						StartAt: makeTime(t, "09:00"),
						EndAt:   makeTime(t, "12:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Roy",
						},
					},
					{
						StartAt: makeTime(t, "12:00"),
						EndAt:   makeTime(t, "14:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Mary",
						},
					},
				},
			},
			expectedOutput: `shift from 09:00 to 17:00
* booked from 09:00 to 12:00 with Roy
* booked from 12:00 to 14:00 with Mary
* unscheduled from 14:00 to 17:00`,
		},
		{
			description: "Shift with a hole at the beggining of the schedule",
			pairUndertest: pairedShiftAppointment{
				shift: Shift{
					StartAt: makeTime(t, "09:00"),
					EndAt:   makeTime(t, "17:00"),
				},
				appointment: []Appointment{
					{
						StartAt: makeTime(t, "10:00"),
						EndAt:   makeTime(t, "12:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Roy",
						},
					},
					{
						StartAt: makeTime(t, "12:00"),
						EndAt:   makeTime(t, "17:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Mary",
						},
					},
				},
			},
			expectedOutput: `shift from 09:00 to 17:00
* unscheduled from 09:00 to 10:00
* booked from 10:00 to 12:00 with Roy
* booked from 12:00 to 17:00 with Mary`,
		},
		{
			description: "Schedule out of schedule",
			pairUndertest: pairedShiftAppointment{
				shift: Shift{
					StartAt: makeTime(t, "09:00"),
					EndAt:   makeTime(t, "17:00"),
				},
				appointment: []Appointment{
					{
						StartAt: makeTime(t, "10:00"),
						EndAt:   makeTime(t, "12:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Roy",
						},
					},
					{
						StartAt: makeTime(t, "09:00"),
						EndAt:   makeTime(t, "10:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Mark",
						},
					},
					{
						StartAt: makeTime(t, "12:00"),
						EndAt:   makeTime(t, "17:00"),
						State:   "booked",
						Patient: domain.Patient{
							PreferredFirstName: "Mary",
						},
					},
				},
			},
			expectedOutput: `shift from 09:00 to 17:00
* booked from 09:00 to 10:00 with Mark
* booked from 10:00 to 12:00 with Roy
* booked from 12:00 to 17:00 with Mary`,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			actualOutput := trial.pairUndertest.ToString()
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
