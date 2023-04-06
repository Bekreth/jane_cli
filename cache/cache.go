package cache

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
)

type treatmentFetcher interface {
	FetchTreatments() ([]domain.Treatment, error)
}

type patientFetcher interface {
	FetchPatients(patientName string) ([]domain.Patient, error)
}

type Cache struct {
	logger           logger.Logger
	patients         map[int]domain.Patient
	treatments       map[int]domain.Treatment
	patientFetcher   patientFetcher
	treatmentFetcher treatmentFetcher
}

func NewCache(
	logger logger.Logger,
	patientFetcher patientFetcher,
	treatmentFetcher treatmentFetcher,
) (Cache, error) {
	return Cache{
		logger:           logger,
		patients:         make(map[int]domain.Patient),
		treatments:       make(map[int]domain.Treatment),
		patientFetcher:   patientFetcher,
		treatmentFetcher: treatmentFetcher,
	}, nil
}
