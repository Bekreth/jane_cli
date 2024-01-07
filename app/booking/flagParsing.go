package booking

import "fmt"

func (state *bookingState) parseTreatmentValue(
	treatmentName string,
	builder bookingBuilder,
) (bookingBuilder, error) {
	if treatmentName == "" {
		return builder, fmt.Errorf("no treatment provided, use the %v flag", treatmentFlag)
	}
	treatments, err := state.fetcher.FindTreatment(treatmentName)
	if err != nil {
		return builder, fmt.Errorf("failed to lookup treatments %v : %v", treatmentName, err)
	}
	builder.treatments = treatments
	if len(treatments) == 0 {
		return builder, fmt.Errorf("no treatment found for %v", treatmentName)
	} else if len(treatments) == 1 {
		builder.targetTreatment = treatments[0]
	} else if len(treatments) > 8 {
		return builder, fmt.Errorf(
			"too many treatments to render nicely for %v",
			treatmentName,
		)
	}
	return builder, nil
}
