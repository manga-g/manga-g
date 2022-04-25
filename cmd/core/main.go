package main

import (
	"fmt"
	"time"

	"manga-g/app"
)

// Entrypoint for the program.
func main() {
<<<<<<< HEAD
	//var starter *app.App
	//url := "https://somemangasite.com/mangaid/pagenumber"
	//MangaG.SaveHtml(url)
=======
    //var starter *app.App
    //MagnaUrl := "https://somemangasite.com/{mangaid}/{pagenumber}"
    //starter.SaveHtml(MangaUrl)
>>>>>>> a595c9b1b3efa46bd391b109b445924772864200

	//html := MangaG.LoadHtml("manga.html")
	//fmt.Println(MangaG.FindMangaTitle(html))
	//imageurl := MagnaG.FindImageUrl(html)
	//fmt.Println(MangaG.FindImageKey(imageurl))
	//cycleImages(MangaG, imageurl, html)

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
