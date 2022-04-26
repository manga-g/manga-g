package main

import (
    "fmt"

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

    //MangaG.SaveHtml(MangaUrl)
    //time.Sleep(time.Second * 8)
    //fmt.Println("Saved HTML")

    fmt.Println("Attempting to detect Manga From Site...")
    //	html := MangaG.LoadHtml("manga.html")
    //	fmt.Println("Attempting to load HTML from a file...")

    html := MangaG.StringifyHtml(MangaUrl)
    fmt.Println("Html was loaded into memory")

    fmt.Println("Got title from header:\n" + MangaG.FindMangaTitle(html))

    ImageUrl, imgErr := MangaG.FindImageUrl(html)
    if imgErr != nil {
        fmt.Println("Error:", imgErr)
    } else {
        fmt.Println("Found Image URL:", ImageUrl[1])
        //fmt.Println(MangaG.FindImageKey(ImageUrl[1]))

        fmt.Println("Attempting to retrieve all manga pages from the site.")
        cycleImages(MangaG, ImageUrl[1], html)
    }

    //MangaG.DeleteFile("manga.html")
    //fmt.Println("Deleted HTML no longer needed")
    //MangaG.DeleteFile("images/")
    //fmt.Println("Deleted images no longer needed")
}

func cycleImages(MangaG *app.App, ImageUrl string, html string) {
    //wg := new(sync.WaitGroup)
    count := MangaG.GetPageCount(html)
    for i := 1; i < count; i++ {
        //	wg.Add(1)

        //go func(MangaG *app.App, ImageUrl string, i int, wg *sync.WaitGroup) {
        fmt.Println("Attempting to download page:", i)
        MangaG.SaveImage(ImageUrl)
        //time.Sleep(time.Second * 8)
        ImageUrl = MangaG.IncrementImageUrl(ImageUrl)
        fmt.Println("Saving Page: ", i)
        //	wg.Done()

        //}(MangaG, ImageUrl, i, wg)

        //time.Sleep(time.Second * 10)
    }
    //wg.Wait()
    fmt.Println("Finished Saving Images")
}
