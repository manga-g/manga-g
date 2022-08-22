package main

import (
	"github.com/manga-g/manga-g/app"
)

// Entrypoint for the program.
func main() {
	basedApiUrl := "http://manga-api.bytecats.codes/"

	app.Init(basedApiUrl)
}
