package util

import (
	"fmt"
	"time"

	"github.com/Bekreth/jane_cli/domain/schedule"
)

const DateFormat = "01.02"
const DateTimeFormat = "01.02T15:04"
const YearDateFormat = "2006.01.02"
const YearDateTimeFormat = "2006.01.02T15:04"

func ParseDate(
	dayFormat string,
	yearFormat string,
	dateString string,
) (schedule.JaneTime, error) {
	dateValue, err := time.Parse(yearFormat, dateString)
	if err == nil {
		return schedule.NewJaneTime(dateValue), err
	} else {
		dateValue, err = time.Parse(dayFormat, dateString)
		if err != nil {
			return schedule.NewJaneTime(time.Now()), fmt.Errorf(
				"unable to parse date %v, please write date in the format '%v' or '%v'",
				dateString,
				dayFormat,
				yearFormat,
			)
		}

		now := time.Now().Local()
		if now.Month() > dateValue.Month() {
			dateValue = dateValue.AddDate(now.Year()+1, 0, 0)
		} else {
			dateValue = dateValue.AddDate(now.Year(), 0, 0)
		}

		return schedule.NewJaneTime(dateValue), nil
	}
}
