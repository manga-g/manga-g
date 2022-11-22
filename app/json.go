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

// ParseChapters parses a list of manga chapters from json to struct
func ParseChapters(results string, chapters *MangaChapters) {
	ParseData(results, chapters)
}

// ParseImages parses a list of manga images from json to struct
func ParseImages(results string, images *MangaImages) {
	ParseData(results, images)
}

// ParseManga parses a manga from json to struct
func ParseManga(results string, manga *Manga) {
	ParseData(results, manga)
}

// ParseMangaList parses a list of manga from json to struct
func ParseMangaList(results string, manga *MangaList) {
	ParseData(results, manga)
}
