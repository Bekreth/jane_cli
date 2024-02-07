package util

import "github.com/Bekreth/jane_cli/app/states"

func MapKeysString(input map[string]string) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}

func MapKeys(input map[string]states.State) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}
