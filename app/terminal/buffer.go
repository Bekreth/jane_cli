package terminal

import (
	"fmt"
)

type Buffer struct {
	contextName    string
	cursorPosition int
	writer         ScreenWriter
	currentValue   string
	previousOutput string
}

func NewBuffer(writer ScreenWriter, contextName string) Buffer {
	return Buffer{
		contextName:    contextName,
		cursorPosition: 0,
		writer:         writer,
		currentValue:   "",
		previousOutput: "",
	}
}

func (buffer Buffer) contextString() string {
	return fmt.Sprintf("%v: %v", buffer.contextName, buffer.currentValue)
}

func (buffer *Buffer) MoveRight() {
	if buffer.cursorPosition < len(buffer.currentValue) {
		buffer.cursorPosition += 1
	}
}

func (buffer *Buffer) SkipRight() {
	charBreaks := indiciesOfChar(buffer.currentValue, ' ')
	for _, charBreak := range append(charBreaks, len(buffer.currentValue)) {
		if charBreak > buffer.cursorPosition {
			buffer.cursorPosition = charBreak
			return
		}
	}
}

func (buffer *Buffer) MoveLeft() {
	if buffer.cursorPosition >= 1 {
		buffer.cursorPosition -= 1
	}
}

func (buffer *Buffer) SkipLeft() {
	possibleBreak := 0
	charBreaks := indiciesOfChar(buffer.currentValue, ' ')
	for _, charBreak := range charBreaks {
		if charBreak > possibleBreak && charBreak < buffer.cursorPosition {
			fmt.Println("Possible", charBreak)
			possibleBreak = charBreak
		}
		if charBreak >= buffer.cursorPosition {
			break
		}
	}
	buffer.cursorPosition = possibleBreak
}

func (buffer *Buffer) AddCharacter(character rune) {
	if character != 0 {
		chared := string(character)
		bufferLength := len(buffer.currentValue)
		if bufferLength == buffer.cursorPosition {
			buffer.currentValue += chared
		} else if buffer.cursorPosition == 0 {
			buffer.currentValue = chared + buffer.currentValue
		} else {
			buffer.currentValue = buffer.currentValue[0:buffer.cursorPosition] +
				chared +
				buffer.currentValue[buffer.cursorPosition:]
		}
		buffer.MoveRight()
		buffer.writer.WriteString(buffer.contextString())
	}
}

func (buffer *Buffer) RemoveCharacter() {
	if len(buffer.currentValue) > 0 && buffer.cursorPosition != 0 {
		if len(buffer.currentValue) == 1 {
			buffer.currentValue = ""
		} else if len(buffer.currentValue) == buffer.cursorPosition {
			buffer.currentValue = buffer.currentValue[0 : len(buffer.currentValue)-1]
		} else {
			buffer.currentValue = buffer.currentValue[0:buffer.cursorPosition-1] +
				buffer.currentValue[buffer.cursorPosition:]
		}
		buffer.MoveLeft()
	}
	buffer.writer.WriteString(buffer.contextString())
}

func (buffer *Buffer) Read() string {
	return buffer.currentValue
}

func (buffer *Buffer) Clear() {
	buffer.currentValue = ""
}

func (buffer *Buffer) Write() {
	buffer.writer.WriteString(buffer.contextString())
}

func (buffer *Buffer) WriteNewLine() {
	buffer.writer.NewLine()
	buffer.Write()
}

func (buffer *Buffer) WriteString(input string) {
	buffer.currentValue = input
	buffer.Write()
}

func (buffer *Buffer) WriteStore() {
	buffer.WriteStoreString(buffer.contextString())
	buffer.currentValue = ""
}

func (buffer *Buffer) WriteStoreString(input string) {
	buffer.writer.NewLine()
	buffer.previousOutput = input
	buffer.writer.WriteString(input)
	buffer.writer.NewLine()
}

func (buffer *Buffer) WritePrevious() {
	buffer.writer.WriteString(buffer.previousOutput)
	buffer.Clear()
	buffer.WriteNewLine()
}
