package app

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// GetMangaAPI function to get the manga information form a url string
func (app *MG) GetMangaAPI(url string) MangaAPI {
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

// FindMangaId from html string
func (app *MG) FindMangaId(html string) string {
	reg, _ := regexp.Compile("/g/[1-9]*/")
	MangaId := strings.Replace(reg.FindString(html), "/g/", "", -1)
	MangaId = strings.Replace(MangaId, "/", "", -1)
	fmt.Println("Manga Id: " + MangaId)
	return MangaId
}

// AddManga adds a new to append add manga to app.Mangas
func (app *MG) AddManga(manga Manga) {
	// append manga to app.Mangas
	app.Mangas = append(app.Mangas, manga)
}

//// GetManga function to get the manga information form a url string
//func (app *App) GetManga(url string) Manga {
//	// get the manga information from the url
//	// and store it in the app.Mangas
//	var manga Manga
//
//	return manga
//}
