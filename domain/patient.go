package domain

import "fmt"

type Patient struct {
	ID                 int    `json:"id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	PreferredFirstName string `json:"preferred_first_name"`
}

func (patient Patient) PrintName() string {
	return fmt.Sprintf("%v %v", patient.PreferredFirstName, patient.LastName)
}

var DefaultPatient = Patient{}
