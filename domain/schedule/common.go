package schedule

import (
	"strings"
	"time"
)

const hourMinuteFormat = "15:04"
const dateTimeFormat = "2006-01-02T15:04:05"

type JaneTime struct {
	time.Time
}

func (janeTime *JaneTime) UnmarshalJSON(bytes []byte) error {
	timeString := strings.Trim(string(bytes), "\"")
	if timeString == "null" {
		janeTime.Time = time.Time{}
		return nil
	}
	parsedTime, err := time.Parse(dateTimeFormat, timeString)
	janeTime.Time = parsedTime
	return err
}
