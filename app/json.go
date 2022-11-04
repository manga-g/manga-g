package app

import (
	"encoding/json"
	"fmt"
)

// ParseData parses json data to struct
func ParseData(results string, data interface{}) {
	if results != "" {
		err := json.Unmarshal([]byte(results), &data)
		if err != nil {
			fmt.Println(results)
			fmt.Println("Error parsing json:", err)
			panic(err)
		}
	}
	if results == "" {
		panic("No results found")
	}
}

// ParseMangaSearch parses a list of manga search results from json to struct
func ParseMangaSearch(results string, manga *MangaList) {
	ParseData(results, manga)
}

func ParseChapters(results string, chapters *MangaChapters) {
	ParseData(results, chapters)
}

func ParseImages(results string, images *MangaImages) {
	ParseData(results, images)
}
