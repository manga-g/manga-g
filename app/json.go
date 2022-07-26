package app

import (
	"encoding/json"
	"fmt"
)

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
	err := json.Unmarshal([]byte(results), &chapters)
	if err != nil {
		fmt.Println("Error parsing json:", err)
	}
}

func ParseImages(results string, images *MangaImages) {
	err := json.Unmarshal([]byte(results), &images)
	if err != nil {
		fmt.Println("Error parsing json:", err)
	}
}
