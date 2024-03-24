package charts

import (
	"fmt"

	"github.com/Bekreth/jane_cli/domain/schedule"
)

type Chart struct {
	ChartEntries []ChartEntry `json:"chart_entries"`
	PerPage      int          `json:"per_page"`
	TotalEntries int          `json:"total_entries"`
}

type ChartEntry struct {
	ID            int                  `json:"id"`
	Title         string               `json:"title"`
	SignedState   string               `json:"signed_state"`
	PatientID     int                  `json:"patient_id"`
	AuthorID      int                  `json:"author_id"`
	AppointmentID int                  `json:"appointment_id"`
	Appointment   schedule.Appointment `json:"appointment"` //TODO: populate cache
	TreatmentID   int                  `json:"treatment_id"`
	CreatedAt     schedule.JaneTime    `json:"create_at"`
	EnteredOn     schedule.JaneTime    `json:"entered_on"`
	ChartParts    []ChartPart          `json:"chart_parts"`
	Snippet       string               `json:"snippet"`
}

func (entry ChartEntry) PrintSelector() string {
	if entry.Title != "" {
		return fmt.Sprintf("%v - %v", entry.EnteredOn.HumanDate(), entry.Title)
	} else if entry.Snippet != "" {
		end := 40
		if len(entry.Snippet) < end {
			end = len(entry.Snippet)
		}
		return fmt.Sprintf("%v - %v", entry.EnteredOn.HumanDate(), entry.Snippet[0:end])
	} else {
		return fmt.Sprintf("%v", entry.EnteredOn.HumanDate())
	}
}

func (entry ChartEntry) PrintText() string {
	output := ""
	for _, part := range entry.ChartParts {
		output += part.Text
	}
	return output
}

var DefaultChartEntry = ChartEntry{}

// Chart::Part::Note
type ChartPart struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type ChartPartUpdate struct {
	ID        int
	TextDelta TextDelta `json:"text_delta"`
	Text      string    `json:"text"`
	Label     string    `json:"note"`
}

type TextDelta struct {
	Ops []Ops `json:"ops"`
}

type Ops struct {
	Insert string `json:"insert"`
}
