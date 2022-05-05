package main

import (
	"io/ioutil"
	"log"
	"os"
)

type EditableBuffer struct {
	Content        string
	CursorPosition int
}

func ReadFileToBuffer(filename string) EditableBuffer {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Couldn't read file: %v", err)
	}

	text := string(content)
	return EditableBuffer{
		Content:        text,
		CursorPosition: 0,
	}
}

func WriteBufferToFile(filename string, buffer EditableBuffer) {
	// TODO use existing file permissions if file exists
	err := os.WriteFile(filename, []byte(buffer.Content), 0644)
	if err != nil {
		log.Fatalf("Couldn't write file: %v", err)
	}
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
	return output
}
