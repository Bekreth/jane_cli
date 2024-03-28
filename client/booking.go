package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bekreth/jane_cli/domain"
	"github.com/Bekreth/jane_cli/domain/schedule"
)

const bookingApi = "appointments/%v/book"

type BookingRequest struct {
	Book        bool    `json:"book"`
	Appointment Booking `json:"appointment"`
}

type Booking struct {
	StaffMemberID int               `json:"staff_member_id"`
	TreatmentID   int               `json:"treatment_id"`
	BookerType    string            `json:"booker_type"`
	PatientID     int               `json:"patient_id"`
	StartAt       schedule.JaneTime `json:"start_at"`
	EndAt         schedule.JaneTime `json:"end_at"`
	State         string            `json:"state"`
	Duration      int               `json:"duration"`
	Break         bool              `json:"break"`
	WithinShift   bool              `json:"within_shiftk"`
	LocationID    int               `json:"location_id"`
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

	requestBody := BookingRequest{
		Book: true,
		Appointment: Booking{
			StaffMemberID: client.user.Auth.UserID,
			TreatmentID:   treatment.ID,
			BookerType:    "StaffMember",
			PatientID:     patient.ID,
			StartAt:       appointment.StartAt,
			EndAt:         appointment.EndAt,
			State:         "reserved",
			Duration:      int(treatment.ScheduledDuration.Duration.Seconds()),
			Break:         false,
			WithinShift:   true,
			LocationID:    client.user.LocationID,
		},
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
	if err != nil {
		client.logger.Infof("failed to build booking request")
		return err
	}
	request.Header = http.Header(commonHeaders)

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to get patient info from Jane")
		return err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return err
	}

	return nil
}
