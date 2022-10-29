package application

import (
	"os"

	"github.com/ayntgl/orchid/widgets"
	"github.com/gdamore/tcell/v2"
)

// Application is a high-level API for writing terminal applications.
type Application struct {
	// Style is the default style to use when clearing the screen.
	Style tcell.Style

	screen tcell.Screen
	root   widgets.Widget
}

func New(root widgets.Widget) *Application {
	return &Application{
		root: root,
	}
}

func (app *Application) Run() error {
	var err error
	app.screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}

	err = app.screen.Init()
	if err != nil {
		return err
	}

	app.screen.SetStyle(app.Style)
	app.screen.Clear()

	width, height := app.screen.Size()
	app.root.SetSize(0, 0, width, height)

	for {
		app.root.Draw(app.screen)
		app.screen.Show()

		e := app.screen.PollEvent()
		switch e := e.(type) {
		case *tcell.EventResize:
			app.screen.Clear()

			width, height := app.screen.Size()
			app.root.SetSize(0, 0, width, height)
		case *tcell.EventKey:
			if e.Key() == tcell.KeyCtrlC {
				app.screen.Fini()
				os.Exit(0)
			}

			app.root.Key(e)
		}
	}
}
