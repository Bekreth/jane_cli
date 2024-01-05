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
const signChartApi = "patients/%v/chart_entries/%v"

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

func (client Client) buildChartAppointmentSetter(chartID int) string {
	return fmt.Sprintf(
		"%v/%v/%v/%v/set_appointment",
		client.getDomain(),
		apiBase2,
		chartLookupApi,
		chartID,
	)
}

func (client Client) buildChartSignRequest(patientID int, chartID int) string {
	return fmt.Sprintf(
		"%v/%v/%v",
		client.getDomain(),
		apiBase2,
		fmt.Sprintf(
			signChartApi,
			patientID,
			chartID,
		),
	)
}

func (client Client) FetchPatientCharts(
	patientID int,
) ([]charts.ChartEntry, error) {
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
		return []charts.ChartEntry{}, err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return []charts.ChartEntry{}, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read response body: %v", err)
	}

	patientCharts := charts.Chart{}
	err = json.Unmarshal(bytes, &patientCharts)
	if err != nil {
		client.logger.Infof("failed to read patient charts: %v", err)
	}

	client.logger.Infof("Got patient chart")
	return patientCharts.ChartEntries, nil
}

/**
* CreatePatientCharts uses a GET request against Jane to create a new, empty chart.
 */
func (client Client) CreatePatientCharts(
	patientID int,
	appointmentID int,
) (charts.ChartEntry, error) {
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
		return charts.ChartEntry{}, err
	}
	client.logger.Debugf("RAW: %v", response)

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("Bad response from Jane: %v", err)
		return charts.ChartEntry{}, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read response body: %v", err)
	}

	newChart := charts.ChartEntry{}
	err = json.Unmarshal(bytes, &newChart)
	if err != nil {
		client.logger.Infof("failed to read patient charts: %v", err)
	}
	client.logger.Debugf("Processed: %v", newChart)

	client.logger.Infof("Got new patient chart")
	return newChart, err
}

/**
 * UpdatePatientChart sends a PUT request for the chart **part** on its ID
 */
func (client Client) UpdatePatientChart(
	chartPartID int,
	chartText string,
) error {
	client.logger.Infof("updating patient chart")
	requestBody := charts.ChartPartUpdate{
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

func (client Client) SetChartingAppointment(chartID int, appointmentID int) error {
	client.logger.Infof("setting appointment ID on chart %v", chartID)
	requestBody := struct {
		AppointmentID int `json:"appointment_id"`
	}{
		AppointmentID: appointmentID,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		client.logger.Infof("failed to serialize chart set appointment ID")
		return err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		client.buildChartAppointmentSetter(chartID),
		strings.NewReader(string(jsonBody)),
	)
	if err != nil {
		client.logger.Infoln("failed to serialize chart set appointment request")
		return err
	}
	request.Header = commonHeaders

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to set chart appointment in Jane: %v", err)
		return err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("bad response from Jane %v: %v", response.StatusCode, err)
		return err
	}
	client.logger.Infof("Updated chart %v with appointment %v", chartID, appointmentID)

	return nil
}

func (client Client) SignChart(chart charts.ChartEntry, patientID int) error {
	client.logger.Infof("signing chart %v", chart.ID)

	jsonBody, err := json.Marshal(chart)
	if err != nil {
		client.logger.Infof("failed to serialize chart for signing")
		return err
	}

	request, err := http.NewRequest(
		http.MethodPut,
		client.buildChartSignRequest(patientID, chart.ID),
		strings.NewReader(string(jsonBody)),
	)
	if err != nil {
		client.logger.Infoln("failed to serialize chart signing request")
		return err
	}
	request.Header = commonHeaders

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("failed to sign chart: %v", err)
		return err
	}

	if err = checkStatusCode(response); err != nil {
		client.logger.Infof("bad response from Jane %v: %v", response.StatusCode, err)
		return err
	}
	client.logger.Infof("Successfully signed chart %v", chart.ID)

	return nil
}
