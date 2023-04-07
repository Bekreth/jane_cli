package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
)

// TODO use this
const bookingTimeFormat = "01.02T15:04"
const bookingDateFlag = "-d"
const appointmentFlag = "-a"
const patientFlag = "-p"

func (state *bookingState) Submit() {
	if state.currentBuffer == ".." {
		state.nextState = state.rootState
		return
	}
	flags := terminal.ParseFlags(state.currentBuffer)
	state.logger.Debugf("submitting query flags: %v", flags)
	missingFlags := map[string]string{
		"-d": "",
		"-a": "",
		"-p": "",
	}

	for key := range missingFlags {
		delete(missingFlags, key)
	}
	if len(missingFlags) != 0 {
		joined := strings.Join(terminal.MapKeysString(missingFlags), ", ")
		notifcation := fmt.Sprintf("missing arguments %v", joined)
		state.writer.WriteString(notifcation)
		state.writer.NewLine()
		return
	}

	// TODO: Date.... deal with year turn over

	state.currentBuffer = ""

	builder := bookingBuilder{
		substate: unknown,
	}
	state.logger.Debugf("builder 1: %v", builder)
	patientName := flags[patientFlag]
	if patientName == "" {
		state.writer.WriteStringf("no name provided, use the %v flag", patientFlag)
		state.writer.NewLine()
		return
	}
	patients, err := state.fetcher.FindPatients(patientName)
	if err != nil {
		state.writer.WriteStringf("failed to lookup patient %v : %v", patientName, err)
		state.writer.NewLine()
		state.nextState = state.rootState
		return
	}
	builder.patients = patients
	if len(patients) == 0 {
		state.writer.WriteStringf("no patients found for %v", patientName)
		state.writer.NewLine()
		return
	} else if len(patients) == 1 {
		builder.targetPatient = patients[0]
	} else if len(patients) > 8 {
		state.writer.WriteStringf("too many patients to render nicely for %v", patientName)
		state.writer.NewLine()
		return
	}

	treatmentName := flags[appointmentFlag]
	if patientName == "" {
		state.writer.WriteStringf("no name provided, use the %v flag", patientFlag)
		state.writer.NewLine()
		return
	}
	treatments, err := state.fetcher.FindTreatment(treatmentName)
	if err != nil {
		state.writer.WriteStringf("failed to lookup treatments %v : %v", treatmentName, err)
		state.writer.NewLine()
		state.nextState = state.rootState
		return
	}
	builder.treatments = treatments
	if len(treatments) == 0 {
		state.writer.WriteStringf("no treatment found for %v", treatmentName)
		state.writer.NewLine()
	} else if len(treatments) == 1 {
		builder.targetTreatment = treatments[0]
	} else if len(treatments) > 8 {
		state.writer.WriteStringf(
			"too many treatments to render nicely for %v",
			treatmentName,
		)
		state.writer.NewLine()
	}

	state.booking = builder
}
