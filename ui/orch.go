package ui

import (
	"github.com/ayntgl/orchid/application"
	"github.com/ayntgl/orchid/widgets"
	"github.com/gdamore/tcell/v2"
)

type MangaListItem string
type MangaListItems []MangaListItem

func (manga MangaListItem) String() string {
	return string(manga)
}

func Orch() {
	block := widgets.NewBlock()
	block.Border = true
	block.Title = "Manga-G"
	block.Borders = widgets.RoundBorders
	block.BorderStyle = block.BorderStyle.Foreground(tcell.ColorPurple)
	block.TitleAlign = widgets.AlignCenter
	block.TitleStyle = block.TitleStyle.Foreground(tcell.ColorPurple)

	list := widgets.NewList()
	list.Title = "Manga List"
	list.ItemStyle = list.ItemStyle.Foreground(tcell.ColorPink)
	list.TitleStyle = list.TitleStyle.Foreground(tcell.ColorPurple)
	list.Border = true
	list.Borders = widgets.RoundBorders
	list.BorderStyle = list.BorderStyle.Foreground(tcell.ColorPink)
	list.TitleAlign = widgets.AlignCenter
	list.Block = block

	list.Append(MangaListItem("One Piece"))
	list.Append(MangaListItem("Bleach"))
	list.Append(MangaListItem("Naruto"))
	list.Append(MangaListItem("Dragon Ball"))
	list.Append(MangaListItem("Hunter X Hunter"))
	list.Append(MangaListItem("Tokyo Ghoul"))
	list.Append(MangaListItem("Black Clover"))
	list.Append(MangaListItem("Fairy Tail"))
	list.Append(MangaListItem("Boku no Hero Academia"))
	list.Append(MangaListItem("One Punch Man"))
	list.Append(MangaListItem("Attack on Titan"))

	tui := application.New(list)
	err := tui.Run()
	if err != nil {
		panic(err)
	}
}
