package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/util"
	"github.com/Bekreth/jane_cli/domain"
)

func (state chartingState) fetchPatients(
	flags map[string]string,
) (interactive.Interactive[domain.Patient], error) {
	patientName, exists := flags[patientFlag]
	if !exists {
		return nil, fmt.Errorf("missing argument -p")
	}

	var err error
	targetPatient, patients, err := util.ParsePatientValue(
		state.fetcher,
		patientName,
	)
	if err != nil {
		return nil, err
	}
	return interactive.NewPatientSelector(targetPatient, patients), nil
}
