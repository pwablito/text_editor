package main

import (
	"log"

	"golang.org/x/term"
)

func GetWidowDimensions() (int, int) {
	width, height, err := term.GetSize(0)
	if err != nil {
		log.Fatal("couldn't get terminal size")
	}
	return width, height
}
