package util

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
)

type PatientFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
}

func ParsePatientValue(
	fetcher PatientFetcher,
	patientName string,
) (domain.Patient, []domain.Patient, error) {
	patient := domain.DefaultPatient
	patients := []domain.Patient{}
	if patientName == "" {
		//TODO: Globalize standard flags
		return patient, patients, fmt.Errorf("no patient provided, use the %v flag", "-p")
	}
	patients, err := fetcher.FindPatients(patientName)
	if err != nil {
		return patient, patients, fmt.Errorf("failed to lookup patient %v : %v", patientName, err)
	}
	if len(patients) == 0 {
		return patient, patients, fmt.Errorf("no patients found for %v", patientName)
	} else if len(patients) == 1 {
		patient = patients[0]
	}
	return patient, patients, nil
}
