//go:build linux

package terminal

import (
	"fmt"
	"os"
)

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
	output := fmt.Sprintf("\u001B[2K\u000D%v %v", writer.contextName, input)
	os.Stdout.Write([]byte(output))
}

func (writer ScreenWriter) NewLine() {
	fmt.Println()
}
