package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Variable to hold mock response data for internal use
var mockResponse struct {
	Result   string      `json:"result"`
	Response string      `json:"response"`
	Data     interface{} `json:"data"`
}

// MockSearchManga simulates searching for manga without API calls.
func MockSearchManga(query string, verbose bool) {
	fmt.Printf("MOCK MODE: Searching for: %s\n", query)

	// Sample manga search response
	mockData := `{ "result": "ok", "response": "collection", "data": [ { "id": "mock-manga-id-1", "type": "manga", "attributes": { "title": {"en": "Mock Manga 1"}, "altTitles": [{"ja": "モックマンガ 1"}], "description": {"en": "This is a mock manga for testing purposes"}, "status": "ongoing", "year": 2023, "tags": [ {"id": "tag-1", "type": "tag", "attributes": {"name": {"en": "Action"}}}, {"id": "tag-2", "type": "tag", "attributes": {"name": {"en": "Adventure"}}} ] }, "relationships": [] }, { "id": "mock-manga-id-2", "type": "manga", "attributes": { "title": {"en": "Mock Manga 2"}, "altTitles": [], "description": {"en": "Another mock manga with a longer description that shows how the output is truncated when it exceeds a certain length"}, "status": "completed", "year": 2022, "tags": [ {"id": "tag-3", "type": "tag", "attributes": {"name": {"en": "Drama"}}} ] }, "relationships": [] } ], "limit": 2, "offset": 0, "total": 2 }`

	if verbose {
		// Assuming printFormattedJSON remains in main or is also moved/made public
		// For now, just print raw mock data if verbose
		fmt.Println("Raw Mock Search Data:", mockData)
	}

	var mangaList MangaList
	err := json.Unmarshal([]byte(mockData), &mockResponse)
	if err != nil {
		fmt.Printf("Error parsing mock data: %v\n", err)
		return
	}

	dataJSON, _ := json.Marshal(mockResponse.Data)
	json.Unmarshal(dataJSON, &mangaList)

	fmt.Printf("Found %d manga results:\n", len(mangaList))
	for i, manga := range mangaList {
		title := mangaList.GetTitle(i)
		fmt.Printf("[%d] ID: %s - Title: %s\n", i+1, manga.ID, title)
		// Print tags and description (assuming min function is available or replaced)
		if len(manga.Attributes.Tags) > 0 { /* ... */
		}
		if desc, ok := manga.Attributes.Description["en"]; ok && desc != "" { /* ... */
		}
		fmt.Println()
	}

	if len(mangaList) > 0 {
		fmt.Print("Select a manga to view details (enter number, or 0 to exit): ")
		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil || choice <= 0 || choice > len(mangaList) {
			fmt.Println("Invalid selection or exiting.")
			return
		}
		selectedID := mangaList[choice-1].ID
		fmt.Println("\n--- Fetching mock details for selected manga ---")
		MockMangaDetails(selectedID, verbose)
	}
}

