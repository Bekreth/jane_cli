package terminal

type Buffer struct {
	writer         ScreenWriter
	currentValue   string
	previousOutput string
}

func NewBuffer(writer ScreenWriter) Buffer {
	return Buffer{
		writer:         writer,
		currentValue:   "",
		previousOutput: "",
	}
}

func (buffer *Buffer) AddCharacter(character rune) {
	if character != 0 {
		buffer.currentValue += string(character)
		buffer.writer.WriteString(buffer.currentValue)
	}
}

func (buffer *Buffer) RemoveCharacter() {
	if len(buffer.currentValue) > 0 {
		buffer.currentValue = buffer.currentValue[0 : len(buffer.currentValue)-1]
	}
	buffer.writer.WriteString(buffer.currentValue)
}

func (buffer *Buffer) Read() string {
	return buffer.currentValue
}

func (buffer *Buffer) PrintHeader() {
	buffer.writer.NewLine()
	buffer.writer.WriteString(buffer.currentValue)
}

func (buffer *Buffer) Clear() {
	buffer.currentValue = ""
}

func (buffer *Buffer) Write() {
	buffer.writer.WriteString(buffer.currentValue)
}

func (buffer *Buffer) WriteString(input string) {
	buffer.currentValue = input
	buffer.Write()
}

func (buffer *Buffer) WriteStore() {
	buffer.WriteStoreString(buffer.currentValue)
	buffer.currentValue = ""
}

func (buffer *Buffer) WriteStoreString(input string) {
	buffer.previousOutput = input
	buffer.writer.WriteString(input)
	buffer.writer.NewLine()
}

func (buffer *Buffer) WritePrevious() {
	buffer.writer.WriteString(buffer.previousOutput)
	buffer.writer.NewLine()
	buffer.writer.WriteString("")
}
