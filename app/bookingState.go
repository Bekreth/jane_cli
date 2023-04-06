package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/logger"
	"github.com/eiannone/keyboard"
)

type bookingDataFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
	FindTreatment(treatmentName string) ([]domain.Treatment, error)
	// BookPatient(
	//
	//	patient domain.Patient,
	//	treatment domain.Treatment,
	//	startTime time.Time,
	//
	// ) error
}

const bookingTimeFormat = "01.02T15:04"
const bookingDateFlag = "-d"
const appointmentFlag = "-a"
const patientFlag = "-p"

type substate = string

const (
	argument            substate = "arguemnt"
	unknown                      = "unknown"
	bookingConfirmation          = "bookingConfirmation"
	patientSelector              = "patientSelector"
	treatmentSelector            = "treatmentSelector"
)

type bookingBuilder struct {
	substate        substate
	patients        []domain.Patient
	targetPatient   domain.Patient
	treatments      []domain.Treatment
	targetTreatment domain.Treatment
	appointmentDate time.Time
}

type bookingState struct {
	logger        logger.Logger
	writer        screenWriter
	rootState     state
	currentBuffer string
	nextState     state
	fetcher       bookingDataFetcher
	booking       bookingBuilder
}

func (bookingState) name() string {
	return "booking"
}

func (booking *bookingState) initialize() {
	booking.logger.Debugf(
		"entering booking. available states %v",
		booking.rootState.name(),
	)
	booking.nextState = booking
	booking.booking = bookingBuilder{
		substate: argument,
	}
	booking.writer.newLine()
	booking.writer.writeString("")
}

func (booking *bookingState) handleKeyinput(character rune, key keyboard.Key) state {
	switch booking.booking.substate {
	case bookingConfirmation:
		booking.confirmBooking(character)
	case treatmentSelector:
		booking.handleTreatmentSelector(character)
	case patientSelector:
		booking.handlePatientSelector(character)
	default:
		keyHandler(key, &booking.currentBuffer, booking.triggerAutocomplete, booking.submit)
		if character != 0 {
			booking.currentBuffer += string(character)
		}
		booking.writer.writeString(booking.currentBuffer)
	}

	if booking.booking.substate != argument {
		if booking.booking.targetPatient == domain.DefaultPatient {
			booking.booking.substate = patientSelector
		} else if booking.booking.targetTreatment == domain.DefaultTreatment {
			booking.booking.substate = treatmentSelector
		} else {
			booking.booking.substate = bookingConfirmation
		}
	}

	switch booking.booking.substate {
	case bookingConfirmation:
		confirmationString := fmt.Sprintf(
			"Book %v for a %v at %v? (Y/n)",
			booking.booking.targetPatient.PreferredFirstName,
			booking.booking.targetTreatment.Name,
			booking.booking.appointmentDate.Format(bookingTimeFormat),
		)
		booking.writer.writeString(confirmationString)
	case treatmentSelector:
		treatmentList := "Select intended treatment\n"
		for i, treatment := range booking.booking.treatments {
			treatmentList += fmt.Sprintf("%v: %v \n", i+1, treatment.Name)
		}
		booking.writer.writeString(treatmentList)
	case patientSelector:
		patientList := "Select intended patient\n"
		for i, patient := range booking.booking.patients {
			patientList += fmt.Sprintf("%v: %v %v \n", i+1, patient.FirstName, patient.LastName)
		}
		booking.writer.writeString(patientList)
	default:
	}

	return booking.nextState
}

func (booking *bookingState) handleTreatmentSelector(character rune) domain.Treatment {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		booking.writer.writeStringf(
			"selector value of %v unacceptable. select a value between 1 and %v",
			string(character),
			len(booking.booking.treatments),
		)
		booking.writer.newLine()
	}
	if index > len(booking.booking.treatments) {
		booking.writer.writeStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index,
			len(booking.booking.treatments),
		)
		booking.writer.newLine()
	}
	booking.logger.Debugf("selected treatment at index %v", index)
	return booking.booking.treatments[index-1]
}

