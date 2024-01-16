package charting

import (
	"fmt"

	"github.com/Bekreth/jane_cli/app/interactive"
	"github.com/Bekreth/jane_cli/app/util"
)

func (state chartingState) handleRead(
	flags map[string]string,
) (chartingBuilder, error) {
	state.logger.Debugln("Handling chart read")
	output := state.builder
	output.flow = read
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

	return output, nil
}
