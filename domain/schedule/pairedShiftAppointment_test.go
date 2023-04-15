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

	includeAll := map[AppointmentType]interface{}{
		Booked:      "",
		Arrived:     "",
		Unscheduled: "",
		Break:       "",
	}

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
			expectedOutput: `shift on Jan 01 from 09:00 to 17:00
09:00 to 12:00: booked with Roy
12:00 to 13:00: break
13:00 to 17:00: booked with Mary`,
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
			expectedOutput: `shift on Jan 01 from 09:00 to 17:00
09:00 to 12:00: booked with Roy
12:00 to 13:00: unscheduled
13:00 to 17:00: booked with Mary`,
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
			expectedOutput: `shift on Jan 01 from 09:00 to 17:00
09:00 to 12:00: booked with Roy
12:00 to 14:00: booked with Mary
14:00 to 17:00: unscheduled`,
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
			expectedOutput: `shift on Jan 01 from 09:00 to 17:00
09:00 to 10:00: unscheduled
10:00 to 12:00: booked with Roy
12:00 to 17:00: booked with Mary`,
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
			expectedOutput: `shift on Jan 01 from 09:00 to 17:00
09:00 to 10:00: booked with Mark
10:00 to 12:00: booked with Roy
12:00 to 17:00: booked with Mary`,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			trial.pairUndertest.include = includeAll
			trial.pairUndertest.showPassedAppointment = true
			actualOutput := trial.pairUndertest.ToString()
			assert.Equal(tt, trial.expectedOutput, actualOutput)
		})
	}
}
