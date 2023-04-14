package terminal

import (
	"strings"
)

func ParseFlags(input string) map[string]string {
	output := map[string]string{}

	words := strings.Split(input, " ")
	previousKey := ""
	for _, word := range words {
		if previousKey != "" {
			output[previousKey] = word
			previousKey = ""
		} else {
			output[word] = ""
		}
		if strings.HasPrefix(word, "-") {
			previousKey = word
		}
	}

	return output
}
