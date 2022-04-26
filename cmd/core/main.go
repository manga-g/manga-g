package main

import (
	"fmt"
	"time"

	"manga-g/app"
)

// Entrypoint for the program.
func main() {
	MangaG := new(app.App)

	//	mangaUrl := "https://somemangasite.com/{mangaid}/{pagenumber}"
	fmt.Println("Starting MangaG...")
	fmt.Println("Please Enter a URL for for a Manga's first page to download:")
	MangaUrl := MangaG.GetInput()
	fmt.Println("trying to grab manga from:", MangaUrl)

	time.Sleep(time.Second * 8)
	MangaG.SaveHtml(MangaUrl, "manga.html")
	fmt.Println("Saved HTML")

	fmt.Println("Attempting to detect Manga From Site...")
	html := MangaG.LoadHtml("manga.html")
	fmt.Println("Attempting to load HTML from a file...")

	//	html := MangaG.StringifyHtml(MangaUrl)
	//	fmt.Println("Html was loaded into memory")

	fmt.Println("Got title from header:\n" + MangaG.FindMangaTitle(html))
	MangaG.NewDir("images")
	ImageUrl, imgErr := MangaG.FindImageUrl(html)
	if imgErr != nil {
		fmt.Println("Error:", imgErr)
	} else {
		fmt.Println("Found Image URL:", ImageUrl[1])
		MangaG.CycleImages(ImageUrl, MangaG.GetPageCount(html))

		//	fmt.Println("Attempting to retrieve all manga pages from the site.")
		//cycleImages(MangaG, ImageUrl[1], MangaG.GetPageCount(html))
	}
	//MangaG.DeleteFile("manga.html")
	//fmt.Println("Deleted HTML no longer needed")
	//MangaG.DeleteFile("images/")
	//fmt.Println("Deleted images no longer needed")
}
