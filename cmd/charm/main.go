package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manga-g/manga-g/app"
)

// use cmd/manga-g/main.go as a reference for handling api data and user input
// convert to the bubbletea model
// press q or ctrl+c to quit
// use j and k to move up and down
// use enter to select
// use tab to switch between search and manga list
// create input box used for search
// input box will be displayed at the top of the screen

// create a model for the program
type Model struct {
	search   string
	manga    app.Manga
	selected int
}

// create a view for the program
func View(m Model) string {
	// display the search box
	// display the manga list
	// display the chapter list
	// display the image list
	return ""
}

// create an update function for the program
func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	// handle user input
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	// handle api data
	// handle search
	// handle manga list
	// handle chapter list
	// handle image list
	// handle errors
	return m, nil
}

func main() {
	// create a model
	// create a program
	// start the program
	model := new(Model)
	model.search = ""
	model.manga = app.Manga{}
	model.selected = 0

	program := new(tea.Program)
	// program = tea.NewProgram(model, tea.WithAltScreen())

	// start the program
	if err := program.Start(); err != nil {
		fmt.Println("failed to start program", err)
	}
}
