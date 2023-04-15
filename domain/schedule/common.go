package schedule

import (
	"strings"
	"time"
)

const humanDateFormat = "Jan 02"
const humanDateTimeFormat = "Jan 02 15:04"
const hourMinuteFormat = "15:04"
const dateTimeFormat = "2006-01-02T15:04:05"
const dateTimeFormatWithTimeStamp = "2006-01-02T15:04:05-07:00"

type JaneTime struct {
	time.Time
}

func NewJaneTime(input time.Time) JaneTime {
	return JaneTime{
		Time: time.Date(
			input.Year(),
			input.Month(),
			input.Day(),
			input.Hour(),
			input.Minute(),
			0,
			0,
			time.Local,
		),
	}
}

func (janeTime *JaneTime) UnmarshalJSON(bytes []byte) error {
	timeString := strings.Trim(string(bytes), "\"")
	if timeString == "null" {
		janeTime.Time = time.Time{}
		return nil
	}
	parsedTime, err := time.Parse(dateTimeFormat, timeString)
	if err != nil {
		parsedTime, err = time.Parse(dateTimeFormatWithTimeStamp, timeString)
	}
	janeTime.Time = NewJaneTime(parsedTime).Time
	return err
}

func (janeTime JaneTime) MarshalJSON() ([]byte, error) {
	timeString := janeTime.Format(dateTimeFormatWithTimeStamp)
	return []byte("\"" + timeString + "\""), nil
}

func (janeTime JaneTime) HumanDate() string {
	return janeTime.Format(humanDateFormat)
}
