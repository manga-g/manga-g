package app

import (
	"fmt"
)

// TODO: define manga-g app object as a struct

// Retry if there is no input, loop the request 3 times
func Retry(query string) {
	var n = 0
	for ; n < 3; n++ {
		fmt.Println("Hint: Choose manga by corresponding number.\nPlease try again\n" + "Search for manga: ")
		query = GetInput()
		if query != "" {
			break
		}
	}
}
func QueryCheck(query string) {
	if query == "" || query == " " {
		Retry(query)
	}
}
