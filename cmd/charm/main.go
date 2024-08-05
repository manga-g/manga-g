package main

import (
	"fmt"
	"log"
	"strings"

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
	choices  []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	selected int      // selected item index
	input    string   // input for search
}

// Init add Init method to Model to satisfy tea.Model interface
func (m Model) Init() tea.Cmd {
	return nil // No initial command
}

// Update create an update function for the program
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
				// Handle search
				query := strings.TrimSpace(m.input)
				if query != "" {
					// Implement search functionality for multiple letters
					log.Printf("Searching for: %s", query)
				}
			} else if m.selected == 1 {
				// Handle manga selection
				log.Println("Manga selected")
			}
		default:
			m.input += msg.String() // Accumulate input for multiple letters
		}
	}
	return m, cmd
}

// View add View method to Model to satisfy tea.Model interface
func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.input) // Display input
	b.WriteString("\n")
	for i, choice := range m.choices {
		cursor := " " // No cursor
		if m.cursor == i {
			cursor = ">" // Cursor for selected item
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}

// main entry point for the program
func main() {
	// create a model
	model := Model{
		choices: []string{"Search", "Manga"},
		cursor:  0,
		input:   "",
	}

	// create a program
	program := tea.NewProgram(model)

	// start the program
	if err := program.Start(); err != nil {
		log.Fatalf("failed to start program: %v", err)
	}
}
