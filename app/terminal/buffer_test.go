package terminal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testWriter struct {
	value string
}

func (writer *testWriter) WriteStringf(input string, args ...any) {
	writer.value = fmt.Sprintf(input, args...)
}

func (writer *testWriter) WriteString(input string) {
	writer.value = input
}

func (writer *testWriter) NewLine() {
	writer.value = "\n"
}

func TestAddCharacter(t *testing.T) {
	trials := []struct {
		description      string
		bufferValue      string
		cursorPosition   int
		expectedValue    string
		expectedPosition int
	}{
		{
			description:      "Cursor at end of string",
			bufferValue:      "sample tex",
			cursorPosition:   10,
			expectedValue:    "sample text",
			expectedPosition: 11,
		},
		{
			description:      "Cursor at start of string",
			bufferValue:      "sample tex",
			cursorPosition:   0,
			expectedValue:    "tsample tex",
			expectedPosition: 1,
		},
		{
			description:      "Cursor in the middle of a string",
			bufferValue:      "sample tex",
			cursorPosition:   5,
			expectedValue:    "samplte tex",
			expectedPosition: 6,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			thisWriter := testWriter{}
			bufferUnderTest := Buffer{
				contextName:    "test",
				cursorPosition: trial.cursorPosition,
				writer:         &thisWriter,
				currentValue:   trial.bufferValue,
				previousOutput: "",
			}

			bufferUnderTest.AddCharacter('t')

			actualOutput := thisWriter.value
			actualBuffer := bufferUnderTest.currentValue
			actualPosition := bufferUnderTest.cursorPosition

			assert.Equal(tt, "test: "+trial.expectedValue, actualOutput)
			assert.Equal(tt, trial.expectedValue, actualBuffer)
			assert.Equal(tt, trial.expectedPosition, actualPosition)
		})
	}
}
func TestRemoveCharacter(t *testing.T) {
	trials := []struct {
		description      string
		bufferValue      string
		cursorPosition   int
		expectedValue    string
		expectedPosition int
	}{
		{
			description:      "empty buffer",
			bufferValue:      "",
			cursorPosition:   0,
			expectedValue:    "",
			expectedPosition: 0,
		},
		{
			description:      "Cursor at end of buffer",
			bufferValue:      "sample text",
			cursorPosition:   11,
			expectedValue:    "sample tex",
			expectedPosition: 10,
		},
		{
			description:      "Cursor at start of buffer",
			bufferValue:      "sample text",
			cursorPosition:   0,
			expectedValue:    "sample text",
			expectedPosition: 0,
		},
		{
			description:      "Cursor in middle of buffer",
			bufferValue:      "sample text",
			cursorPosition:   6,
			expectedValue:    "sampl text",
			expectedPosition: 5,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			thisWriter := testWriter{}
			bufferUnderTest := Buffer{
				contextName:    "test",
				cursorPosition: trial.cursorPosition,
				writer:         &thisWriter,
				currentValue:   trial.bufferValue,
				previousOutput: "",
			}

			bufferUnderTest.RemoveCharacter()

			actualOutput := thisWriter.value
			actualBuffer := bufferUnderTest.currentValue
			actualPosition := bufferUnderTest.cursorPosition

			assert.Equal(tt, "test: "+trial.expectedValue, actualOutput)
			assert.Equal(tt, trial.expectedValue, actualBuffer)
			assert.Equal(tt, trial.expectedPosition, actualPosition)
		})
	}
}

func TestSkipRight(t *testing.T) {
	//                 *  *   *      *
	//             01234567890123456789012
	sampleText := "This is the sample text"

	trials := []struct {
		description      string
		startingPosition int
		expectedPosition int
	}{
		{
			description:      "Cursor at end of string",
			startingPosition: len(sampleText),
			expectedPosition: len(sampleText),
		},
		{
			description:      "Cursor at beginning of string",
			startingPosition: 0,
			expectedPosition: 4,
		},
		{
			description:      "Cursor in middle of string",
			startingPosition: 9,
			expectedPosition: 11,
		},
		{
			description:      "Cursor before end",
			startingPosition: 19,
			expectedPosition: 23,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			writer := testWriter{}
			bufferUnderTest := Buffer{
				contextName:    "test",
				cursorPosition: trial.startingPosition,
				writer:         &writer,
				currentValue:   sampleText,
				previousOutput: "",
			}

			bufferUnderTest.SkipRight()
			actualPosition := bufferUnderTest.cursorPosition

			assert.Equal(tt, trial.expectedPosition, actualPosition)
		})
	}
}

func TestSkipLeft(t *testing.T) {
	//                 *  *   *      *
	//             01234567890123456789012
	sampleText := "This is the sample text"

	trials := []struct {
		description      string
		startingPosition int
		expectedPosition int
	}{
		{
			description:      "Cursor at the end of string",
			startingPosition: len(sampleText),
			expectedPosition: 18,
		},
		{
			description:      "Cursor at beginning of string",
			startingPosition: 0,
			expectedPosition: 0,
		},
		{
			description:      "Cursor in middle of string",
			startingPosition: 9,
			expectedPosition: 7,
		},
		{
			description:      "Cursor after beginning",
			startingPosition: 2,
			expectedPosition: 0,
		},
	}

	for _, trial := range trials {
		t.Run(trial.description, func(tt *testing.T) {
			writer := testWriter{}
			bufferUnderTest := Buffer{
				contextName:    "test",
				cursorPosition: trial.startingPosition,
				writer:         &writer,
				currentValue:   sampleText,
				previousOutput: "",
			}

			bufferUnderTest.SkipLeft()
			actualPosition := bufferUnderTest.cursorPosition

			assert.Equal(tt, trial.expectedPosition, actualPosition)
		})
	}
}
