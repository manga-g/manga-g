package app

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

// ValidateUrl checks if url is valid
func (app *App) ValidateUrl(UrlToCheck string) bool {
	_, err := url.ParseRequestURI(UrlToCheck)
	if err != nil {
		return false
	}
	return true
}

// get title form http request

// StringToInt to change string to int
func (app *App) StringToInt(str string) int {
	var i int
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// FindImageUrl write a function to find image url from html
func (app *App) FindImageUrl(html string) ([]string, error) {
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

// FindMangaTitle from html string
func (app *App) FindMangaTitle(html string) string {
	TitleReg, _ := regexp.Compile("<title>(.*)</title>")
	Title := TitleReg.FindString(html)
	Title = strings.Replace(Title, "<title>", "", -1)
	Title = strings.Replace(Title, "</title>", "", -1)
	fmt.Println("Manga Title: " + Title)
	return Title
}
