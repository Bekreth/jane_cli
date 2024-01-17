package terminal

import (
	"fmt"
	"os"

	tsize "github.com/kopoli/go-terminal-size"
)

type Window interface {
	ClearLine()
	MoveToLineStart()
	MoveToPreviousLine()
}

type Terminal struct {
	tsize.Size
	window Window
}

func NewTerminal(size tsize.Size, window Window) ScreenWriter {
	return Terminal{
		Size:   size,
		window: window,
	}
}

func (terminal Terminal) WriteStringf(input string, args ...any) {
	terminal.WriteString(fmt.Sprintf(input, args...))
}

func (terminal Terminal) WriteString(input string) {
	rows := (len(input) - 1) / terminal.Width
	terminal.window.MoveToLineStart()
	terminal.window.ClearLine()
	for i := rows; i > 0; i-- {
		terminal.window.MoveToPreviousLine()
		terminal.window.ClearLine()
	}
	os.Stdout.Write([]byte(input))
}

func (terminal Terminal) NewLine() {
	fmt.Println()
}
