package charting

import (
	"strings"
)

func (state *chartingState) printHelp() {
	helpString := []string{
		"Charting command is to view, create, and edit charts",
		"Available subcommands",
		"\tread\tRead the chart for a given patient.  Requires patient name.  If no date is provided, a list of prior dates will be presented",
		"\tcreate\tMake a new chart.  Defaults to creating the chart for the most recent appointment.  Provide patient name and date if you want to create a chart for a previous appointment",
		//TODO: Check this
		"Available flags",
		"\t-d\tDate.  Format is MM.DDTHH.MM",
		"\t-p\tThe name of the patient (First, last, or preffered)",
	}
	state.buffer.WriteStoreString(strings.Join(helpString, "\n"))
}
