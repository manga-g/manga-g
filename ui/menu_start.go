package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Search Manga", "Search Comics", "Exit"}

type MenuModel struct {
	Cursor int
	Choice string
}

func (m MenuModel) Init() tea.Cmd {
	// cmd := tea.EnterAltScreen
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.Choice = choices[m.Cursor]
			return m, tea.Quit

		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(choices) {
				m.Cursor = 0
			}

		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString("Choose an action to perform:\n\n")

	for i := 0; i < len(choices); i++ {
		if m.Cursor == i {
			s.WriteString("(x) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

// func main() {
// 	p := tea.NewProgram(model{})

// 	// Run returns the model as a tea.Model.
// 	m, err := p.Run()
// 	if err != nil {
// 		fmt.Println("Oh no:", err)
// 		os.Exit(1)
// 	}

// 	// Assert the final tea.Model to our local model and print the choice.
// 	if m, ok := m.(model); ok && m.choice != "" {
// 		fmt.Printf("\n---\nYou chose %s!\n", m.choice)
// 	}
// }
