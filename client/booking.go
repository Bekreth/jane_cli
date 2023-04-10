package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

const bookingApi = "appointments/%v/book"

type BookingRequest struct {
	Book bool `json:"book"`
}

type Booking struct {
	StaffMemberID int               `json:"staff_member_id"`
	TreatmentID   int               `json:"treatment_id"`
	BookerType    string            `json:"booker_type"`
	PatientID     int               `json:"patient_id"`
	StartAt       schedule.JaneTime `json:"start_at"`
	EndAt         schedule.JaneTime `json:"end_at"`
	Duration      int               `json:"duration"`
	Break         bool              `json:"break"`
	RoomID        int               `json:"room_id"`
}

func (client Client) buildBookingRequest(appointmentID int) string {
	return fmt.Sprintf("%v/%v/%v",
		client.getDomain(),
		apiBase3,
		fmt.Sprintf(bookingApi, appointmentID),
	)
}

func (client Client) BookAppointment(
	appointment Appointment,
	treatment domain.Treatment,
	patient domain.Patient,
) error {
	client.logger.Debugf("booking an appointment")

	requestBody := Booking{
		StaffMemberID: client.user.Auth.UserID,
		TreatmentID:   treatment.ID,
		BookerType:    "StaffMember",
		PatientID:     patient.ID,
		StartAt:       appointment.StartAt,
		EndAt:         appointment.EndAt,
		Duration:      int(treatment.ScheduledDuration.Duration.Seconds()),
		Break:         false,
		RoomID:        client.user.RoomID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize booking request")
		return err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		client.buildBookingRequest(appointment.ID),
		strings.NewReader(string(jsonBody)),
	)

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to get patient info from Jane")
		return err
	}
	//TODO v.v.v.v.v.v.v.v
	rBody, _ := ioutil.ReadAll(response.Body)
	client.logger.Debugf("RESPONSE: %v", string(rBody))
	//TODO ^.^.^.^.^.^.^.^

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return err
	}

	return nil
}
