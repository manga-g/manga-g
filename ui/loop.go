package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type MangaModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func (m *MangaModel) Init() tea.Cmd {
	return nil
}
func InitModel() tea.Model {
	return &MangaModel{
		choices:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		selected: make(map[int]struct{}),
	}

}

func (m *MangaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+q":
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
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m *MangaModel) View() string {
	footer := "\nPress ctrl+c to quit.\n"
	header := "Enter a manga URL:\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		header += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	header += footer
	return header
}

func StartLoop() {
	program := tea.NewProgram(InitModel(), tea.WithAltScreen())
	err := program.Start()
	if err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, err.Error())
		if err2 != nil {
			return
		}
		os.Exit(1)
	}
}
