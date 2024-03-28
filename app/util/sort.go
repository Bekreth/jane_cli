package util

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

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

type Treatments []domain.Treatment

func (treatments Treatments) Len() int {
	return len(treatments)
}
func (treatments Treatments) Less(i int, j int) bool {
	return treatments[i].ID < treatments[j].ID
}
func (treatments Treatments) Swap(i int, j int) {
	treatments[i], treatments[j] = treatments[j], treatments[i]
}

type Patients []domain.Patient

func (patients Patients) Len() int {
	return len(patients)
}
func (patients Patients) Less(i int, j int) bool {
	return patients[i].PrintName() < patients[j].PrintName()
}
func (patients Patients) Swap(i int, j int) {
	patients[i], patients[j] = patients[j], patients[i]
}
