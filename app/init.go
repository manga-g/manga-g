package app

import (
	"fmt"
)

func Init(basedApiUrl string) {
	CheckApi(basedApiUrl)
	StartMenu(basedApiUrl)
	fmt.Println("Done")
}
