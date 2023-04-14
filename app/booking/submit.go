package booking

import (
	"fmt"
	"strings"

	"github.com/Bekreth/jane_cli/app/terminal"
	"github.com/Bekreth/jane_cli/app/util"
)

const bookingTimeYearFormat = "06.01.02T15:04"
const bookingTimeFormat = "01.02T15:04"
const bookingDateFlag = "-d"
const treatmentFlag = "-t"
const patientFlag = "-p"

func (state *bookingState) Submit() {
	flags := terminal.ParseFlags(state.buffer.Read())
	state.buffer.Clear()
	if _, exists := flags[".."]; exists {
		state.nextState = state.rootState
		return
	} else if _, exists := flags["help"]; exists {
		state.printHelp()
		return
	}

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
		state.buffer.WriteStoreString(notifcation)
		return
	}
	builder := bookingBuilder{
		substate: unknown,
	}

	builder, err := state.parsePatientValue(flags[patientFlag], builder)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	builder, err = state.parseTreatmentValue(flags[treatmentFlag], builder)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	builder.appointmentDate, err = util.ParseDate(
		bookingTimeFormat,
		bookingTimeYearFormat,
		flags[bookingDateFlag],
	)
	if err != nil {
		state.buffer.WriteStoreString(err.Error())
		return
	}

	state.booking = builder
}

func (state *bookingState) printHelp() {
	// TODO: automate this list of elements
	state.buffer.WriteStoreString(fmt.Sprintf(
		"booking command is used to create new bookings:\n%v\n%v\n%v",
		"\t-d\tWhen to creating the appointment in the format of MM.DDTHH.MM",
		"\t-t\tThe treatment to use",
		"\t-p\tThe name of the patient (First, last, or preffered)",
	))
}

func (state *bookingState) parsePatientValue(
	patientName string,
	builder bookingBuilder,
) (bookingBuilder, error) {
	if patientName == "" {
		return builder, fmt.Errorf("no patient provided, use the %v flag", patientFlag)
	}
	patients, err := state.fetcher.FindPatients(patientName)
	if err != nil {
		return builder, fmt.Errorf("failed to lookup patient %v : %v", patientName, err)
	}
	builder.patients = patients
	if len(patients) == 0 {
		return builder, fmt.Errorf("no patients found for %v", patientName)
	} else if len(patients) == 1 {
		builder.targetPatient = patients[0]
	} else if len(patients) > 8 {
		return builder, fmt.Errorf("too many patients to render nicely for %v", patientName)
	}
	return builder, nil
}

func (state *bookingState) parseTreatmentValue(
	treatmentName string,
	builder bookingBuilder,
) (bookingBuilder, error) {
	if treatmentName == "" {
		return builder, fmt.Errorf("no treatment provided, use the %v flag", treatmentFlag)
	}
	treatments, err := state.fetcher.FindTreatment(treatmentName)
	if err != nil {
		return builder, fmt.Errorf("failed to lookup treatments %v : %v", treatmentName, err)
	}
	builder.treatments = treatments
	if len(treatments) == 0 {
		return builder, fmt.Errorf("no treatment found for %v", treatmentName)
	} else if len(treatments) == 1 {
		builder.targetTreatment = treatments[0]
	} else if len(treatments) > 8 {
		return builder, fmt.Errorf(
			"too many treatments to render nicely for %v",
			treatmentName,
		)
	}
	return builder, nil
}
