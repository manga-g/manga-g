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
	fmt.Println("Please Enter a URL to download:")
	MangaUrl := MangaG.GetInput()
	//fmt.Println("trying to grab manga from:", MangaUrl)

	//MangaG.SaveHtml(MangaUrl)
	//time.Sleep(time.Second * 8)
	//fmt.Println("Saved HTML")

	fmt.Println("Attempting to detect Manga From Site...")

	//	html := MangaG.LoadHtml("manga.html")
	//	fmt.Println("Attempting to load HTML from a file...")

	html := MangaG.StringifyHtml(MangaUrl)
	fmt.Println("Html was loaded into memory")

	fmt.Println("Got title from header:\n" + MangaG.FindMangaTitle(html))

	ImageUrl, imgerr := MangaG.FindImageUrl(html)
	if imgerr != nil {
		fmt.Println("Error:", imgerr)
	} else {
		fmt.Println("Found Image URL:", ImageUrl)
	}

	//fmt.Println(MangaG.FindImageKey(ImageUrl))

	//fmt.Println("Attempting to retrieve all manga pages from the site.")
	//	cycleImages(MangaG, ImageUrl, html)

	//MangaG.DeleteFile("manga.html")
	//fmt.Println("Deleted HTML no longer needed")
	//MangaG.DeleteFile("images/")
	//fmt.Println("Deleted images no longer needed")
}

func cycleImages(starter *app.App, imageurl string, html string) {
	pagecount := starter.GetPageCount(html)
	for i := 1; i < pagecount; i++ {
		starter.SaveImage(imageurl)
		time.Sleep(time.Second * 6)

		imageurl = starter.IncrementImageUrl(imageurl)

		fmt.Println("Saving Page: ", i)

		time.Sleep(time.Second * 6)
	}
}
