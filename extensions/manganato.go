package extensions

import (
    "regexp"
    "strings"

    "github.com/manga-g/manga-g/app"
)

const url = "https://manganato.com/search/story/"

// NatoSearch a function to search the manga in the Nato website
func NatoSearch(app *app.MG, query string) string {

    html := app.StringifyHtml(url + query)
    return html
}

// NatoParse a function to parse the html result to return only the div's with the class "search-story-item"
func NatoParse(html string) []string {
    var stories []string
    rx := regexp.MustCompile("<div class=\"search-story-item\">.*?</div>")

    for _, match := range rx.FindAllStringSubmatch(html, -1) {
        match[0] = strings.ReplaceAll(match[0], "'", "")
        match[0] = strings.ReplaceAll(match[0], "\"", "")
        match[0] = strings.ReplaceAll(match[0], "\\", "")
        match[0] = strings.ReplaceAll(match[0], " ", "")
        match[0] = strings.ReplaceAll(match[0], "\n", "")
        match[0] = strings.ReplaceAll(match[0], "\t", "")
        match[0] = strings.ReplaceAll(match[0], "'", "")
        match[0] = strings.ReplaceAll(match[0], "alt=", "")
        match[0] = strings.ReplaceAll(match[0], ">", "")
        //match[0] = match[0][:len(match[0])-2]

        stories = append(stories, match[0])
    }
    return stories
}
