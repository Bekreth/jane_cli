package app

func mapKeys(input map[string]state) []string {
	output := []string{}
	for key := range input {
		output = append(output, key)
	}
	return output
}
