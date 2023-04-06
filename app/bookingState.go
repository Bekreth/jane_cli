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

type patientFetcher interface {
	FindPatients(patientName string) ([]domain.Patient, error)
}

type appointmentBooker interface {
}

const bookingTimeFormat = "01.02T15:04"
const bookingDateFlag = "-d"
const appointmentFlag = "-a"
const patientFlag = "-p"

type bookingBuilder struct {
	building            bool
	patients            []domain.Patient
	targetPatient       domain.Patient
	appointmentDate     time.Time
	appointmentDuration time.Duration
}

type bookingState struct {
	logger         logger.Logger
	writer         screenWriter
	rootState      state
	currentBuffer  string
	nextState      state
	patientFetcher patientFetcher
	booking        bookingBuilder
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
	booking.booking = bookingBuilder{}
	booking.writer.newLine()
	booking.writer.writeString("")
}

func (booking *bookingState) handleKeyinput(character rune, key keyboard.Key) state {
	if booking.booking.building {
		return booking.handlePatientSelector(character)
	}
	keyHandler(key, &booking.currentBuffer, booking.triggerAutocomplete, booking.submit)

	if character != 0 {
		booking.currentBuffer += string(character)
	}

	booking.writer.writeString(booking.currentBuffer)
	return booking.nextState
}

func (booking *bookingState) handlePatientSelector(character rune) state {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		booking.writer.writeStringf(
			"selector value of %v unacceptable. aborting booking",
			string(character),
		)
		booking.writer.newLine()
		return booking.rootState
	}
	if index > len(booking.booking.patients) {
		booking.writer.writeStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index, len(booking.booking.patients),
		)
		booking.writer.newLine()
		return booking
	}
	booking.booking.targetPatient = booking.booking.patients[index-1]
	booking.logger.Debugf("selected patient at index %v", index)
	//TODO
	return booking.rootState
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

	booking.currentBuffer = ""
	patientName := flags[patientFlag]
	if patientName == "" {
		booking.writer.writeStringf("no name provided, use the %v flag", patientFlag)
		booking.writer.newLine()
		return
	}
	patients, err := booking.patientFetcher.FindPatients(patientName)
	if err != nil {
		booking.writer.writeStringf("failed to lookup patient %v : %v", patientName, err)
		booking.writer.newLine()
		return
	}
	if len(patients) == 0 {
		booking.writer.writeStringf("no patients found for %v", patientName)
		booking.writer.newLine()
		return
	}
	//TODO: Make this configurable
	if len(patients) > 8 {
		booking.writer.writeStringf("too many patients to render nicely for %v", patientName)
		booking.writer.newLine()
		return
	}

	booking.booking = bookingBuilder{
		building: true,
		patients: patients,
	}
	patientList := "Select intended patient\n"
	for i, patient := range patients {
		patientList += fmt.Sprintf("%v: %v %v \n", i+1, patient.FirstName, patient.LastName)
	}
	booking.writer.writeString(patientList)
}
