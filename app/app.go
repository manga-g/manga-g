package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// GetInput Function to get user input from the command line
func (app *App) GetInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return ""
	}
	return input
}

//// GetManga function to get the manga information form a url string
//func (app *App) GetManga(url string) Manga {
//	// get the manga information from the url
//	// and store it in the app.Mangas
//	var manga Manga
//
//	return manga
//}

// AddManga adds a new to append add manga to app.Mangas
func (app *App) AddManga(manga Manga) {
	// append manga to app.Mangas
	app.Mangas = append(app.Mangas, manga)
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

// StringToInt to change string to int
func (app *App) StringToInt(str string) int {
	var i int
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// GetMangaPages function to get the manga pages form a url string
//func (app *App) GetMangaPages(api MangaAPI) {
//	// get the page number from the url
//	// get the image file location from the body of the request
//	// and store it in the app.Mangas
//	// get the api endpoint from the url
//	// get the manga id from the url
//	// get the page number from the url if any
//	// add them to the manga api struct
//}

// FindImageUrl write a function to find image url from html
func (app *App) FindImageUrl(html string) ([]string, error) {

	// find the image url from the html
	//bytes := []byte(html)

	// find the image url from the html
	// find urls from string using regex
	// find all urls ending in an image extension
	// return the first url

	reg, _ := regexp.Compile("(https?://[a-z]\\d+.[a-z]+.[a-z]+/[a-z]+/\\d+/\\d+.(png|jpg|gif))")
	match := reg.FindStringSubmatch(html)
	if match != nil {
		return match, nil
	}
	return nil, errors.New("no image url found")
}

// FindImageKey from html and url
func (app *App) FindImageKey(url string) string {
	reg, _ := regexp.Compile("galleries/(\\d+)/\\d+.(jpg|png|gif)")
	// try to find the image key from the url
	ImageKey := reg.FindStringSubmatch(url)
	return ImageKey[1]
}

// GetImageNumber from url
func (app *App) GetImageNumber(url string) string {
	reg, _ := regexp.Compile("\\d+.(jpg|png|gif)")
	return reg.FindString(url)
}

// FindMangaId from html string
func (app *App) FindMangaId(html string) string {
	reg, _ := regexp.Compile("/g/[1-9]*/")
	MangaId := strings.Replace(reg.FindString(html), "/g/", "", -1)
	MangaId = strings.Replace(MangaId, "/", "", -1)
	fmt.Println("Manga Id: " + MangaId)
	return MangaId
}

// FindMangaTitle from html string
func (app *App) FindMangaTitle(html string) string {
	TitleReg, _ := regexp.Compile("<title>(.*)</title>")
	Title := TitleReg.FindString(html)
	Title = strings.Replace(Title, "<title>", "", -1)
	Title = strings.Replace(Title, "</title>", "", -1)
	fmt.Println("Manga Title: " + Title)
	return Title
}

// IncrementImageUrl to increment the image url
func (app *App) IncrementImageUrl(url string) string {
	ImageName := app.GetImageNumber(url)
	ImageNumber := strings.Replace(ImageName, ".png", "", -1)
	ImageNumber = strings.Replace(ImageNumber, ".jpg", "", -1)
	ImageNumberInt := app.StringToInt(ImageNumber)
	ImageNumberInt++

	fmt.Println("ImageNumber incremented: ", ImageNumber)
	ImageNumber = strconv.Itoa(ImageNumberInt)

	//ImageNumber = strings.Repeat("0", 4-len(ImageNumber)) + ImageNumber
	ImageNumber = ImageNumber + ".jpg"
	url = strings.Replace(url, ImageName, ImageNumber, -1)
	final := strings.Replace(url, ImageNumber, "", -1) + ImageNumber
	fmt.Print("Final url: ", final, "\n")
	return final
}

func (app *App) GetPageCount(html string) int {
	reg, _ := regexp.Compile("<span class=\"num-pages\">[1-9]*</span>")
	PageCount := reg.FindString(html)
	PageCount = strings.Replace(PageCount, "<span class=\"num-pages\">", "", -1)
	PageCount = strings.Replace(PageCount, "</span>", "", -1)
	PageCountInt := app.StringToInt(PageCount)
	fmt.Println("Page Count: ", PageCountInt)
	return PageCountInt
}

func (app *App) CycleImages(ImageUrl []string, max int) {
	wg := new(sync.WaitGroup)
	for i := 1; i < max; i++ {
		wg.Add(1)
		go func(ImageUrl []string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Println("Attempting to download page:", i)
			app.SaveImage(ImageUrl[1], app.GetImageNumber(ImageUrl[1]))
			ImageUrl[1] = app.IncrementImageUrl(ImageUrl[1])
		}(ImageUrl, i, wg)
		wg.Wait()
	}
	fmt.Println("Finished downloading all pages.")
}
