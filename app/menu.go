package app

import (
	"fmt"
	"net/url"
	"os"
)

// DisplayOptions displays Menu options
func DisplayOptions() {
	fmt.Println("Enter your choice: (1-3)")
	fmt.Println("1.) Search Manga")
	fmt.Println("2.) Search Comics")
	fmt.Println("3.) Exit")
}

// StartMenu starts the menu loop
func StartMenu(basedApiUrl string) {

	for started := true; started == true; {
		if !started {
			break
		}
		DisplayOptions()
		backToMenu := false
		switch menuChoice := GetInput(); menuChoice {
		case "1":
			fmt.Print("Search for manga: ")
			query := GetInput()
			if query == "<" {
				backToMenu = true
			}
			if query != "<" {
				QueryCheck(query)
				query = url.QueryEscape(query)
				MkSearch(basedApiUrl, query)
				started = false
			}

		case "2":
			fmt.Print("Search for comic: ")
			query := GetInput()
			if query == "<" {
				backToMenu = true
			}
			if query != "<" {
				QueryCheck(query)
				query = url.QueryEscape(query)
				ComicSearch()
				started = false
			}
			if query == "<" {
				backToMenu = true
			}

		case "3":
			os.Exit(0)
		default:
			if backToMenu == false {
				fmt.Println("Invalid Option")
			}
		}
	}
}
