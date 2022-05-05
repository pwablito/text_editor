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
	p.Text = buffer.GetTermUiCompatibleOutput()
	p.SetRect(0, 0, width, height)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.MouseEvent:
				// Not handled yet
			case ui.ResizeEvent:
				// Not handled yet
			case ui.KeyboardEvent:
				rebuildUI := false
				switch e.ID {
				case "<C-w>":
					return
				case "<C-s>":
					WriteBufferToFile(filename, buffer)
					break
				case "<C-r>":
					// Reload from file
					buffer = ReadFileToBuffer(filename)
					rebuildUI = true
				case "<Right>":
					rebuildUI = buffer.ArrowRight()
				case "<Left>":
					rebuildUI = buffer.ArrowLeft()
				case "<Space>":
					buffer.Insert(rune(' '))
					rebuildUI = true
				case "<Enter>":
					buffer.Insert(rune('\n'))
					rebuildUI = true
				case "<Backspace>":
					rebuildUI = buffer.Backspace()
				case "<Delete>":
					rebuildUI = buffer.Delete()
				default:
					if len(e.ID) == 1 {
						buffer.Insert(rune(e.ID[0]))
						rebuildUI = true
					} else {
						log.Printf("Unhandled input: %s", e.ID)
					}
				}
				if rebuildUI {

					p.Text = buffer.GetTermUiCompatibleOutput()
					ui.Render(p)
				}
			}
		case <-ticker:
			// Handle event loop ticker
		}
	}
}
