package app

import (
	"encoding/json"
	"fmt"
)

// ParseData parses json data to struct
func ParseData(results string, data interface{}) error {
	if results == "" {
		return fmt.Errorf("empty response from API")
	}

	// First check if this is a MangaChapters request by looking at the type
	if chapters, ok := data.(*MangaChapters); ok {
		err := json.Unmarshal([]byte(results), chapters)
		if err != nil {
			return err
		}
		return nil
	}

	// Regular handling for other types
	err := json.Unmarshal([]byte(results), &data)
	if err != nil {
		return err
	}
	return nil
}

// ParseMangaSearch parses manga search results from json to struct
func ParseMangaSearch(results string, manga *MangaList) {
	// Create a temporary response object
	var response struct {
		Result   string      `json:"result"`
		Response string      `json:"response"`
		Data     interface{} `json:"data"`
	}

	err := ParseData(results, &response)
	if err != nil {
		*manga = MangaList{} // Empty list on error
		return
	}

	// Check for errors in API response
	if response.Result != "ok" {
		*manga = MangaList{} // Empty list on error
		return
	}

	// Then extract and parse the data field
	dataJSON, err := json.Marshal(response.Data)
	if err != nil {
		*manga = MangaList{} // Empty list on error
		return
	}

	// Unmarshal into the manga list
	err = json.Unmarshal(dataJSON, manga)
	if err != nil {
		*manga = MangaList{} // Empty list on error
	}
}

// identifyJsonStructure attempts to figure out the structure of JSON for debugging
func identifyJsonStructure(jsonStr string) {
	var generic interface{}
	_ = json.Unmarshal([]byte(jsonStr), &generic)
}

// ParseChapters parses manga chapters from json to struct
func ParseChapters(results string, chapters *MangaChapters) {
	// The MangaChapters struct now matches the API response format,
	// so we can directly unmarshal
	err := json.Unmarshal([]byte(results), chapters)
	if err != nil {
		fmt.Printf("DEBUG: Failed to parse chapters: %v\n", err)
		return
	}

	// Check if we got valid results
	if chapters.Result != "ok" {
		fmt.Printf("DEBUG: API response result is not 'ok': %s\n", chapters.Result)
		return
	}

	// Process the data to fill compatibility fields even if there's no error
	chapters.Process()
}

// ParseImages parses a list of manga images from json to struct
func ParseImages(results string, images *MangaImages) {
	// First try to parse as AtHomeResponse
	var atHome AtHomeResponse
	err := ParseData(results, &atHome)
	if err != nil {
		return
	}

	// Process the AtHomeResponse data into MangaImages format
	*images = ProcessAtHomeResponse(atHome)
}

// ParseManga parses a manga from json to struct
func ParseManga(results string, manga *Manga) {
	err := ParseData(results, manga)
	if err != nil {
	}
}

// ParseMangaList parses a list of manga from json to struct
func ParseMangaList(results string, manga *MangaList) {
	ParseMangaSearch(results, manga)
}
