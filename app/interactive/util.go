package interactive

import (
	"fmt"
	"strconv"
	"strings"
)

const ESCAPE_STRING = "(or ESC to back out)"

func PrintSelectorList[R interface{}](interactive Interactive[R]) string {
	output := []string{}
	pageNumber, totalPages, elementCount := interactive.PagingInfo()
	if elementCount == 0 {
		return ""
	}
	pageInfo := ""
	if totalPages > 1 {
		pageInfo = fmt.Sprintf(
			"Showing page %v of %v, 'f' to page forwards, 'b' to page backwards ",
			pageNumber+1,
			totalPages,
		)
	}

	header := fmt.Sprintf(
		"%v. %v%v",
		interactive.PossibleSelections()[0].PrintHeader(),
		pageInfo,
		ESCAPE_STRING,
	)
	for i, possible := range interactive.PossibleSelections() {
		if i == 0 {
			output = append(output, header)
		}
		output = append(output, fmt.Sprintf("%v: %v", i+1, possible.PrintSelector()))
	}
	return strings.Join(output, "\n")
}

func ElementSelector[R interface{}](
	character rune,
	input []Selection[R],
) (Selection[R], error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("Input has size 0")
	}
	index, err := strconv.Atoi(string(character))
	if err != nil || index == 0 {
		return nil, fmt.Errorf(
			"selector value of '%v' unacceptable. select a value between 1 and %v",
			string(character),
			len(input),
		)
	}
	if index > len(input) {
		return nil, fmt.Errorf(
			"selector value of '%v' is too large.  select a value between 1 and %v",
			index,
			len(input),
		)
	}
	return input[index-1], nil
}

func mod(a, b int) int {
	return (a%b + b) % b
}
