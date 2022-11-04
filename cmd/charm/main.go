package main

import "fmt"

// make a model for the tui
type model struct {
	// add a string to the model
	text string
}

func main() {

	m := new(model)
	m.text = "Hello, World!"

	fmt.Println(m.text)
}
