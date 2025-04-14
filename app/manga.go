package app

import (
	"fmt"
)

// Manga Type to hold all the information about a manga for api
type Manga struct {
	SelectedMangaID string
	SearchList      *MangaList
	ChapterList     *MangaChapters
	ImageList       *MangaImages
}

// MangaDex API v5 response structure for manga search
type MangaDexResponse struct {
	Result   string      `json:"result"`
	Response string      `json:"response"`
	Data     interface{} `json:"data"`
	Limit    int         `json:"limit"`
	Offset   int         `json:"offset"`
	Total    int         `json:"total"`
}

// MangaList holds list of manga for search results
// Updated to match MangaDex API v5 format
type MangaList []struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Title       map[string]string   `json:"title"`
		AltTitles   []map[string]string `json:"altTitles"`
		Description map[string]string   `json:"description"`
		Status      string              `json:"status"`
		Year        int                 `json:"year"`
		Tags        []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Name map[string]string `json:"name"`
			} `json:"attributes"`
		} `json:"tags"`
	} `json:"attributes"`
	Relationships []struct {
		ID         string      `json:"id"`
		Type       string      `json:"type"`
		Attributes interface{} `json:"attributes,omitempty"`
	} `json:"relationships"`
}

// Helper method to get the title in a readable format
func (m MangaList) GetTitle(index int) string {
	// First check if there's an English title
	if title, ok := m[index].Attributes.Title["en"]; ok {
		return title
	}

	// Then check for any title
	for _, title := range m[index].Attributes.Title {
		return title
	}

	return "Unknown Title"
}

// MangaChapters holds data about manga chapters
// Updated to match MangaDex API v5 format
type MangaChapters struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Volume             string `json:"volume"`
			Chapter            string `json:"chapter"`
			Title              string `json:"title"`
			TranslatedLanguage string `json:"translatedLanguage"`
			ExternalURL        string `json:"externalUrl"`
			PublishAt          string `json:"publishAt"`
			ReadableAt         string `json:"readableAt"`
			CreatedAt          string `json:"createdAt"`
			UpdatedAt          string `json:"updatedAt"`
		} `json:"attributes"`
	} `json:"data"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`

	// Compatibility with old code
	Status             string   `json:"-"`
	LastUpdateDatetime string   `json:"-"`
	ViewCount          string   `json:"-"`
	Author             string   `json:"-"`
	Chapters           []string `json:"-"`
	LinkInfo           []string `json:"-"`
	DateInfo           []string `json:"-"`
	ChapterID          []string `json:"-"`
	Description        string   `json:"-"`
}

// Process prepares the compatibility fields from the API response
func (mc *MangaChapters) Process() {
	// Fill in compatibility fields
	mc.ChapterID = make([]string, len(mc.Data))
	mc.Chapters = make([]string, len(mc.Data))

	for i, chapter := range mc.Data {
		mc.ChapterID[i] = chapter.ID
		title := chapter.Attributes.Title
		if title == "" {
			title = "Chapter " + chapter.Attributes.Chapter
		}
		mc.Chapters[i] = title
	}
}

// AtHomeResponse holds the response from the at-home server
type AtHomeResponse struct {
	Result  string `json:"result"`
	BaseURL string `json:"baseUrl"`
	Chapter struct {
		Hash      string   `json:"hash"`
		Data      []string `json:"data"`      // Original quality
		DataSaver []string `json:"dataSaver"` // Compressed quality
	} `json:"chapter"`
}

// MangaImages stores list of manga images for downloading
type MangaImages []struct {
	ImageUrl string `json:"$img_link"`
	Alt      string `json:"$img_alt"`
}

// ProcessAtHomeResponse generates image URLs from the at-home response
func ProcessAtHomeResponse(atHome AtHomeResponse) MangaImages {
	// Create image URLs from at-home data
	var images MangaImages

	// Use the data-saver (compressed) images for smaller downloads
	for _, filename := range atHome.Chapter.DataSaver {
		imgURL := fmt.Sprintf("%s/data-saver/%s/%s",
			atHome.BaseURL, atHome.Chapter.Hash, filename)

		images = append(images, struct {
			ImageUrl string `json:"$img_link"`
			Alt      string `json:"$img_alt"`
		}{
			ImageUrl: imgURL,
			Alt:      filename,
		})
	}

	return images
}
