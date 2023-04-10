package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Bekreth/jane_cli/domain"
)

const treatmentApi = "treatments/boot_treatments"

func (client Client) buildTreatmentRequest() string {
	return fmt.Sprintf("%v/%v/%v",
		client.getDomain(),
		apiBase2,
		treatmentApi,
	)
}

func (client Client) FetchTreatments() ([]domain.Treatment, error) {
	client.logger.Debugf("fetching treatments")
	treatmentList := []domain.Treatment{}

	request, err := http.NewRequest(
		http.MethodGet,
		client.buildTreatmentRequest(),
		nil,
	)
	if err != nil {
		client.logger.Infof("failed to build treatment request")
		return treatmentList, err
	}

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to get treatment info from Jane")
		return treatmentList, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return treatmentList, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read message body")
		return treatmentList, err
	}

	err = json.Unmarshal(bytes, &treatmentList)
	if err != nil {
		client.logger.Infof("failed to deserialize into patient struct")
		return treatmentList, err
	}

	output := []domain.Treatment{}
	for _, treatment := range treatmentList {
		client.logger.Debugf("%v %v", treatment.StaffMemberID, treatment.Name)
		if treatment.StaffMemberID == client.user.Auth.UserID {
			output = append(output, treatment)
		}
	}

	return output, nil
}
