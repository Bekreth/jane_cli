//go:build windows

package terminal

import (
	"fmt"
	"os"
)

const ESC = "\u001B"

type WindowsWindow struct {
}

func NewWindow() Window {
	return WindowsWindow{}
}

func (WindowsWindow) ClearLine() {
	os.Stdout.Write([]byte(fmt.Sprintf("%v%v", ESC, "[2M")))
}

func (WindowsWindow) MoveToLineStart() {
	os.Stdout.Write([]byte("\u000D"))
}

func (WindowsWindow) MoveToPreviousLine() {
	os.Stdout.Write([]byte(fmt.Sprintf("%v%v", ESC, "[A")))
}

/*
TODO: This section needs to be uncommented, but there's currently a problem with
JAWS that causes the first letter of a strings starting with /u001B[nK to be read
several times in a row.  This can be replicated in powershell or command prompt with
the command `echo ^[[2Khello` which reads as "h h hello".  Until this issue is resolved,
the text manipulation for the sighted is deprioritized for those requiring screen
readers

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
*/
//TODO: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
/*
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
	output := fmt.Sprintf("\u000D%v %v%v", writer.contextName, input, " ")
	os.Stdout.Write([]byte(output))
}

func (writer ScreenWriter) NewLine() {
	fmt.Println()
}
*/
