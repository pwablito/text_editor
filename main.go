package main

import (
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var logger *log.Logger

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Missing filename argument")
	}
	log_file, _ := os.Create(".text_editor.log") // TODO catch errors
	logger = log.New(log_file, "", 0)
	filename := os.Args[1]
	buffer := ReadFileToBuffer(filename)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = filename
	p.Text = buffer.GetTermUiCompatibleOutput()

	setWindowSize(p)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		rebuildUI := false
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.MouseEvent:
				// Not handled yet
			case ui.ResizeEvent:
				setWindowSize(p)
				ui.Render(p)
			case ui.KeyboardEvent:
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
				case "<Down>":
					rebuildUI = buffer.ArrowDown()
				case "<Space>":
					buffer.Insert(rune(' '))
					rebuildUI = true
				case "<Enter>":
					buffer.Insert(rune('\n'))
					rebuildUI = true
				case "<Tab>":
					buffer.Insert(rune('\t'))
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
						logger.Printf("Unhandled input: %s", e.ID)
					}
				}
			}
		case <-ticker:
			// Handle event loop ticker
		}
		if rebuildUI {
			p.Text = buffer.GetTermUiCompatibleOutput()
			ui.Render(p)
		}
	}
}
