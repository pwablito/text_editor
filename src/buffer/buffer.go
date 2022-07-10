package buffer

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

var Logger *log.Logger

type EditableBuffer struct {
	Content        string
	CursorPosition int
}

func ReadFileToBuffer(filename string) *EditableBuffer {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		Logger.Fatalf("Couldn't read file: %v", err)
	}

	text := string(content)
	return &EditableBuffer{
		Content:        text,
		CursorPosition: 0,
	}
}

func WriteBufferToFile(filename string, buffer *EditableBuffer) error {
	// TODO use existing file permissions if file exists
	var permission fs.FileMode = 0644
	info, err := os.Stat(filename)
	if err == nil {
		permission = info.Mode()
	}
	if !IsWritable(filename) {
		return errors.New(fmt.Sprintf("Can't write to %s", filename))
	}
	err = os.WriteFile(filename, []byte(buffer.Content), permission)
	if err != nil {
		Logger.Fatalf("Couldn't write file: %v", err)
	}
	return nil
}

func (buffer *EditableBuffer) Insert(char rune) {
	buffer.Content = buffer.contentBefore(buffer.CursorPosition) + string(char) + buffer.contentAfter(buffer.CursorPosition)
	buffer.CursorPosition++
}

func (buffer EditableBuffer) GetTermUiCompatibleOutput() string {
	cursor_rendered := buffer.contentBefore(buffer.CursorPosition) + "|" + buffer.contentAfter(buffer.CursorPosition)
	output := ""
	for _, char := range cursor_rendered {
		switch char {
		case '\t':
			output += "    "
			break
		default:
			output += string(char)
		}
	}
	return output
}

func IsWritable(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}

func (buffer EditableBuffer) contentBefore(position int) string {
	if position == 0 {
		return ""
	}
	return buffer.Content[:position]
}

func (buffer EditableBuffer) contentAfter(position int) string {
	if position == len(buffer.Content) {
		return ""
	}
	return buffer.Content[position:]
}

func (buffer *EditableBuffer) ArrowRight() (rebuildUI bool) {
	rebuildUI = false
	if buffer.CursorPosition != len(buffer.Content) {
		buffer.CursorPosition++
		rebuildUI = true
	}
	return
}

func (buffer *EditableBuffer) ArrowLeft() (rebuildUI bool) {
	rebuildUI = false
	if buffer.CursorPosition != 0 {
		buffer.CursorPosition--
		rebuildUI = true
	}
	return
}

func (buffer *EditableBuffer) Backspace() (rebuildUI bool) {
	rebuildUI = false
	if buffer.CursorPosition != 0 {
		buffer.Content = buffer.contentBefore(buffer.CursorPosition-1) + buffer.contentAfter(buffer.CursorPosition)
		buffer.CursorPosition--
		rebuildUI = true
	}
	return
}

func (buffer *EditableBuffer) Delete() (rebuildUI bool) {
	rebuildUI = false
	if buffer.CursorPosition != len(buffer.Content) {
		buffer.Content = buffer.contentBefore(buffer.CursorPosition) + buffer.contentAfter(buffer.CursorPosition+1)
		rebuildUI = true
	}
	return
}

func (buffer EditableBuffer) positionInLine() int {
	offset := 0
	for buffer.Content[buffer.CursorPosition-offset] != '\n' && offset != buffer.CursorPosition {
		offset++
	}
	return offset
}

func (buffer *EditableBuffer) MoveToStartOfNextLine() (rebuildUI bool) {
	rebuildUI = false
	for buffer.CursorPosition != len(buffer.Content) && buffer.Content[buffer.CursorPosition] != '\n' {
		buffer.CursorPosition++
		rebuildUI = true
	}
	if buffer.CursorPosition < len(buffer.Content) {
		buffer.CursorPosition++
		rebuildUI = true
	}
	return
}

func (buffer *EditableBuffer) seekPosInLine(pos int) (rebuildUI bool) {
	rebuildUI = false
	numForward := 0
	for numForward <= pos && len(buffer.Content) > numForward+buffer.CursorPosition && buffer.Content[buffer.CursorPosition+numForward+1] != '\n' {
		numForward++
		rebuildUI = true
	}
	buffer.CursorPosition += numForward
	return
}

func SplitTextToLines(text string) []string {
	lines := make([]string, 0)
	current_line := ""
	for _, character := range text {
		if character == '\n' {
			lines = append(lines, current_line)
			current_line = ""
		} else {
			current_line += string(character)
		}
	}
	lines = append(lines, current_line)
	return lines
}

func (buffer EditableBuffer) CursorLineNumber() int {
	return len(SplitTextToLines(buffer.contentBefore(buffer.CursorPosition))) - 1
}
func (buffer EditableBuffer) CursorPositionInLine() int {
	lines := SplitTextToLines(buffer.contentBefore(buffer.CursorPosition))
	return len(lines[len(lines)-1])
}

func (buffer *EditableBuffer) ArrowDown() bool {
	oldPosition := buffer.CursorPosition
	pos := buffer.positionInLine()
	buffer.MoveToStartOfNextLine()
	buffer.seekPosInLine(pos)
	return buffer.CursorPosition != oldPosition
}

func (buffer *EditableBuffer) ArrowUp() bool {
	oldCursorPosition := buffer.CursorPosition
	cursorLineNumber := buffer.CursorLineNumber()
	if cursorLineNumber == 0 {
		buffer.CursorPosition = 0
	} else {
		Logger.Print("Only know how to do up arrow when already at the top line")
	}
	return buffer.CursorPosition != oldCursorPosition
}
