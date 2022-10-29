package widgets

import (
	"github.com/gdamore/tcell/v2"
)

var _ Widget = NewBlock()

// Block represents a frame with optional and customizable elements such as borders and a title.
type Block struct {
	Title      string
	TitleAlign Align
	TitleStyle tcell.Style

	Border      bool
	Borders     Borders
	BorderStyle tcell.Style

	x, y, width, height int
}

func NewBlock() *Block {
	return &Block{
		Border:  true,
		Borders: LineBorders,
	}
}

func (b *Block) InnerSize() (int, int, int, int) {
	x, y, width, height := b.Size()
	if b.Border {
		x++
		y++
		width -= 2
		height -= 2
	}

	return x, y, width, height
}

// ---

func (w *Block) Size() (int, int, int, int) {
	return w.x, w.y, w.width, w.height
}

func (w *Block) SetSize(x int, y int, width int, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

func (w *Block) Key(event *tcell.EventKey) {}

func (w *Block) Draw(screen tcell.Screen) {
	if w.Border {
		x, y, width, height := w.Size()
		width--
		height--

		// Top corner borders
		screen.SetContent(x, y, w.Borders.TopLeft, nil, w.BorderStyle)
		screen.SetContent(width, y, w.Borders.TopRight, nil, w.BorderStyle)
		// Bottom corner borders
		screen.SetContent(x, height, w.Borders.BottomLeft, nil, w.BorderStyle)
		screen.SetContent(width, height, w.Borders.BottomRight, nil, w.BorderStyle)

		// Horizontal borders
		for x := x + 1; x < width; x++ {
			screen.SetContent(x, y, w.Borders.Top, nil, w.BorderStyle)
			screen.SetContent(x, height, w.Borders.Bottom, nil, w.BorderStyle)
		}

		// Vertical borders
		for y := y + 1; y < height; y++ {
			screen.SetContent(x, y, w.Borders.Left, nil, w.BorderStyle)
			screen.SetContent(width, y, w.Borders.Right, nil, w.BorderStyle)
		}

		var tmpX int
		switch w.TitleAlign {
		case AlignLeft:
			tmpX = x + 1
		case AlignCenter:
			tmpX = (width - len(w.Title)) / 2
		case AlignRight:
			tmpX = width - len(w.Title)
		}

		for _, r := range w.Title {
			screen.SetContent(tmpX, y, r, nil, w.TitleStyle)
			tmpX++
		}
	}
}
