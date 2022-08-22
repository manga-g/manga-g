package app

import (
	"fmt"
	"net/url"
	"os"
)

func Init(basedApiUrl string) {
	CheckApi(basedApiUrl)
	fmt.Println("Enter your choice: (1-3)")
	fmt.Println("1.) Search Manga")
	fmt.Println("2.) Search Comics")
	fmt.Println("3.) Search Exit")
	menuChoice := GetInput()
	QueryCheck(menuChoice)
	if menuChoice == "1" {
		fmt.Print("Search for manga: ")
		query := GetInput()
		QueryCheck(query)
		query = url.QueryEscape(query)
		MkSearch(basedApiUrl, query)
	}
	if menuChoice == "2" {
		fmt.Print("Search for comic: ")
		query := GetInput()
		QueryCheck(query)
		query = url.QueryEscape(query)
		ComicSearch()
	}
	if menuChoice == "3" {
		os.Exit(0)
	}

	fmt.Println("Done")
}
