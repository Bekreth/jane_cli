//go:build windows

package terminal

import (
	"fmt"
	"os"
	"syscall"
)

type ScreenWriter struct {
	contextName string
}

var (
	kernel         *syscall.LazyDLL  = syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode *syscall.LazyProc = kernel.NewProc("SetConsoleMode")
)

const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x4

func NewScreenWriter(contextName string) ScreenWriter {
	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil {
		panic(err) //TODO: Handle this better
	}

	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	returnCode, _, err := setConsoleMode.Call(uintptr(syscall.Stdout), uintptr(mode))
	if returnCode == 0 {
		panic(err)
	}

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
