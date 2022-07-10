package main

import (
    "fmt"

    "github.com/manga-g/manga-g/prone"
)

// Entrypoint for the program.
func main() {
    fmt.Println("Manga Crawler Feature Test")
    //title := app.GetInput()

    //fmt.Println(extensions.NatoSearch(app, title))

    scraper := prone.NewScraper("www.manganato.com", "a[href]", "href", "")
    scraper.Scrape()

}
