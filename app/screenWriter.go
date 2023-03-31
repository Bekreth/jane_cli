package app

import "fmt"

type screenWriter struct {
	contextName string
}

func (writer screenWriter) writeString(input string) {
	fmt.Print("\n\033[1A\033[K")
	fmt.Print("\r", writer.contextName, " ", input)
}

func (writer screenWriter) newLine() {
	fmt.Println()
}
