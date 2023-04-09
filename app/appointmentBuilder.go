package app

import (
	"time"

	"github.com/Bekreth/jane_cli/client"
	"github.com/Bekreth/jane_cli/domain"
)

type appointmentBuilder struct {
	client client.Client
}

func (builder appointmentBuilder) BookPatient(
	patient domain.Patient,
	treatment domain.Treatment,
	startTime time.Time,
) error {
	return nil
}
