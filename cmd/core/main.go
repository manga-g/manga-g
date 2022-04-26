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
	fmt.Println("trying to grab manga from:", MangaUrl)

	MangaG.SaveHtml(MangaUrl)
	fmt.Println("Saved HTML")

	html := MangaG.LoadHtml("manga.html")
	fmt.Println(MangaG.FindMangaTitle(html))

	ImageUrl := MangaG.FindImageUrl(html)
	fmt.Println(MangaG.FindImageKey(ImageUrl))
	cycleImages(MangaG, ImageUrl, html)
	MangaG.DeleteHtml("manga.html")
	fmt.Println("Deleted HTML no longer needed")
}

func cycleImages(starter *app.App, imageurl string, html string) {
	pagecount := starter.GetPageCount(html)
	for i := 1; i < pagecount; i++ {
		starter.SaveImage(imageurl)
		time.Sleep(time.Second * 8)

		imageurl = starter.IncrementImageUrl(imageurl)

		fmt.Println("Saving Page: ", i)

		time.Sleep(time.Second * 8)
	}
}
