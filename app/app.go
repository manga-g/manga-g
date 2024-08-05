package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Padding(1, 2)
	inputStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).BorderStyle(lipgloss.NormalBorder()).Padding(1, 2)
	messageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Background(lipgloss.Color("235")).Padding(1, 2)
)

type model struct {
	input    textinput.Model
	cursor   int
	choices  []string
	selected map[int]struct{}
	errorMsg string
}

func Init() {
	p := tea.NewProgram(model{
		input:    textinput.New(), // Updated to use New() instead of NewModel()
		choices:  []string{"Search Manga", "Exit"},
		selected: make(map[int]struct{}),
	})
	if err := p.Start(); err != nil { // Reverted to use Start() instead of Run()
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				// Trigger search functionality
				query := strings.TrimSpace(m.input.Value())
				err := QueryCheck(query)
				if err != nil {
					m.errorMsg = err.Error()
				} else {
					// Call your search function here
					// Example: SearchManga(query)
					m.input.Reset() // Reset input after search
				}
			} else {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Welcome to Manga-g!") + "\n")
	b.WriteString(inputStyle.Render(m.input.View()) + "\n")
	b.WriteString(messageStyle.Render("Press Ctrl+C or 'q' to quit.\n\n"))

	if m.errorMsg != "" {
		b.WriteString(messageStyle.Render(m.errorMsg + "\n\n"))
	}

	for i, choice := range m.choices {
		cursor := " " // Default cursor
		if m.cursor == i {
			cursor = ">" // Highlighted cursor
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}

	return b.String()
}

func EndMessage() {
	fmt.Println(messageStyle.Render("\nManga-g has completed.\nStart program again to search for another manga."))
}

// QueryCheck validates the user's query input
func QueryCheck(query string) error {
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("query cannot be empty, please enter a valid search term")
	}
	return nil
}
