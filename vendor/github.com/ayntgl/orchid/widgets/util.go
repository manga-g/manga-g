package widgets

import "github.com/gdamore/tcell/v2"

type Borders struct {
	Top         rune
	Bottom      rune
	Left        rune
	Right       rune
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
}

var LineBorders = Borders{
	Top:         tcell.RuneHLine,
	Bottom:      tcell.RuneHLine,
	Left:        tcell.RuneVLine,
	Right:       tcell.RuneVLine,
	TopLeft:     tcell.RuneULCorner,
	TopRight:    tcell.RuneURCorner,
	BottomLeft:  tcell.RuneLLCorner,
	BottomRight: tcell.RuneLRCorner,
}

var RoundBorders = Borders{
	Top:         LineBorders.Top,
	Bottom:      LineBorders.Bottom,
	Left:        LineBorders.Left,
	Right:       LineBorders.Right,
	TopLeft:     '╭',
	TopRight:    '╮',
	BottomLeft:  '╰',
	BottomRight: '╯',
}

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

func Print(screen tcell.Screen, x int, y int, style tcell.Style, text string) {
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}
