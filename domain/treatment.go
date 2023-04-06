package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SecondsDuration struct {
	time.Duration
}

func (duration *SecondsDuration) UnmarshalJSON(bytes []byte) error {
	durationString := strings.Trim(string(bytes), "\"")
	if durationString == "null" {
		duration.Duration = 0
		return nil
	}
	secondCount, err := strconv.Atoi(durationString)
	if err != nil {
		duration.Duration = 0
		return fmt.Errorf("provided string '%v' can't be parsed to int", durationString)
	}
	duration.Duration = time.Second * time.Duration(secondCount)
	return nil
}

type Treatment struct {
	ID                int             `json:"id"`
	BillingName       string          `json:"billing_name"`
	ScheduledDuration SecondsDuration `json:"scheduled_duration"`
}
