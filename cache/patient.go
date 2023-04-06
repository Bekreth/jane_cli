package cache

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/domain"
)

func (cache Cache) FindPatients(patientName string) ([]domain.Patient, error) {
	possibleMatches := []domain.Patient{}
	for _, patient := range cache.patients {
		if matchingPatient(patient, patientName) {
			possibleMatches = append(possibleMatches, patient)
		}
	}
	if len(possibleMatches) > 0 {
		return possibleMatches, nil
	}
	cache.logger.Debugf("no patients found in cache, looking up via Jane client")

	fetchedPatients, err := cache.patientFetcher.FetchPatients(patientName)
	if err != nil {
		return possibleMatches, err
	}

	for _, patient := range fetchedPatients {
		cache.patients[patient.ID] = patient
	}
	cache.logger.Debugf("searching for patients again")
	for _, patient := range cache.patients {
		if matchingPatient(patient, patientName) {
			possibleMatches = append(possibleMatches, patient)
		}
	}

	return possibleMatches, err
}

func matchingPatient(patient domain.Patient, nameToCheck string) bool {
	firstName := strings.ToLower(patient.FirstName)
	lastName := strings.ToLower(patient.LastName)
	preferred := strings.ToLower(patient.PreferredFirstName)
	loweredName := strings.ToLower(nameToCheck)

	byFirst := strings.HasPrefix(firstName, loweredName)
	if byFirst {
		fmt.Printf("Matched by first %v: %v\n", loweredName, patient)
	}
	byLast := strings.HasPrefix(lastName, loweredName)
	if byFirst {
		fmt.Printf("Matched by last %v: %v\n", loweredName, patient)
	}
	byPref := strings.HasPrefix(preferred, loweredName)
	if byFirst {
		fmt.Printf("Matched by pref %v: %v\n", loweredName, patient)
	}

	output := byFirst || byLast || byPref
	return output
}
