package cache

import (
	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
	"github.com/Bekreth/jane_cli/logger"
)

type treatmentFetcher interface {
	FetchTreatments() ([]domain.Treatment, error)
}

type patientFetcher interface {
	FetchPatients(patientName string) ([]domain.Patient, error)
}

type scheduleFetcher interface {
	// TODO: autosave schedule
	FetchSchedule(
		startDate schedule.JaneTime,
		endDate schedule.JaneTime,
	) (schedule.Schedule, error)
}

type dataFetchers interface {
	treatmentFetcher
	patientFetcher
	scheduleFetcher
}

type Cache struct {
	logger           logger.Logger
	patients         map[int]domain.Patient
	treatments       map[int]domain.Treatment
	appointments     map[int]schedule.Appointment
	patientFetcher   patientFetcher
	treatmentFetcher treatmentFetcher
	scheduleFetcher  scheduleFetcher
}

func NewCache(
	logger logger.Logger,
	dataFetchers dataFetchers,
) (Cache, error) {
	return Cache{
		logger:           logger,
		patients:         make(map[int]domain.Patient),
		treatments:       make(map[int]domain.Treatment),
		appointments:     make(map[int]schedule.Appointment),
		patientFetcher:   dataFetchers,
		treatmentFetcher: dataFetchers,
		scheduleFetcher:  dataFetchers,
	}, nil
}
