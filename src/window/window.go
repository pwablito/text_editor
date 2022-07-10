package window

import (
	"log"

	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/term"
)

var Logger *log.Logger

func GetWidowDimensions() (int, int) {
	width, height, err := term.GetSize(1)
	if err != nil {
		log.Fatal("couldn't get terminal size")
	}
	return width, height
}

func SetWindowSize(p *widgets.Paragraph) {
	width, height := GetWidowDimensions()
	p.SetRect(0, 0, width, height)
}
