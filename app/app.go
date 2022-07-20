package app

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetInput Function to get user input from the command line
func GetInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return ""
	}
	return input
}

// ValidateUrl checks if url is valid
func ValidateUrl(UrlToCheck string) bool {
	_, err := url.ParseRequestURI(UrlToCheck)
	if err != nil {
		return false
	}
	return true
}

// StringToInt to change string to int
func StringToInt(str string) int {
	var i int
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return 0
	}
	return i
}

type MangaInfo []struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// create a function to parse json into struct
func ParseMangaSearch(results string, manga *MangaInfo) MangaInfo {

	err := json.Unmarshal([]byte(results), &manga)
	if err != nil {
		fmt.Println("Error parsing json:", err)
	}

	return *manga
}
