package app

import (
	"fmt"

	"github.com/byte-cats/filechick"
	"github.com/manga-g/manga-g/config"
)

// EndMessage prints a message when the app is finished
func EndMessage() {
	fmt.Println("\nManga-g has completed.\nStart program again to search for another manga.")
}

// Init app
func Init() {
	settings := config.NewSettings()
	// settings.SetApiUrl(config.GetEnvVar("api_url"))
	// config.SetDownloadPath(settings)
	CheckApi(settings.ApiUrl)
	StartMenu(settings.ApiUrl)
	EndMessage()
}

// Retry if there is no input, loop the request 3 times
func Retry() {
	var n = 0
	for ; n < 3; n++ {
		fmt.Println("Hint: Choose manga by corresponding number.\nPlease try again\n" + "Search for manga: ")
		query := filechick.GetInput()
		if query != "" {
			break
		}
	}
}

// QueryCheck checks if the query is empty
func QueryCheck(query string) {
	if query == "" || query == " " {
		Retry()
	}
}
