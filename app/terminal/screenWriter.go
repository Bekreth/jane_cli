package terminal

import "fmt"

type ScreenWriter struct {
	contextName string
}

func NewScreenWriter(contextName string) ScreenWriter {
	return ScreenWriter{
		contextName: contextName,
	}
}

func (writer ScreenWriter) WriteStringf(input string, args ...any) {
	writer.WriteString(fmt.Sprintf(input, args...))
}

func (writer ScreenWriter) WriteString(input string) {
	fmt.Print("\n\033[1A\033[K")
	fmt.Print("\r", writer.contextName, " ", input)
}

func (writer ScreenWriter) NewLine() {
	fmt.Println()
}