func (booking *bookingState) handlePatientSelector(character rune) domain.Patient {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		booking.writer.writeStringf(
			"selector value of %v unacceptable.  select a value between 1 and %v",
			string(character),
			len(booking.booking.patients),
		)
		booking.writer.newLine()
	}
	if index > len(booking.booking.patients) {
		booking.writer.writeStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index, len(booking.booking.patients),
		)
		booking.writer.newLine()
	}
	booking.logger.Debugf("selected patient at index %v", index)
	return booking.booking.patients[index-1]
}

func (booking *bookingState) confirmBooking(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		booking.writer.writeString("submitting booking")
		booking.writer.newLine()
		//TODO
		//err := booking.fetcher.BookPatient(
		//	booking.booking.targetPatient,
		//	booking.booking.targetTreatment,
		//	booking.booking.appointmentDate,
		//)
		//if err != nil {
		//	booking.writer.writeStringf("failed to make appointment: %v", err)
		//	booking.writer.newLine()
		//}
	case "n":
		fallthrough
	case "N":
		booking.writer.writeString("aborting booking")
		booking.writer.newLine()
		booking.nextState = booking.rootState
	default:
		booking.writer.writeStringf(
			"input of %v not support.  Confirm booking (Y/n)?",
			string(character),
		)
		booking.writer.newLine()
	}
}

func (booking *bookingState) shutdown() {
	booking.currentBuffer = ""
}

func (booking *bookingState) triggerAutocomplete() {
}

func (booking *bookingState) submit() {
	flags := parseFlags(booking.currentBuffer)
	booking.logger.Debugf("submitting query flags: %v", flags)
	missingFlags := map[string]string{
		"-d": "",
		"-a": "",
		"-p": "",
	}

	for key := range missingFlags {
		delete(missingFlags, key)
	}
	if len(missingFlags) != 0 {
		joined := strings.Join(mapKeysString(missingFlags), ", ")
		notifcation := fmt.Sprintf("missing arguments %v", joined)
		booking.writer.writeString(notifcation)
		booking.writer.newLine()
		return
	}

	// TODO: Date.... deal with year turn over

	booking.currentBuffer = ""

	builder := bookingBuilder{
		substate: unknown,
	}
	booking.logger.Debugf("builder 1: %v", builder)
	patientName := flags[patientFlag]
	if patientName == "" {
		booking.writer.writeStringf("no name provided, use the %v flag", patientFlag)
		booking.writer.newLine()
		return
	}
	patients, err := booking.fetcher.FindPatients(patientName)
	if err != nil {
		booking.writer.writeStringf("failed to lookup patient %v : %v", patientName, err)
		booking.writer.newLine()
		return
	}
	builder.patients = patients
	if len(patients) == 0 {
		booking.writer.writeStringf("no patients found for %v", patientName)
		booking.writer.newLine()
		return
	} else if len(patients) == 1 {
		builder.targetPatient = patients[0]
	} else if len(patients) > 8 {
		booking.writer.writeStringf("too many patients to render nicely for %v", patientName)
		booking.writer.newLine()
		return
	}

	booking.logger.Debugf("builder 2: %v", builder)
	treatmentName := flags[appointmentFlag]
	if patientName == "" {
		booking.writer.writeStringf("no name provided, use the %v flag", patientFlag)
		booking.writer.newLine()
		return
	}
	treatments, err := booking.fetcher.FindTreatment(treatmentName)
	if err != nil {
		booking.writer.writeStringf("failed to lookup treatments %v : %v", treatmentName, err)
		booking.writer.newLine()
		return
	}
	builder.treatments = treatments
	if len(treatments) == 0 {
		booking.writer.writeStringf("no treatment found for %v", treatmentName)
		booking.writer.newLine()
	} else if len(treatments) == 1 {
		builder.targetTreatment = treatments[0]
	} else if len(treatments) > 8 {
		booking.writer.writeStringf(
			"too many treatments to render nicely for %v",
			treatmentName,
		)
		booking.writer.newLine()
	}

	booking.logger.Debugf("builder 3: %v", builder)
	booking.booking = builder
}
