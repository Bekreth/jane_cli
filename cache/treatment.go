package cache

import (
	"strings"

	"github.com/Bekreth/jane_cli/domain"
)

func (cache Cache) FindTreatment(treatmentName string) ([]domain.Treatment, error) {
	possibleMatches := []domain.Treatment{}
	for _, treatment := range cache.treatments {
		if matchingTreatment(treatment, treatmentName) {
			possibleMatches = append(possibleMatches, treatment)
		}
	}

	if len(possibleMatches) > 0 {
		return possibleMatches, nil
	}
	cache.logger.Debugf("no treatments found in cache, looking up via Jane client")

	fetchedTreatments, err := cache.treatmentFetcher.FetchTreatments()
	if err != nil {
		return possibleMatches, err
	}

	for _, treatment := range fetchedTreatments {
		cache.treatments[treatment.ID] = treatment
	}
	cache.logger.Debugf("searching for patients again")
	for _, treatment := range cache.treatments {
		if matchingTreatment(treatment, treatmentName) {
			possibleMatches = append(possibleMatches, treatment)
		}
	}

	return possibleMatches, nil
}

func matchingTreatment(treatment domain.Treatment, nameToCheck string) bool {
	return strings.Contains(treatment.BillingName, nameToCheck)
}
