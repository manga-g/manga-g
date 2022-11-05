package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manga-g/manga-g/app"
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
*/

// Model create a model for the program
type Model struct {
	search   string
	manga    app.Manga
	selected int
}

// View create a view for the program
func View(m Model) string {
	// display the search box
	if m.search != "" {
		return fmt.Sprintf("Search: %s", m.search)
	}
	// 	 display the manga list
	// 	 return fmt.Sprintf("Manga: %s", m.manga.Name)
	// 	 display the chapter list
	// 	 return fmt.Sprintf("Chapters: %s", m.manga.Chapters)
	//
	// 	 display the image list
	// 	 return fmt.Sprintf("Images: %s", m.manga.Chapters)
	//
	// 	 display the chapter list
	// 	 return fmt.Sprintf("Chapters: %s", m.manga.Chapters)
	//
	// 	 display the image list
	// 	 return fmt.Sprintf("Images: %s", m.manga.Chapters)

	return ""
}

// Update create an update function for the program
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
	// basedApiUrl := "http://manga-api.bytecats.codes/"
	// app.Init(basedApiUrl)

	// handle search
	if m.search != "" {
		// handle search input
		// handle search results
		// handle manga selection
	}

	return m, nil
}

// main entry point for the program
func main() {

	// create a model
	model := new(Model)
	model.search = ""
	model.manga = app.Manga{}
	model.selected = 0

	// create a program
	program := new(tea.Program)

	// start the program
	if err := program.Start(); err != nil {
		fmt.Println("failed to start program", err)
	}
}
