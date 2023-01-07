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
	// fmt.Println("2.) Search Comics")
	fmt.Println("2.) Exit")
}

// StartMenu starts the menu loop
func StartMenu(basedApiUrl string) {

	for started := true; started; {
		if !started {
			break
		}
		DisplayOptions()
		backToMenu := false
		switch menuChoice := filechick.GetInput(); menuChoice {
		case "1":
			fmt.Print("Search for manga: ")
			query := filechick.GetInput()
			if query == "<" {
				backToMenu = true
			}
			if query != "<" {
				QueryCheck(query)
				query = url.QueryEscape(query)
				results := MkSearch(query)
				fmt.Println(results)
				started = false
			}

		// case "2":
		//	fmt.Print("Search for comic: ")
		//	query := filechick.GetInput()
		//	if query == "<" {
		//		backToMenu = true
		//	}
		//	if query != "<" {
		//		QueryCheck(query)
		//		query = url.QueryEscape(query)
		//		ComicSearch()
		//		started = false
		//	}
		//	if query == "<" {
		//		backToMenu = true
		//	}

		case "2":
			os.Exit(0)

		default:
			if !backToMenu {
				fmt.Println("Invalid Option")
			}
		}
	}
}
