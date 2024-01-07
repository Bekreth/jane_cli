package util

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain"
)

type TreatmentFetcher interface {
	FindTreatment(treatmentName string) ([]domain.Treatment, error)
}

func ParseTreatmentFlag(
	fetcher TreatmentFetcher,
	treatmentName string,
) (domain.Treatment, []domain.Treatment, error) {
	treatment := domain.DefaultTreatment
	treatments := []domain.Treatment{}

	if treatmentName == "" {
		return treatment, treatments, fmt.Errorf("no treatment provided, use the %v flag", treatmentName)
	}
	treatments, err := fetcher.FindTreatment(treatmentName)
	if err != nil {
		return treatment, treatments, fmt.Errorf("failed to lookup treatments %v : %v", treatmentName, err)
	}

	if len(treatments) == 0 {
		return treatment, treatments, fmt.Errorf("no treatment found for %v", treatmentName)
	} else if len(treatments) == 1 {
		treatment = treatments[0]
	}
	return treatment, treatments, nil
}
