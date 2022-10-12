package ui

//
//import (
//    tea "github.com/charmbracelet/bubbletea"
//)
//
//type model struct {
//    choices  []string         // items on the to-do list
//    cursor   int              // which to-do list item our cursor is pointing at
//    selected map[int]struct{} // which to-do items are selected
//    input    string           // the user's input
//}
//
//func initialModel(array []string) model {
//    return model{
//        choices:  array,
//        cursor:   0,
//        selected: make(map[int]struct{}),
//    }
//}
//
//func (m *model) Add() {
//    m.choices = append(m.choices, m.input)
//    m.input = ""
//}
//
//func (m *model) Remove() {
//    if len(m.selected) == 0 {
//        return
//    }
//    newChoices := make([]string, 0, len(m.choices)-len(m.selected))
//    for i, choice := range m.choices {
//        if _, ok := m.selected[i]; !ok {
//            newChoices = append(newChoices, choice)
//        }
//    }
//    m.choices = newChoices
//    m.selected = make(map[int]struct{})
//}
//
//func (m *model) Up() {
//    if m.cursor > 0 {
//        m.cursor--
//    }
//}
//
//func (m *model) Down() {
//    if m.cursor < len(m.choices)-1 {
//        m.cursor++
//    }
//}
//
//func (m model) Init() tea.Cmd {
//    // Just return `nil`, which means "no I/O right now, please."
//    return nil
//}
//
//// create input
//func (m *model) Input() {
//
//}
//
//// View
//func (m *model) View() string {
//    // The header
//    s := "Manga Choice:\n\n"
//    // Iterate over our choices
//    for i, choice := range m.choices {
//        // Is the cursor pointing at this choice?
//        cursor := " " // no cursor
//        if m.cursor == i {
//            cursor = ">" // cursor!
//        }
//
//    }
//    // The footer
//    s += "\n\n"
//    return s
//}
//
//// Update
//func (m *model) Update() {
//
//}
