package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

/*
use cmd/manga-g/main.go as a reference for handling api data and user input
convert to the bubbletea model
press q or ctrl+c to quit
use j and k to move up and down
use enter to select
use tab to switch between search and manga list
create input box used for search
input box will be displayed at the top of the screen
create a list of manga
list will be displayed below the input box
create a list of chapters
list will be displayed below the manga list
download the selected chapter into a folder named after the manga and chapter number in the download folder
*/

// Model create a model for the program
type Model struct {
	// search   string
	// manga    app.Manga

	choices  []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	selected int
}

// Update create an update function for the program
func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter":
			m.selected = m.cursor
			if m.selected == 0 {
				// nothing

			} else if m.selected == 1 {
				// search for manga

			}

		}

	}

	return m, nil
}

// View create a view for the program
func View(m Model) string {

	return ""
}

// main entry point for the program
func main() {

	// create a model
	model := new(Model)
	// model.search = ""
	// model.manga = app.Manga{}

	model.choices = []string{"Search", "Manga"}
	model.cursor = 0

	model.selected = 0

	// create a program
	program := new(tea.Program)

	// start the program
	if err := program.Start(); err != nil {
		fmt.Println("failed to start program", err)
	}

}
