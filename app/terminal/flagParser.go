package terminal

import (
	"strings"
)

func ParseFlags(input string) map[string]string {
	output := map[string]string{}

	words := splitWords(input)
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

func splitWords(input string) []string {
	var words []string
	insideQuotes := false
	currentWord := ""
	for _, char := range input + " " {
		if char == '"' {
			insideQuotes = !insideQuotes
			continue
		}
		if char == ' ' && !insideQuotes {
			words = append(words, currentWord)
			currentWord = ""
			continue
		}
		currentWord += string(char)
	}
	return words
}
