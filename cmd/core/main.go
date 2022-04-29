package main

import (
	"fmt"

	"manga-g/app"
)

// Entrypoint for the program.
func main() {
	MangaG := new(app.App)

	//	mangaUrl := "https://somemangasite.com/{mangaid}/{pagenumber}"
	//fmt.Println("Starting MangaG...")
	fmt.Println("Please Enter a URL for for a Manga's first page to download:")
	MangaUrl := MangaG.GetInput()
	fmt.Println("trying to grab manga from:", MangaUrl)

	//	time.Sleep(time.Second * 8)
	//	MangaG.SaveHtml(MangaUrl, "manga.html")
	//	fmt.Println("Saved HTML")

	//	fmt.Println("Attempting to detect Manga From Site...")
	//	html := MangaG.LoadHtml("manga.html")
	//	fmt.Println("Attempting to load HTML from a file...")

	html := MangaG.StringifyHtml(MangaUrl)
	//	fmt.Println("Html was loaded into memory")

	fmt.Println("Attempting to retrieve all manga pages from the site.")
	//fmt.Println("Got title from header:\n" + MangaG.FindMangaTitle(html))
	MangaG.NewDir("images")
	// trim last two chars from found[1]
	//found[1] = found[1][:len(found[1])-2]
	//fmt.Println("found[1]:", found[1])
	found := MangaG.FindImageUrls(html)
	MangaG.DownloadAllWebp(found)
	//MangaG.DownloadWebp(found[1], "images/001.webp")

	//MangaG.DeleteFile("manga.html")
	//	fmt.Println("Deleted HTML no longer needed")
	//MangaG.DeleteFile("images/")
	//fmt.Println("Deleted images no longer needed")
}
