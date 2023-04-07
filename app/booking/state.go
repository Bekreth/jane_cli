package booking

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/app/terminal"
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

// TODO use this
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
	logger    logger.Logger
	writer    terminal.ScreenWriter
	fetcher   bookingDataFetcher
	rootState terminal.State

	currentBuffer string
	nextState     terminal.State
	booking       bookingBuilder
}

func NewState(
	logger logger.Logger,
	writer terminal.ScreenWriter,
	fetcher bookingDataFetcher,
	rootState terminal.State,
) terminal.State {
	return &bookingState{
		logger:    logger,
		writer:    writer,
		fetcher:   fetcher,
		rootState: rootState,
	}
}

func (bookingState) Name() string {
	return "booking"
}

func (booking *bookingState) Initialize() {
	booking.logger.Debugf(
		"entering booking. available states %v",
		booking.rootState.Name(),
	)
	booking.nextState = booking
	booking.booking = bookingBuilder{
		substate: argument,
	}
	booking.writer.NewLine()
	booking.writer.WriteString("")
}

func (booking *bookingState) HandleKeyinput(
	character rune,
	key keyboard.Key,
) terminal.State {
	switch booking.booking.substate {
	case bookingConfirmation:
		booking.confirmBooking(character)
	case treatmentSelector:
		booking.booking.targetTreatment = booking.handleTreatmentSelector(character)
	case patientSelector:
		booking.booking.targetPatient = booking.handlePatientSelector(character)
	default:
		terminal.KeyHandler(
			key,
			&booking.currentBuffer,
			booking.triggerAutocomplete,
			booking.Submit,
		)
		if character != 0 {
			booking.currentBuffer += string(character)
		}
		booking.writer.WriteString(booking.currentBuffer)
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
		booking.writer.WriteString(confirmationString)
	case treatmentSelector:
		treatmentList := "Select intended treatment\n"
		for i, treatment := range booking.booking.treatments {
			treatmentList += fmt.Sprintf("%v: %v \n", i+1, treatment.Name)
		}
		booking.writer.WriteString(treatmentList)
	case patientSelector:
		patientList := "Select intended patient\n"
		for i, patient := range booking.booking.patients {
			patientList += fmt.Sprintf("%v: %v %v \n", i+1, patient.FirstName, patient.LastName)
		}
		booking.writer.WriteString(patientList)
	default:
	}

	return booking.nextState
}

func (booking *bookingState) handleTreatmentSelector(character rune) domain.Treatment {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		booking.writer.WriteStringf(
			"selector value of %v unacceptable. select a value between 1 and %v",
			string(character),
			len(booking.booking.treatments),
		)
		booking.writer.NewLine()
	}
	if index > len(booking.booking.treatments) {
		booking.writer.WriteStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index,
			len(booking.booking.treatments),
		)
		booking.writer.NewLine()
	}
	booking.logger.Debugf("selected treatment at index %v", index)
	return booking.booking.treatments[index-1]
}

func (booking *bookingState) handlePatientSelector(character rune) domain.Patient {
	index, err := strconv.Atoi(string(character))
	if err != nil {
		booking.writer.WriteStringf(
			"selector value of %v unacceptable.  select a value between 1 and %v",
			string(character),
			len(booking.booking.patients),
		)
		booking.writer.NewLine()
	}
	if index > len(booking.booking.patients) {
		booking.writer.WriteStringf(
			"selector value of %v is too large.  select a value between 1 and %v",
			index, len(booking.booking.patients),
		)
		booking.writer.NewLine()
	}
	booking.logger.Debugf("selected patient at index %v", index)
	return booking.booking.patients[index-1]
}

func (booking *bookingState) confirmBooking(character rune) {
	switch string(character) {
	case "y":
		fallthrough
	case "Y":
		booking.writer.WriteString("submitting booking")
		booking.writer.NewLine()
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
		booking.writer.WriteString("aborting booking")
		booking.writer.NewLine()
		booking.nextState = booking.rootState
	default:
		booking.writer.WriteStringf(
			"input of %v not support.  Confirm booking (Y/n)?",
			string(character),
		)
		booking.writer.NewLine()
	}
}

func (booking *bookingState) triggerAutocomplete() {
}

func (booking *bookingState) Submit() {
	if booking.currentBuffer == ".." {
		booking.nextState = booking.rootState
		return
	}
	flags := terminal.ParseFlags(booking.currentBuffer)
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
		joined := strings.Join(terminal.MapKeysString(missingFlags), ", ")
		notifcation := fmt.Sprintf("missing arguments %v", joined)
		booking.writer.WriteString(notifcation)
		booking.writer.NewLine()
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
		booking.writer.WriteStringf("no name provided, use the %v flag", patientFlag)
		booking.writer.NewLine()
		return
	}
	patients, err := booking.fetcher.FindPatients(patientName)
	if err != nil {
		booking.writer.WriteStringf("failed to lookup patient %v : %v", patientName, err)
		booking.writer.NewLine()
		booking.nextState = booking.rootState
		return
	}
	builder.patients = patients
	if len(patients) == 0 {
		booking.writer.WriteStringf("no patients found for %v", patientName)
		booking.writer.NewLine()
		return
	} else if len(patients) == 1 {
		builder.targetPatient = patients[0]
	} else if len(patients) > 8 {
		booking.writer.WriteStringf("too many patients to render nicely for %v", patientName)
		booking.writer.NewLine()
		return
	}

	treatmentName := flags[appointmentFlag]
	if patientName == "" {
		booking.writer.WriteStringf("no name provided, use the %v flag", patientFlag)
		booking.writer.NewLine()
		return
	}
	treatments, err := booking.fetcher.FindTreatment(treatmentName)
	if err != nil {
		booking.writer.WriteStringf("failed to lookup treatments %v : %v", treatmentName, err)
		booking.writer.NewLine()
		booking.nextState = booking.rootState
		return
	}
	builder.treatments = treatments
	if len(treatments) == 0 {
		booking.writer.WriteStringf("no treatment found for %v", treatmentName)
		booking.writer.NewLine()
	} else if len(treatments) == 1 {
		builder.targetTreatment = treatments[0]
	} else if len(treatments) > 8 {
		booking.writer.WriteStringf(
			"too many treatments to render nicely for %v",
			treatmentName,
		)
		booking.writer.NewLine()
	}

	booking.booking = builder
}
