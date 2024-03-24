package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state chartingState) handleCreate(
	flags map[string]string,
) (chartingBuilder, error) {
	state.logger.Debugln("Handling chart create")
	output := state.builder
	output.flow = create
	output.substate = unknown

	// Setup Patient
	patientName, exists := flags[patientFlag]
	if !exists {
		return output, fmt.Errorf("missing argument -p")
	}

	var err error
	targetPatient, patients, err := util.ParsePatientValue(
		state.fetcher,
		patientName,
	)
	if err != nil {
		return output, err
	}
	output.patientSelector = interactive.NewPatientSelector(targetPatient, patients)

	// Setup Date of Appointment
	date, exists := flags[dateFlag]
	if exists {
		parsedDate, err := util.ParseDate(
			util.DateFormat,
			util.YearDateFormat,
			date,
		)
		if err != nil {
			return output, err
		}
		output.date = parsedDate
	} else {
		return output, fmt.Errorf("missing argument -d")
	}

	// Setup Note
	output.note, exists = flags[noteFlag]

	return output, nil
}
