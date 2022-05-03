package main

import (
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/term"
)

func main() {
	if term.IsTerminal(0) {
		println("in a term")
	} else {
		println("not in a term")
	}
	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}
	println("width:", width, "height:", height)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = "Editor"
	p.Text = "Welcome to the editor"
	p.SetRect(0, 0, width, height)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return // (exit)
			}
		case <-ticker:
			// Handle event loop ticker
		}
	}
}
