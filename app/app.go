package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

// GetInput Function to get user input from the command line
func GetInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return ""
	}
	return input
}

// access html and get manga information

// GetManga function to get the manga information form a url string
func (app *App) GetManga(url string) Manga {
	// get the manga information from the url
	// and store it in the app.Mangas
	var manga Manga

	return manga
}

// AddManga adds a new to append add manga to app.Mangas
func (app *App) AddManga(manga Manga) {
	// append manga to app.Mangas
	app.Mangas = append(app.Mangas, manga)
}

type MangaAPI struct {
	HostName   string
	Api        string
	MangaId    int
	PageNumber int
}

// get title form http request

// GetMangaAPI function to get the manga information form a url string
func (app *App) GetMangaAPI(url string) MangaAPI {
	var api MangaAPI
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(result.Body)

	// get the manga information from the url
	api.HostName = result.Request.URL.Host
	api.Api = result.Request.URL.Path

	// get the integer after the /g/ in the url
	api.MangaId = app.StringToInt(result.Request.URL.Query()["g"][0])

	// get the  page integer from the hostname/g/{manga_id}/{page} in the url
	// find image file location from body of the request

	// and store it in the app.Mangas

	// get the api endpoint from the url
	// get the manga id from the url
	// get the page number from the url if any
	// add them to the manga api struct

	return api
}

// function to change string to int
func (app *App) StringToInt(str string) int {
	var i int
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// GetMangaPages function to get the manga pages form a url string
func (app *App) GetMangaPages(api MangaAPI) {
	// get the page number from the url
	// get the image file location from the body of the request
	// and store it in the app.Mangas
	// get the api endpoint from the url
	// get the manga id from the url
	// get the page number from the url if any
	// add them to the manga api struct
}

// SaveHtml write a function to save html from url to file
func (app *App) SaveHtml(url string) {
	// get the html from the url
	results, _ := http.Get(url)
	bytes, err := io.ReadAll(results.Body)
	if err != nil {
		fmt.Println(err)
	}

	// save it to a file
	file, err2 := os.Create("manga.html")
	if err2 != nil {
		return
	}

	// turn result into a string

	// write the string to the file
	_, err3 := file.WriteString(string(bytes))
	if err3 != nil {
		return
	}

}

// LoadHtml load html from file to string
func (app *App) LoadHtml(file string) string {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	// read the file
	bytes, _ := io.ReadAll(f)
	// turn result into a string
	return string(bytes)
}

// FindImageUrl write a function to find image url from html
func (app *App) FindImageUrl(html string) string {
	// find the image url from the html
	//bytes := []byte(html)

	// find the image url from the html
	// find urls from string using regex
	// find all urls ending in an image extension
	// return the first url
	reg, _ := regexp.Compile("https://[a-z]\\d.[a-z]*.[a-z]*/galleries/[1-9]*/[1-9]*.(jpg|png|gif)")
	return reg.FindString(html)
}

// Find imagekey from html and url
func (app *App) FindImageKey(url string) string {
	reg, _ := regexp.Compile("galleries/(\\d+)/\\d+.(jpg|png|gif)")
	ImageKey := reg.FindStringSubmatch(url)[1]
	return ImageKey
}

// Get image number from url
func (app *App) GetImageNumber(url string) string {
	reg, _ := regexp.Compile("\\d+.(jpg|png|gif)")
	return reg.FindString(url)
}
func (app *App) SaveImage(url string) {

	filename := app.GetImageNumber(url)
	filename = strings.Replace(filename, ".png", ".jpg", -1)
	filename = "images/" + filename
	fmt.Println("Image being written to file location: " + filename)
	// if directory doesn't exist, create it
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		err := os.Mkdir("images", 0777)
		if err != nil {
			println("Error creating directory: " + err.Error())
		}
	} else if err != nil {
		fmt.Println("Error creating directory: %s", err)
	}

	f, err2 := os.Create(filename)
	if err2 != nil {
		fmt.Println("Error creating image file", err2)
	}
	defer func(Body io.ReadCloser) {
		err3 := Body.Close()
		if err3 != nil {
			fmt.Println("Error Closing http body", err3)
		}
	}(f)
	results, _ := http.Get(url)
	bytes, err4 := io.ReadAll(results.Body)
	if err4 != nil {
		fmt.Println(err4)
	}
	_, err5 := f.Write(bytes)
	if err5 != nil {
		fmt.Println("Error writing to file", err5)
	} else {
		fmt.Println("Image saved")
	}
}
func (app *App) FindMangaId(html string) string {
	reg, _ := regexp.Compile("/g/[1-9]*/")
	MangaId := strings.Replace(reg.FindString(html), "/g/", "", -1)
	MangaId = strings.Replace(MangaId, "/", "", -1)
	return MangaId
}
func (app *App) FindMangaTitle(html string) string {
	TitleReg, _ := regexp.Compile("<title>(.*)</title>")
	Title := TitleReg.FindString(html)
	Title = strings.Replace(Title, "<title>", "", -1)
	Title = strings.Replace(Title, "</title>", "", -1)
	return Title
}
func (app *App) IncrementImageUrl(url string) string {
	ImageName := app.GetImageNumber(url)
	ImageNumber := strings.Replace(ImageName, ".png", "", -1)
	ImageNumber = strings.Replace(ImageNumber, ".jpg", "", -1)
	ImageNumberInt := app.StringToInt(ImageNumber)
	ImageNumberInt++

	fmt.Println("ImageNumber incremented: ", ImageNumber)

	ImageNumber = strconv.Itoa(ImageNumberInt)

	//ImageNumber = strings.Repeat("0", 4-len(ImageNumber)) + ImageNumber
	ImageNumber = ImageNumber + ".png"
	url = strings.Replace(url, ImageName, ImageNumber, -1)
	return strings.Replace(url, ImageNumber, "", -1) + ImageNumber
}

// DownloadAllPages function to download all images from a gallery url and save them to a folder same name as manga title
func (app *App) DownloadAllPages(galleryurl string) int {
	for i := 1; i <= app.GetPageCount(galleryurl); i++ {
		app.SaveImage(app.IncrementImageUrl(galleryurl))
		galleryurl = app.IncrementImageUrl(galleryurl)
		fmt.Println("Downloading page: ", i)
	}
	return 0
}

func (app *App) GetPageCount(html string) int {
	reg, _ := regexp.Compile("<span class=\"num-pages\">[1-9]*</span>")
	PageCount := reg.FindString(html)
	PageCount = strings.Replace(PageCount, "<span class=\"num-pages\">", "", -1)
	PageCount = strings.Replace(PageCount, "</span>", "", -1)
	PageCountInt := app.StringToInt(PageCount)
	return PageCountInt
}
