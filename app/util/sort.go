package util

import "github.com/Bekreth/jane_cli/domain/schedule"

type AppointmentByDate []schedule.Appointment

func (byDate AppointmentByDate) Len() int {
	return len(byDate)
}
func (byDate AppointmentByDate) Less(i int, j int) bool {
	return byDate[i].StartAt.Time.Unix() < byDate[j].StartAt.Time.Unix()
}
func (byDate AppointmentByDate) Swap(i int, j int) {
	byDate[i], byDate[j] = byDate[j], byDate[i]
}

//TODO: Alphabetical patients
