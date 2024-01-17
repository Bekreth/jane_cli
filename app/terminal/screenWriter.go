package terminal

type ScreenWriter interface {
	WriteStringf(input string, args ...any)
	WriteString(input string)
	NewLine()
}
