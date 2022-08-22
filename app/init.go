package app

import (
	"fmt"
	"net/url"
)

func Init(basedApiUrl string) {
	CheckApi(basedApiUrl)
	fmt.Print("Search for manga: ")
	query := GetInput()
	QueryCheck(query)
	query = url.QueryEscape(query)
	MkSearch(basedApiUrl, query)
	fmt.Println("Done")
}
