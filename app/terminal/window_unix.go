//go:build linux

package terminal

import (
	"fmt"
	"os"
)

const ESC = "\u001B"

type LinuxWindow struct {
}

func NewWindow() Window {
	return LinuxWindow{}
}

func (LinuxWindow) ClearLine() {
	os.Stdout.Write([]byte(fmt.Sprintf("%v%v", ESC, "[2K")))
}

func (LinuxWindow) MoveToLineStart() {
	os.Stdout.Write([]byte("\u000D"))
}

func (LinuxWindow) MoveToPreviousLine() {
	os.Stdout.Write([]byte(fmt.Sprintf("%v%v", ESC, "[A")))
}
