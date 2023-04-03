package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Bekreth/jane_cli/domain/schedule"
)

const apiBase = "admin/api/v2"
const calendar = "calendar"
const timeFormat = "2006-01-02"

func (client Client) buildScheduleRequest(startDate time.Time, endDate time.Time) string {
	return fmt.Sprintf(
		"%v/%v/%v?start_date=%v&end_date=%v&staff_member_ids[]=%v",
		client.getDomain(),
		apiBase,
		calendar,
		startDate.Format(timeFormat),
		endDate.Format(timeFormat),
		client.auth.UserID,
	)
}

func (client Client) FetchSchedule(
	startDate time.Time,
	endDate time.Time,
) (schedule.Schedule, error) {
	client.logger.Debugf("fetching scheudle")

	request, err := http.NewRequest(
		http.MethodGet,
		client.buildScheduleRequest(startDate, endDate),
		nil,
	)

	if err != nil {
		client.logger.Infof("failed to build fetch schedule request: %v", err)
	}

	response, err := client.janeClient.Do(request)
	if err != nil {
		client.logger.Infof("got a bad fetch schedule response: %v", err)
		return schedule.Schedule{}, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.logger.Infof("failed to read response body: %v", err)
	}

	fetchedSchedule := schedule.Schedule{}
	err = json.Unmarshal(bytes, &fetchedSchedule)
	if err != nil {
		client.logger.Infof("failed to read schedule: %v", err)
	}

	client.logger.Debugf(
		"Got schedule for %v to %v",
		startDate.Format(timeFormat),
		endDate.Format(timeFormat),
	)
	return fetchedSchedule, nil
}
