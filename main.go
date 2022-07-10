package main

import (
	"log"
	"os"
	buf_pkg "text_editor/src/buffer"
	"text_editor/src/window"
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
	buf_pkg.Logger = logger
	window.Logger = logger
	filename := os.Args[1]
	buffer := buf_pkg.ReadFileToBuffer(filename)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = filename
	p.Text = buffer.GetTermUiCompatibleOutput()

	window.SetWindowSize(p)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		rebuildUI := false
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.MouseEvent:
				logger.Print("Mouse inputs have not yet been implemented")
				switch e.ID {
				case "<MouseRelease>":
					logger.Print("Mouse release not implemented")
					break
				case "<MouseLeft>":
					logger.Print("Mouse left not implemented")
					break
				case "<MouseRight>":
					logger.Print("Mouse right not implemented")
					break
				default:
					logger.Printf("Unhandled mouse input: %s", e.ID)
				}
			case ui.ResizeEvent:
				window.SetWindowSize(p)
				rebuildUI = true
			case ui.KeyboardEvent:
				switch e.ID {
				case "<C-w>":
					// Exit the program
					return
				case "<C-s>":
					// Save the buffer
					buf_pkg.WriteBufferToFile(filename, buffer)
					break
				case "<C-r>":
					// Reload from file
					buffer = buf_pkg.ReadFileToBuffer(filename)
					rebuildUI = true
				case "<C-<Space>>":
					// Enter command mode
				case "<Right>":
					rebuildUI = buffer.ArrowRight()
				case "<Left>":
					rebuildUI = buffer.ArrowLeft()
				case "<Down>":
					rebuildUI = buffer.ArrowDown()
				case "<Up>":
					rebuildUI = buffer.ArrowUp()
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