// MockMangaDetails simulates getting manga details without API calls.
func MockMangaDetails(mangaID string, verbose bool) {
	fmt.Printf("MOCK MODE: Getting manga details for ID: %s\n", mangaID)

	// Sample manga details response
	mockData := `{ "result": "ok", "response": "entity", "data": { "id": "` + mangaID + `", "type": "manga", "attributes": { "title": {"en": "Mock Manga Details"}, "altTitles": [{"ja": "モックマンガ詳細"}], "description": {"en": "This is a detailed description of a mock manga for testing the API parsing functions."}, "status": "ongoing", "year": 2023, "contentRating": "safe", "tags": [ {"id": "tag-1", "type": "tag", "attributes": {"name": {"en": "Action"}}}, {"id": "tag-2", "type": "tag", "attributes": {"name": {"en": "Adventure"}}} ] }, "relationships": [ {"id": "author-1", "type": "author", "attributes": {"name": "Mock Author"}} ] } }`

	// Chapter feed mock data
	mockChapterData := `{ "result": "ok", "response": "collection", "data": [ {"id": "mock-chapter-1", "type": "chapter", "attributes": { "volume": "1", "chapter": "1", "title": "Mock Chapter 1", "translatedLanguage": "en", "publishAt": "2023-01-01T00:00:00Z", "createdAt": "2023-01-01T00:00:00Z", "updatedAt": "2023-01-01T00:00:00Z" } }, {"id": "mock-chapter-2", "type": "chapter", "attributes": { "volume": "1", "chapter": "2", "title": "Mock Chapter 2", "translatedLanguage": "en", "publishAt": "2023-01-02T00:00:00Z", "createdAt": "2023-01-02T00:00:00Z", "updatedAt": "2023-01-02T00:00:00Z" } }, {"id": "mock-chapter-3", "type": "chapter", "attributes": { "volume": "1", "chapter": "3", "title": "Mock Chapter 3", "translatedLanguage": "en", "publishAt": "2023-01-03T00:00:00Z", "createdAt": "2023-01-03T00:00:00Z", "updatedAt": "2023-01-03T00:00:00Z" } } ], "limit": 3, "offset": 0, "total": 3 }`

	if verbose {
		fmt.Println("Mock Manga Details:", mockData)
	}

	var response struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				Title       map[string]string `json:"title"`
				Description map[string]string `json:"description"`
				Status      string            `json:"status"`
				Year        int               `json:"year"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(mockData), &response); err != nil {
		fmt.Printf("Error parsing mock JSON: %v\n", err)
		return
	}

	fmt.Printf("ID: %s\n", response.Data.ID)
	// Assuming getFirstValue is available or replaced
	// fmt.Printf("Title: %s\n", getFirstValue(response.Data.Attributes.Title))
	fmt.Printf("Title: %s\n", "Mock Title") // Placeholder
	fmt.Printf("Status: %s\n", response.Data.Attributes.Status)
	fmt.Printf("Year: %d\n", response.Data.Attributes.Year)
	// fmt.Printf("Description: %s\n", getFirstValue(response.Data.Attributes.Description))
	fmt.Printf("Description: %s\n", "Mock Description") // Placeholder

	// Process chapters
	fmt.Println("\nMock Chapters:")
	if verbose {
		fmt.Println("Mock Chapter Data:", mockChapterData)
	}

	var mangaChapters MangaChapters // Use type directly from app package
	if err := json.Unmarshal([]byte(mockChapterData), &mangaChapters); err != nil {
		fmt.Printf("Error parsing mock chapter data: %v\n", err)
		return
	}

	mangaChapters.Process() // Assumes Process method exists on MangaChapters

	fmt.Printf("Found %d chapters:\n", len(mangaChapters.Chapters))
	for i, chapter := range mangaChapters.Chapters {
		fmt.Printf("[%d] %s (ID: %s)\n", i+1, chapter, mangaChapters.ChapterID[i])
	}

	if len(mangaChapters.Chapters) > 0 {
		fmt.Print("\nSelect a chapter (enter number, or 0 to exit): ")
		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil || choice <= 0 || choice > len(mangaChapters.Chapters) {
			fmt.Println("Invalid selection or exiting.")
			return
		}
		selectedIndex := choice - 1
		selectedID := mangaChapters.ChapterID[selectedIndex]
		fmt.Println("\n--- Mock Chapter Details ---")
		mockAtHomeServer(selectedID) // Call internal helper
	}
}

// mockAtHomeServer simulates fetching chapter image server info (internal helper)
func mockAtHomeServer(chapterID string) {
	fmt.Printf("MOCK MODE: Getting at-home server for chapter ID: %s\n", chapterID)

	mockAtHome := AtHomeResponse{ // Use type directly
		Result:  "ok",
		BaseURL: "https://uploads.mangadex.org",
		Chapter: struct {
			Hash      string   `json:"hash"`
			Data      []string `json:"data"`
			DataSaver []string `json:"dataSaver"`
		}{
			Hash:      "abcdef123456",
			Data:      []string{"01.png", "02.png", "03.png"},
			DataSaver: []string{"01.jpg", "02.jpg", "03.jpg"},
		},
	}

	fmt.Printf("Base URL: %s\n", mockAtHome.BaseURL)
	fmt.Printf("Chapter Hash: %s\n", mockAtHome.Chapter.Hash)
	fmt.Printf("Image Count: %d\n", len(mockAtHome.Chapter.Data))

	if len(mockAtHome.Chapter.Data) > 0 {
		fmt.Printf("First image (full quality): %s/data/%s/%s\n",
			mockAtHome.BaseURL, mockAtHome.Chapter.Hash, mockAtHome.Chapter.Data[0])
		fmt.Printf("First image (data saver): %s/data-saver/%s/%s\n",
			mockAtHome.BaseURL, mockAtHome.Chapter.Hash, mockAtHome.Chapter.DataSaver[0])
	}

	fmt.Print("\nSimulate download of mock images? (y/n): ")
	var answer string
	fmt.Scanln(&answer)
	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		simulateMockDownload(mockAtHome) // Call internal helper
	}
}

// simulateMockDownload simulates downloading images in mock mode (internal helper)
func simulateMockDownload(atHome AtHomeResponse) { // Use type directly
	outputDir := fmt.Sprintf("manga/mock_chapter_%s", atHome.Chapter.Hash)
	fmt.Printf("Mock Mode: Would download %d images to: %s\n", len(atHome.Chapter.Data), outputDir)

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	for i, filename := range atHome.Chapter.Data {
		imageURL := fmt.Sprintf("%s/data/%s/%s", atHome.BaseURL, atHome.Chapter.Hash, filename)
		outputPath := fmt.Sprintf("%s/%03d_%s", outputDir, i+1, filename)
		fmt.Printf("Mock downloading image %d/%d: %s\n", i+1, len(atHome.Chapter.Data), filename)
		out, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			continue
		}
		out.WriteString(fmt.Sprintf("Mock image data for %s\nURL: %s", filename, imageURL))
		out.Close()
		time.Sleep(200 * time.Millisecond)
		// Need a min function here or replace logic
		// fmt.Printf("\r[%-50s] %d%%", strings.Repeat("=", (i+1)*50/len(atHome.Chapter.Data)), (i+1)*100/len(atHome.Chapter.Data))
		fmt.Printf("\rProgress: %d / %d", i+1, len(atHome.Chapter.Data)) // Simplified progress
	}

	fmt.Println("\nMock download complete!")
	fmt.Printf("Mock images saved to: %s\n", outputDir)
}
