package charts

import (
	"github.com/Bekreth/jane_cli/domain/schedule"
)

type Chart struct {
	ChartEntries []ChartEntry `json:"chart_entries"`
	PerPage      int          `json:"per_page"`
	TotalEntries int          `json:"total_entries"`
}

type ChartEntry struct {
	ID            int                  `json:"id"`
	SignedState   string               `json:"signed_state"`
	PatientID     int                  `json:"patient_id"`
	AuthorID      int                  `json:"author_id"`
	AppointmentID int                  `json:"appointment_id"`
	Appointment   schedule.Appointment `json:"appointment"` //TODO: populate cache
	TreatmentID   int                  `json:"treatment_id"`
	CreatedAt     schedule.JaneTime    `json:"create_at"`
	EnteredOn     schedule.JaneTime    `json:"entered_on"`
	ChartParts    []ChartPart          `json:"chart_part"`
	Snippet       string               `json:"snippet"`
}

var DefaultChartEntry = ChartEntry{}

// Chart::Part::Upload
// Chart::Part::ChiefComplaint
// Chart::Part::CheckBoxes
// Chart::Part::Note
type ChartPart struct {
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
