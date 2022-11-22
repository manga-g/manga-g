package app

import (
	"fmt"

	"github.com/byte-cats/filechick"
)

// TODO: define manga-g app object as a struct

// EndMessage prints a message when the app is finished
func EndMessage() {
	fmt.Println("\nManga-g has completed.\nStart program again to search for another manga.")
}

// Init app
func Init(basedApiUrl string) {
	CheckApi(basedApiUrl)
	StartMenu(basedApiUrl)
	EndMessage()
}

// Retry if there is no input, loop the request 3 times
func Retry(query string) {
	var n = 0
	for ; n < 3; n++ {
		fmt.Println("Hint: Choose manga by corresponding number.\nPlease try again\n" + "Search for manga: ")
		query = filechick.GetInput()
		if query != "" {
			break
		}
	}
}

// QueryCheck checks if the query is empty
func QueryCheck(query string) {
	if query == "" || query == " " {
		Retry(query)
	}
}
