package widgets

import "github.com/gdamore/tcell/v2"

// Widget is the base interface for all widgets.
type Widget interface {
	Size() (x int, y int, width int, height int)
	SetSize(x int, y int, width int, height int)

	Key(event *tcell.EventKey)
	Draw(screen tcell.Screen)
}
