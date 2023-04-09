package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const appointmentApi = "appointments"

type AppointmentRequest struct {
	Appointment Appointment `json:"appointment"`
	Book        bool        `json:"book"`
}

type Appointment struct {
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
	Break         bool      `json:"break"`
	LocationID    int       `json:"location_id"`
	RoomID        int       `json:"room_id"`
	StaffMemberID int       `json:"staff_member_id"`
}

type AppointmentResponse struct {
	ID int `json:"id"`
}

func (client Client) buildAppointmentRequest() string {
	return fmt.Sprintf("%v/%v/%v",
		client.getDomain(),
		apiBase2,
		appointmentApi,
	)
}

func (client Client) CreateAppointment(
	startDate time.Time,
	endDate time.Time,
	employeeBreak bool,
) (int, error) {
	client.logger.Debugf("creating appointment")

	requestBody := AppointmentRequest{
		Appointment: Appointment{
			StartAt:       startDate,
			EndAt:         endDate,
			Break:         employeeBreak,
			LocationID:    1,
			RoomID:        1,
			StaffMemberID: client.auth.UserID,
		},
		Book: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize booking request")
		return 0, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		client.buildAppointmentRequest(),
		strings.NewReader(string(jsonBody)),
	)
	if err != nil {
		client.logger.Infof("failed to serialize appointment request: %v", requestBody)
		return 0, err
	}

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to create appoint in Jane: %v", err)
		return 0, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("bad response from Jane: %v", err)
		return 0, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read message body")
		return 0, err
	}

	output := AppointmentResponse{}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		client.logger.Infof("failed to deserialize into patient struct")
		return 0, err
	}

	return output.ID, nil
}
