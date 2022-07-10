package extensions

import (
    "fmt"
    "net/url"
    "regexp"

    "github.com/manga-g/manga-g/app"
)

const site = "https://manganato.com/search/story/"

// NatoParse a function to parse the html result to return only the div's with the class "search-story-item"
func NatoParse(html string) []string {
    var stories []string
    // get all sub-matches between <div class="search-story-item"> and </div>
    re := regexp.MustCompile(`<div class="search-story-item">(.*?)</div>`)
    for _, match := range re.FindAllStringSubmatch(html, -1) {
        stories = append(stories, match[0])
    }
    return stories
}

// NatoSearch a function to search the manga in the Nato website
func NatoSearch(app *app.MG, query string) []string {
    fmt.Println("Searching: " + site + url.QueryEscape(query))
    html, err := app.CustomRequest(site + url.QueryEscape(query))
    if err != nil {
        fmt.Println(err)
    }
    app.SaveHtml(site+query, "nato.html")
    return NatoParse(html)
}

func NatoHtmlTest(app *app.MG) {
    fmt.Println(NatoParse(app.LoadHtml("nato.html")))
}
