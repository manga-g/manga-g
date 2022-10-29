package widgets

import (
	"github.com/gdamore/tcell/v2"
)

var _ Widget = NewList()

type ListItem interface {
	String() string
}

// List represents a vertical list of selectable items.
type List struct {
	*Block
	ItemStyle         tcell.Style
	SelectedItemStyle tcell.Style

	selected int
	items    []ListItem
}

func NewList() *List {
	return &List{
		Block:             NewBlock(),
		SelectedItemStyle: tcell.StyleDefault.Foreground(tcell.ColorPurple),
	}
}

// Append appends an item to the list.
func (l *List) Append(item ListItem) {
	l.items = append(l.items, item)
}

// Remove removes an item from the list.
func (l *List) Remove(item ListItem) {
	for i, itm := range l.items {
		if item == itm {
			l.items = append(l.items[:i], l.items[i+1:]...)
		}
	}
}

// ---

func (w *List) Key(event *tcell.EventKey) {
	w.Block.Key(event)

	switch event.Key() {
	case tcell.KeyUp:
		w.selected--
		if w.selected < 0 {
			w.selected = 0
		}
	case tcell.KeyDown:
		w.selected++
		if w.selected >= len(w.items) {
			w.selected = len(w.items) - 1
		}
	}
}

func (w *List) Draw(screen tcell.Screen) {
	w.Block.Draw(screen)

	x, y, _, height := w.InnerSize()

	for i, item := range w.items {
		if i >= height {
			break
		}

		itemStyle := w.ItemStyle
		if w.selected == i {
			itemStyle = w.SelectedItemStyle
		}

		Print(screen, x, y, itemStyle, item.String())
		y++
	}
}
