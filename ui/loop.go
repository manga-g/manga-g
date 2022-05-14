package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct{}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.Key:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	return "Ok"
}

func StartLoop() {
	err := tea.NewProgram(&Model{}, tea.WithAltScreen()).Start()
	if err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, err.Error())
		if err2 != nil {
			return
		}
		os.Exit(1)
	}

}
