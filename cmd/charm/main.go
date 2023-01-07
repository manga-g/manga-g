package main

import "github.com/manga-g/manga-g/ui"

/*
use cmd/manga-g/main.go as a reference for handling api data and user input
convert to the bubbletea model
press q or ctrl+c to quit
use j and k to move up and down
use enter to select
use tab to switch between search and manga list
create input box used for search
input box will be displayed at the top of the screen
create a list of manga
list will be displayed below the input box
create a list of chapters
list will be displayed below the manga list
download the selected chapter into a folder named after the manga and chapter number in the download folder
*/

func main() {
	ui.StartRenderingUI()
}
