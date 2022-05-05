package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type EditableBuffer struct {
	Content        string
	CursorPosition int
}

func ReadFileToBuffer(filename string) *EditableBuffer {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Couldn't read file: %v", err)
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
		log.Fatalf("Couldn't write file: %v", err)
	}
	return nil
}

func (buffer *EditableBuffer) Insert(char rune) {
	buffer.Content = buffer.contentBefore(buffer.CursorPosition) + string(char) + buffer.contentAfter(buffer.CursorPosition)
	buffer.CursorPosition++
}

func (buffer EditableBuffer) GetTermUiCompatibleOutput() string {
	output := ""
	for _, char := range buffer.Content {
		switch char {
		case '\t':
			output += "    "
		default:
			output += string(char)
		}
	}
	output = buffer.contentBefore(buffer.CursorPosition) + "|" + buffer.contentAfter(buffer.CursorPosition)
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
