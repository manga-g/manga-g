package app

type MangaList []struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Author             string `json:"author"`
	LatestChapter      string `json:"latestChapter"`
	LastUpdateDatetime string `json:"lastUpdateDatetime"`
	ImagePreview       string `json:"imagePreview"`
	ViewCount          string `json:"viewCount"`
}
type MangaChapters struct {
	Status             string   `json:"status"`
	LastUpdateDatetime string   `json:"lastUpdateDatetime"`
	ViewCount          string   `json:"viewCount"`
	Author             string   `json:"author"`
	Chapters           []string `json:"chapters"`
	LinkInfo           []string `json:"linkInfo"`
	DateInfo           []string `json:"dateInfo"`
	ChapterID          []string `json:"chapterId"`
	Description        string   `json:"description"`
}

type MangaImages []struct {
	ImageUrl string `json:"$img_link"`
	Alt      string `json:"$img_alt"`
}
