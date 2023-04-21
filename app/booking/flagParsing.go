package booking

import "fmt"

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

func (state *bookingState) parseAppointmentValue(
	patientName string,
	builder bookingBuilder,
) (bookingBuilder, error) {
	appointments, err := state.fetcher.FindAppointments(
		builder.appointmentDate.ThisDay(),
		builder.appointmentDate.NextDay(),
		patientName,
	)
	if err != nil {
		return builder, fmt.Errorf("failed to lookup appointments : %v", err)
	}
	builder.appointments = appointments
	if len(builder.appointments) == 0 {
		outputMessage := "no appointments found on %v"
		if patientName != "" {
			outputMessage += " for " + patientName
		}
		return builder, fmt.Errorf(outputMessage, builder.appointmentDate.HumanDate())
	} else if len(builder.appointments) == 1 {
		builder.targetAppointment = builder.appointments[0]
	} else if len(builder.appointments) > 20 {
		return builder, fmt.Errorf(
			"too many appointments to render nicely on %v",
			builder.appointmentDate.HumanDate(),
		)
	}
	return builder, nil
}
