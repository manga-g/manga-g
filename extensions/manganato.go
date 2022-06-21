package extensions

import (
    "fmt"
    "regexp"

    "github.com/manga-g/manga-g/app"
)

const url = "https://manganato.com/search/story/"

// NatoParse a function to parse the html result to return only the div's with the class "search-story-item"
func NatoParse(html string) []string {
    var stories []string
    rx := regexp.MustCompile("<div class=\"search-story-item\">.*?</div>")

    for _, match := range rx.FindAllStringSubmatch(html, -1) {
        stories = append(stories, match[0])
    }
    return stories
}

// NatoSearch a function to search the manga in the Nato website
func NatoSearch(app *app.MG, query string) []string {
    html, err := app.CustomRequest(url + query)
    if err != nil {
        fmt.Println(err)
    }
    return NatoParse(html)
}
