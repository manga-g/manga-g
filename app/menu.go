package app

import (
	"fmt"
	"net/url"
	"os"

	"github.com/byte-cats/filechick"
)

// DisplayOptions displays Menu options
func DisplayOptions() {
	fmt.Println("Enter your choice: (1-2)")
	fmt.Println("1.) Search Manga")
	fmt.Println("2.) Exit")
}

// StartMenu starts the menu loop
// func StartMenu(basedApiUrl string) {
// 	for {
// 		DisplayOptions()
// 		menuChoice := filechick.GetInput()
// 		handleMenuChoice(menuChoice, basedApiUrl)
// 	}
// }

// handleMenuChoice processes the user's menu choice
func handleMenuChoice(choice, basedApiUrl string) {
	switch choice {
	case "1":
		searchManga(basedApiUrl)
	case "2":
		os.Exit(0)
	default:
		fmt.Println("Invalid Option")
	}
}

// searchManga handles the manga search functionality
func searchManga(basedApiUrl string) {
	fmt.Print("Search for manga: ")
	query := filechick.GetInput()
	if query == "<" {
		return // Back to menu
	}
	QueryCheck(query)
	query = url.QueryEscape(query)
	MkSearch(basedApiUrl, query)
}
