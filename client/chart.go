package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Bekreth/jane_cli/domain/charts"
)

const chartLookupApi = "chart_entries"
const updateChartApi = "chart_parts"
const newChartApi = "chart_entries/new"

func (client Client) buildChartLookupRequest(patientID int) string {
	pageID := 1
	return fmt.Sprintf(
		"%v/%v/%v?treatment_plan_id=null&patient_id=%v&page=%v", 
		client.getDomain(), 
		apiBase2,
		chartLookupApi,
		patientID,
		pageID,
	)
}

func (client Client) buildNewChartRequest(
	patientID int,
	appointmentID int,
) string {
	return fmt.Sprintf(
		"%v/%v/%v?patient_id=%v&appointment_id=%v&template_identifier=system_template_note", 
		client.getDomain(), 
		apiBase2,
		newChartApi,
		patientID,
		appointmentID,
	)
}

func (client Client) buildChartUpdateRequest(chartPartID int) string {
	return fmt.Sprintf(
		"%v/%v/%v/%v", 
		client.getDomain(), 
		apiBase2,
		updateChartApi,
		chartPartID,
	)
}

func (client Client) FetchPatientCharts(
	patientID int,
) ([]charts.Chart, error){
	client.logger.Debugf("fetching patient charts")
	request, err := http.NewRequest(
		http.MethodGet,
		client.buildChartLookupRequest(patientID),
		nil,
	)
	if err != nil {
		client.logger.Infof("failed to build fetch patients charts request: %v", err)
	}
	request.Header = commonHeaders

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("got a bad fetch patient charts response: %v", err)
		return []charts.Chart{}, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return []charts.Chart{}, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read response body: %v", err)
	}

	patientCharts := []charts.Chart{}
	err = json.Unmarshal(bytes, &patientCharts)
	if err != nil {
		client.logger.Infof("failed to read patient charts: %v", err)
	}

	client.logger.Infof("Got patient chart")
	return patientCharts, nil
}

func (client Client) CreatePatientCharts(
	patientID int,
	appointmentID int,
) (charts.Chart, error) {
	client.logger.Debugf("fetching new patient chart")
	request, err := http.NewRequest(
		http.MethodGet,
		client.buildNewChartRequest(patientID, appointmentID),
		nil,
	)
	request.Header = http.Header(commonHeaders)

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("got a bad new chart response: %v", err)
		return charts.Chart{}, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return charts.Chart{}, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read response body: %v", err)
	}

	newChart := charts.Chart{}
	err = json.Unmarshal(bytes, &newChart)
	if err != nil {
		client.logger.Infof("failed to read patient charts: %v", err)
	}

	client.logger.Infof("Got new patient chart")
	return newChart, err
}

func (client Client) UpdatePatientChart(
	chartPartID int,
	chartText string,
) error {
	client.logger.Infof("updating patient chart")
	requestBody := charts.ChartPart{
		TextDelta: charts.TextDelta{
			Ops: []charts.Ops{
				{
					Insert: chartText,
				},
			},
		},
		Text:  chartText,
		Label: "Note",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize chart update request")
		return err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		client.buildChartUpdateRequest(chartPartID),
		strings.NewReader(string(jsonBody)),
	)
	if err != nil {
		client.logger.Infoln("failed to serialize chart update request")
		return err
	}
	request.Header = commonHeaders
	
	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to update chart in Jane: %v", err)
		return err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("bad response from Jane %v: %v", response.StatusCode, err)
		return err
	}
	client.logger.Infof("Updated chart part %v", chartPartID)

	return nil
}
