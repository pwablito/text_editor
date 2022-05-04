package main

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing filename argument")
	}
	filename := os.Args[1]
	buffer := ReadFileToBuffer(filename)

	width, height := GetWidowDimensions()
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = filename
	p.Text = buffer.Content
	p.SetRect(0, 0, width, height)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<C-w>":
				return // (exit)
			case "<C-s>":
				WriteBufferToFile(filename, buffer)
			}
		case <-ticker:
			// Handle event loop ticker
		}
	}
}
