package charting

import "github.com/Bekreth/jane_cli/app/util"

func (state *chartingState) fetchPatients(flags map[string]string) {
	builder := chartingBuilder{
		substate: unknown,
		flow:     read,
	}

	patientName, exists := flags[patientFlag]
	if !exists {
		state.buffer.WriteStoreString("missing argument -p")
		return
	}

	var err error
	builder.targetPatient, builder.patients, err = util.ParsePatientValue(
		state.fetcher,
		patientName,
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	state.builder = builder
}
