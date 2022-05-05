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
				switch e.ID {
				case "<C-w>":
					return
				case "<C-s>":
					WriteBufferToFile(filename, buffer)
					break
				case "<Right>":
					if buffer.CursorPosition != len(buffer.Content) {
						buffer.CursorPosition++
						p.Text = buffer.GetTermUiCompatibleOutput()
						ui.Render(p)
					}
					break
				case "<Left>":
					if buffer.CursorPosition != 0 {
						buffer.CursorPosition--
						p.Text = buffer.GetTermUiCompatibleOutput()
						ui.Render(p)
					}
					break
				case "<Space>":
					buffer.Insert(rune(' '))
					p.Text = buffer.GetTermUiCompatibleOutput()
					ui.Render(p)
				case "<Enter>":
					buffer.Insert(rune('\n'))
					p.Text = buffer.GetTermUiCompatibleOutput()
					ui.Render(p)
				default:
					if len(e.ID) == 1 {
						buffer.Insert(rune(e.ID[0]))
						p.Text = buffer.GetTermUiCompatibleOutput()
						ui.Render(p)
					} else {
						log.Printf("Unhandled input: %s", e.ID)
					}
				}
			}

		case <-ticker:
			// Handle event loop ticker
		}
	}
}
