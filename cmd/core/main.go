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
            fmt.Println("Invalid URL BRO try again.")
        }
    } else {
        fmt.Println("Could not connect to the internet.")
    }
}

func DoStuff(MangaG *app.MG, MangaUrl string) {
    //	fmt.Println("trying to grab manga from:", MangaUrl)
    //time.Sleep(time.Second * 8)
    //MangaG.SaveHtml(MangaUrl, "manga.html")
    //fmt.Println("Saved HTML")

    //fmt.Println("Attempting to detect Manga From Site...")
    //html := MangaG.LoadHtml("manga.html")
    //fmt.Println("Attempting to load HTML from a file...")

    html := MangaG.StringifyHtml(MangaUrl)
    //fmt.Println(html)
    fmt.Println("Html was loaded into memory")

    //	fmt.Println("Attempting to retrieve all manga pages from the site.")
    MangaG.NewDir("images")
    newFolderName := MangaG.TitleToDirName(MangaG.FindMangaTitle(html))
    // turn newFolderName to a string
    nfnString := string(newFolderName)
    // if nfnString length is longer than 10 characters then set nfnString to the first 10 characters
    nameLimit := 15

    if nameLimit < len(nfnString) {
        nfnString = nfnString[:nameLimit]
    }
    fmt.Println(nfnString)
    MangaG.NewDir("images/" + nfnString)

    //	found := MangaG.FindImageUrls(html)
    //	MangaG.DownloadAllWebp(found)
    //  MangaG.DownloadWebp(found[1], "images/001.webp")

    //MangaG.DeleteFile("manga.html")
    //fmt.Println("Deleted HTML no longer needed")

    //MangaG.DeleteFile("images/")
    //fmt.Println("Deleted images no longer needed")
}
