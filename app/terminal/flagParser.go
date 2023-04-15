package terminal

import (
	"strings"
)

func ParseFlags(input string) map[string]string {
	output := map[string]string{}

	words := strings.Split(input, " ")
	previousKey := ""
	for _, word := range words {

		isFlag := strings.HasPrefix(word, "-")
		if isFlag {
			output[word] = ""
			previousKey = word
		} else if previousKey == "" {
			output[word] = ""
		} else {
			output[previousKey] = word
			previousKey = ""
		}
	}

	return output
}
