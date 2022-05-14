package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func InitStyle() lipgloss.Style {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(0).
		PaddingLeft(4).
		Width(100)
	//fmt.Println(style.Render("Welcome to Manga G"))
	return style
}

func Render(style lipgloss.Style, str string) {
	fmt.Println(style.Render(str))
}
