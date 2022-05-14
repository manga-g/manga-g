package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MangaModel struct {
	loader    spinner.Model
	textInput textinput.Model

	typing  bool
	loading bool

	choices  []string
	cursor   int
	selected map[int]struct{}

	err error
}

func (m *MangaModel) Init() tea.Cmd {
	return textinput.Blink
}
func InitModel() *MangaModel {
	return &MangaModel{
		typing:    true,
		loading:   false,
		loader:    spinner.New(),
		textInput: textinput.New(),

		choices:  []string{"1", "2", "3", "4", "5"},
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
				m.typing = false
				m.loading = true
				delete(m.selected, m.cursor)
				return m, spinner.Tick
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}
	}
	if m.typing {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	if m.loading {
		var cmd tea.Cmd
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *MangaModel) View() string {
	footer := "\nPress ctrl+c to quit.\n"
	if m.typing {
		header := "Enter a manga URL:\n"
		return header + m.textInput.View() + footer
	}
	header := "Manga List:\n"

	for i, choice := range m.choices {
		if m.loading {
			return m.loader.View() + "Please wait a moment...\n" + footer
		}
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

func StartProgram() {
	mod := InitModel()
	mod.textInput.Focus()
	mod.loader.Spinner = spinner.MiniDot
	program := tea.NewProgram(mod)

	err := program.Start()
	if err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, err.Error())
		if err2 != nil {
			return
		}
		os.Exit(1)
	}
}
