package main

import (
	"github.com/manga-g/manga-g/app"
)

// Entrypoint for the program.
func main() {
	app.Init("http://manga-api.bytecats.codes/")
}
