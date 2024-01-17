package terminal

import "fmt"

type Buffer struct {
	contextName    string
	writer         ScreenWriter
	currentValue   string
	previousOutput string
}

func NewBuffer(writer ScreenWriter, contextName string) Buffer {
	return Buffer{
		contextName:    contextName,
		writer:         writer,
		currentValue:   "",
		previousOutput: "",
	}
}

func (buffer Buffer) contextString() string {
	return fmt.Sprintf("%v: %v", buffer.contextName, buffer.currentValue)
}

func (buffer *Buffer) AddCharacter(character rune) {
	if character != 0 {
		buffer.currentValue += string(character)
		buffer.writer.WriteString(buffer.contextString())
	}
}

func (buffer *Buffer) RemoveCharacter() {
	if len(buffer.currentValue) > 0 {
		buffer.currentValue = buffer.currentValue[0 : len(buffer.currentValue)-1]
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
