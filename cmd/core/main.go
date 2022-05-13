package main

import (
	"fmt"

	"manga-g/app"
)

// Entrypoint for the program.
func main() {
	MangaG := new(app.MG)
	fmt.Println("Starting MangaG...")
	if MangaG.Connected() {
		fmt.Println("Please Enter a URL for a Manga's first page to download:")
		MangaUrl := MangaG.GetInput()
		if MangaG.ValidateUrl(MangaUrl) {
			DoStuff(MangaG, MangaUrl)
		} else {
			fmt.Println("Invalid URL please try again.")
			fmt.Println("Example: https://www.mangaeden.com/en/en-manga/one-piece/")
			fmt.Println("Exiting...")
		}
	} else {
		fmt.Println("Could not connect to the internet.")
		fmt.Println("Exiting...")
	}
}
