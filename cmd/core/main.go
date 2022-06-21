package main

import (
    "fmt"

    "github.com/manga-g/manga-g/app"
    "github.com/manga-g/manga-g/extensions"
)

// Entrypoint for the program.
func main() {
    app := new(app.MG)
    fmt.Println("Manga Title Name?")
    title := app.GetInput()
    fmt.Println(extensions.NatoSearch(app, title))

}
