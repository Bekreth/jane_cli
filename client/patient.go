package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Bekreth/jane_cli/domain"
)

const patientLookup = "patients/lookup"

//const patientLookup = "patient_lookup/lookup"

type PatientRequest struct {
	Autocomplete bool   `json:"autocomplete"`
	Limit        int    `json:"limit"`
	Name         string `json:"q"`
}

func (client Client) buildPatientRequest() string {
	return fmt.Sprintf("%v/%v/%v",
		client.getDomain(),
		apiBase2,
		patientLookup,
	)
}

func (client Client) FetchPatients(patientName string) ([]domain.Patient, error) {
	client.logger.Debugf("fetching patient")
	output := []domain.Patient{}

	requestBody := PatientRequest{
		Autocomplete: true,
		Limit:        50,
		Name:         patientName,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize patient request")
		return output, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		client.buildPatientRequest(),
		strings.NewReader(string(jsonBody)),
	)
	//	client.logger.Debugf("REQEUST: %v", request)
	if err != nil {
		client.logger.Infof("failed to build patient request")
		return output, err
	}
	request.Header = http.Header(commonHeaders)

	response, err := client.janeClient.Do(request)
	//	client.logger.Debugf("RESPONSE: %v", request)
	if err != nil {
		client.logger.Infof("failed to get patient info from Jane")
		return output, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return output, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read message body")
		return output, err
	}

	err = json.Unmarshal(bytes, &output)
	if err != nil {
		client.logger.Infof("failed to deserialize into patient struct")
		return output, err
	}

	return output, nil
}
