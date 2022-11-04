package app

import (
	"encoding/json"
	"fmt"
)

// ParseData generic parsing function
func ParseData(results string, data *MangaImages) {
	err := json.Unmarshal([]byte(results), &data)
	if err != nil {
		fmt.Println("Error parsing json:", err)
	}
}

// ParseMangaSearch parses a list of manga search results from json to struct
func ParseMangaSearch(results string, manga *MangaList) {
	if results != "" {
		err := json.Unmarshal([]byte(results), &manga)
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

func ParseChapters(results string, chapters *MangaChapters) {
	ParseChapters(results, chapters)
}

func ParseImages(results string, images *MangaImages) {
	ParseData(results, images)
}
