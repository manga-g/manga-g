package main

import (
	"fmt"
	"strconv"

	"github.com/manga-g/manga-g/app"
	_ "github.com/manga-g/manga-g/app"
)

// Entrypoint for the program.
func main() {
	apiUrl := "http://localhost:3000/"
	fmt.Print("Search for manga: ")
	query := app.GetInput()
	apiUrl += "nato/search?q=" + query
	res, _ := app.CustomRequest(apiUrl)
	var mangaList app.MangaInfo
	app.ParseMangaSearch(res, &mangaList)
	titles := []string{}
	for i, manga := range mangaList {
		titles = append(titles, fmt.Sprintf("%d. %s", i+1, manga.Title))
	}

	number := len(titles)
	SelectMessage := "Select a title: (1 - " + strconv.Itoa(number) + ")"
	fmt.Println(SelectMessage)
	for _, title := range titles {
		fmt.Println(title)
	}
	// mangaChoice := app.GetInput()

}
