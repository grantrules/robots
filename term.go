package main

import (
	"fmt"

	"golang.org/x/term"
)

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(0)
	if err != nil {
		// handle error
	}
	return width, height
}

// print character at position x, y in the terminal
func printAt(x, y int, char string) {
	fmt.Printf("\033[%d;%dH%s", y, x, char)
}

// clear the terminal
func clear() {
	fmt.Print("\033[H\033[2J")
}

// create dialog surrounded with a border using unicode characters in terminal in the center with width and height
func dialog(x, y, width, height int, title string) {
	// top border
	printAt(x, y, "╔")
	for i := 1; i < width-1; i++ {
		printAt(x+i, y, "═")
	}
	printAt(x+width-1, y, "╗")

	// sides
	for i := 1; i < height-1; i++ {
		printAt(x, y+i, "║")
		printAt(x+width-1, y+i, "║")
	}

	// title
	printAt(x+1, y+1, title)

	// bottom border
	printAt(x, y+height-1, "╚")
	for i := 1; i < width-1; i++ {
		printAt(x+i, y+height-1, "═")
	}
	printAt(x+width-1, y+height-1, "╝")
}
