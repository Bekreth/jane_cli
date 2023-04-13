package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Bekreth/jane_cli/domain/schedule"
)

const appointmentApi = "appointments"

type AppointmentRequest struct {
	Appointment Appointment `json:"appointment"`
	Book        bool        `json:"book"`
}

type Appointment struct {
	ID            int               `json:"id,omitempty"`
	StartAt       schedule.JaneTime `json:"start_at"`
	EndAt         schedule.JaneTime `json:"end_at"`
	Break         bool              `json:"break"`
	LocationID    int               `json:"location_id"`
	StaffMemberID int               `json:"staff_member_id"`
}

func (client Client) buildAppointmentRequest() string {
	return fmt.Sprintf("%v/%v/%v",
		client.getDomain(),
		apiBase2,
		appointmentApi,
	)
}

func (client Client) CreateAppointment(
	startDate schedule.JaneTime,
	endDate schedule.JaneTime,
	employeeBreak bool,
) (Appointment, error) {
	client.logger.Debugf("creating appointment")
	output := Appointment{}

	requestBody := AppointmentRequest{
		Appointment: Appointment{
			StartAt:       startDate,
			EndAt:         endDate,
			Break:         employeeBreak,
			LocationID:    client.user.LocationID,
			StaffMemberID: client.user.Auth.UserID,
		},
		Book: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize booking request")
		return output, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		client.buildAppointmentRequest(),
		strings.NewReader(string(jsonBody)),
	)
	if err != nil {
		client.logger.Infof("failed to serialize appointment request: %v", requestBody)
		return output, err
	}
	request.Header = commonHeaders

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to create appoint in Jane: %v", err)
		return output, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("bad response from Jane: %v", err)
		return output, err
	}
	bytes, err := io.ReadAll(response.Body)

	if err != nil {
		client.logger.Infof("failed to read message body")
		return output, err
	}

	err = json.Unmarshal(bytes, &output)
	if err != nil {
		client.logger.Infof("failed to deserialize into appointment struct: %v", err)
		return output, err
	}

	client.logger.Infof("created appointment %v at %v", output.ID, output.StartAt)

	return output, nil
}
