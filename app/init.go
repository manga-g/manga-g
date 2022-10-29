package app

import (
	"fmt"
)

func EndProgram() {
	fmt.Println("\nManga-g has completed.\nStart program again to search for another manga.")
}

func Init(basedApiUrl string) {
	CheckApi(basedApiUrl)
	StartMenu(basedApiUrl)
	EndProgram()
}
