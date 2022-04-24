package main

import (
    "fmt"
    "time"

    "manga-g/app"
)

// Entrypoint for the program.
func main() {
    //var starter *app.App
    //MagnaUrl := "https://somemangasite.com/{mangaid}/{pagenumber}"
    //starter.SaveHtml(MangaUrl)

    //html := starter.LoadHtml("manga.html")
    //fmt.Println(starter.FindMangaTitle(html))
    //imageurl := starter.FindImageUrl(html)
    //fmt.Println(starter.FindImageKey(imageurl))
    //cycleImages(starter, imageurl, html)

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
