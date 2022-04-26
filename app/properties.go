package app

type Manga struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Genre       string `json:"genre"`
	Chapter     int    `json:"chapter"`
	Pages       Pages
}

type Page struct {
	Number      int
	ImageUrl    string
	ImageKey    string
	Description string
}

type Pages []Page

type App struct {
	Query    string
	MangaAPI MangaAPI
	Mangas   []Manga
	// Viewer
}

type MangaAPI struct {
	HostName   string
	Api        string
	MangaId    int
	PageNumber int
}
