package app

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseData parses json data to struct
func ParseData(results string, data interface{}) error {
	if results == "" {
		return fmt.Errorf("empty response from API")
	}

	err := json.Unmarshal([]byte(results), &data)
	if err != nil {
		fmt.Println("Raw JSON response:", results[:min(500, len(results))])
		fmt.Println("Error parsing json:", err)
		return err
	}
	return nil
}

// ParseMangaSearch parses manga search results from json to struct
func ParseMangaSearch(results string, manga *MangaList) {
	// First parse into a MangaDex response container
	var response struct {
		Result string      `json:"result"`
		Data   interface{} `json:"data"`
	}

	err := ParseData(results, &response)
	if err != nil {
		fmt.Println("DEBUG: Error parsing initial response:", err)
		*manga = MangaList{} // Empty list on error
		return
	}

	// Check for errors in API response
	if response.Result != "ok" {
		fmt.Println("DEBUG: API returned non-OK result:", response.Result)
		*manga = MangaList{} // Empty list on error
		return
	}

	// Debug output of the data structure
	fmt.Printf("DEBUG: Response structure: %T\n", response.Data)

	// Then extract and parse the data field
	dataJSON, err := json.Marshal(response.Data)
	if err != nil {
		fmt.Println("DEBUG: Error re-marshaling data:", err)
		*manga = MangaList{} // Empty list on error
		return
	}

	// Debug output of the JSON
	fmt.Println("DEBUG: Data JSON (first 200 chars):", string(dataJSON[:min(200, len(dataJSON))]))

	// Unmarshal into the manga list
	err = json.Unmarshal(dataJSON, manga)
	if err != nil {
		fmt.Println("DEBUG: Error parsing manga data:", err)
		// Try to identify the issue
		identifyJsonStructure(string(dataJSON))
		*manga = MangaList{} // Empty list on error
	}

	fmt.Printf("DEBUG: Successfully parsed %d manga items\n", len(*manga))
}

// identifyJsonStructure attempts to figure out the structure of JSON for debugging
func identifyJsonStructure(jsonStr string) {
	var generic interface{}
	err := json.Unmarshal([]byte(jsonStr), &generic)
	if err != nil {
		fmt.Println("DEBUG: Failed to parse JSON for structure identification:", err)
		return
	}

	switch data := generic.(type) {
	case []interface{}:
		fmt.Printf("DEBUG: JSON is an array with %d items\n", len(data))
		if len(data) > 0 {
			fmt.Printf("DEBUG: First item type: %T\n", data[0])
			// If it's a map, show the keys
			if m, ok := data[0].(map[string]interface{}); ok {
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				fmt.Println("DEBUG: Keys in first item:", strings.Join(keys, ", "))
			}
		}
	case map[string]interface{}:
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		fmt.Println("DEBUG: JSON is an object with keys:", strings.Join(keys, ", "))
	default:
		fmt.Printf("DEBUG: JSON is of type %T\n", data)
	}
}

// ParseChapters parses manga chapters from json to struct
func ParseChapters(results string, chapters *MangaChapters) {
	err := ParseData(results, chapters)
	if err != nil {
		fmt.Println("DEBUG: Error parsing chapters response:", err)
		return
	}

	// Process the data to fill compatibility fields
	chapters.Process()
	fmt.Printf("DEBUG: Processed %d chapters\n", len(chapters.Chapters))
}

// ParseImages parses a list of manga images from json to struct
func ParseImages(results string, images *MangaImages) {
	// First try to parse as AtHomeResponse
	var atHome AtHomeResponse
	err := ParseData(results, &atHome)
	if err != nil {
		fmt.Println("DEBUG: Error parsing at-home response:", err)
		return
	}

	// Debug output
	fmt.Printf("DEBUG: AtHome server: %s, Hash: %s, Pages: %d\n",
		atHome.BaseURL, atHome.Chapter.Hash, len(atHome.Chapter.DataSaver))

	// Process the AtHomeResponse data into MangaImages format
	*images = ProcessAtHomeResponse(atHome)
	fmt.Printf("DEBUG: Processed %d image URLs\n", len(*images))
}

// ParseManga parses a manga from json to struct
func ParseManga(results string, manga *Manga) {
	err := ParseData(results, manga)
	if err != nil {
		fmt.Println("DEBUG: Error parsing manga:", err)
	}
}

// ParseMangaList parses a list of manga from json to struct
func ParseMangaList(results string, manga *MangaList) {
	ParseMangaSearch(results, manga)
}
