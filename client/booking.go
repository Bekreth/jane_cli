package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const bookingApi = "appointments/%v/book"

type BookingRequest struct {
	Book bool `json:"book"`
}
type Booking struct {
	StaffMemberID int       `json:"staff_member_id"`
	TreatmentID   int       `json:"treatment_id"`
	BookerType    string    `json:"booker_type"` //StaffMember
	PatientID     int       `json:"patient_id"`
	StartAt       time.Time `json:"start_at"` // 2023-04-04T17:00:00-07:00
	EndAt         time.Time `json:"end_at"`
	//TODO: Fix this v.v.v.v.v.v
	Duration int  `json:"duration"` //Duration in Seconds
	Break    bool `json:"break"`
	RoomID   int  `json:"room_id"`
}

func (client Client) buildBookingRequest(appointmentID int) string {
	return fmt.Sprintf("%v/%v/%v",
		client.getDomain(),
		apiBase3,
		fmt.Sprintf(bookingApi, appointmentID),
	)
}

func (client Client) BookAppointment(appointmentID int) error {
	client.logger.Debugf("booking an appointment")

	//TODO: Fix this v.v.v.v.v.v
	requestBody := BookingRequest{}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize booking request")
		return err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		client.buildBookingRequest(appointmentID),
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
