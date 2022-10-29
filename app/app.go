package app

import (
	"fmt"
)

// Retry if there is no input, loop the request 3 times
func Retry(query string) {
	var n = 0
	for ; n < 3; n++ {
		fmt.Println("You should choose the number corresponding to the manga you want to read.\nTry again,please :)\n" + "Search for manga: ")
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
