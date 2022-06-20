package app

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"manga-g/ui"
)

// GetInput Function to get user input from the command line
func (app *MG) GetInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return ""
	}
	return input
}

// ValidateUrl checks if url is valid
func (app *MG) ValidateUrl(UrlToCheck string) bool {
	_, err := url.ParseRequestURI(UrlToCheck)
	if err != nil {
		return false
	}
	return true
}

// StringToInt to change string to int
func (app *MG) StringToInt(str string) int {
	var i int
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return 0
	}
	return i
}

// FindImageUrl write a function to find image url from html
func (app *MG) FindImageUrl(html string) ([]string, error) {
	reg, _ := regexp.Compile("(https?://[a-z]\\d+.[a-z]+.[a-z]+/[a-z]+/\\d+/\\d+.(png|jpg|gif))")
	match := reg.FindStringSubmatch(html)
	if match != nil {
		return match, nil
	}
	return nil, errors.New("could not find image url")
}

// FindImageKey from html and url
func (app *MG) FindImageKey(url string) string {
	reg, _ := regexp.Compile("galleries/(\\d+)/\\d+.(jpg|png|gif)")
	// try to find the image key from the url
	ImageKey := reg.FindStringSubmatch(url)
	return ImageKey[1]
}

// GetImageNumber from url
func (app *MG) GetImageNumber(url string) string {
	reg, _ := regexp.Compile("\\d+.(jpg|png|gif)")
	return reg.FindString(url)
}

// FindMangaTitle from html string
func (app *MG) FindMangaTitle(html string) string {
	TitleReg, _ := regexp.Compile("<title>(.*)</title>")
	Title := TitleReg.FindString(html)
	Title = strings.Replace(Title, "<title>", "", -1)
	Title = strings.Replace(Title, "</title>", "", -1)
	fmt.Println("Manga Title: " + Title)
	return Title
}

// InitMangaG initializes the MangaG struct.
func InitMangaG() {
	MangaG := new(MG)
	style := ui.InitStyle()

	//ui.Render(style, "Starting Manga-G...")
	fmt.Println("Starting Manga-G...")
	if MangaG.Connected() {
		//ui.Render(style, "Enter a URL for a Manga's first page to download:")
		fmt.Println("Enter a URL for a Manga's first page to download:")
		MangaUrl := MangaG.GetInput()
		if MangaG.ValidateUrl(MangaUrl) {
			DoStuff(MangaG, MangaUrl)
		} else {
			//ui.Render(style, "Invalid URL please try again.")
			//ui.Render(style, "Example: https://www.mangaeden.com/en/en-manga/one-piece/")
			//ui.Render(style, "Exiting...")
			fmt.Println("Invalid URL please try again.")
			fmt.Println("Example: https://www.mangaeden.com/en/en-manga/one-piece/")
			fmt.Println("Exiting...")

		}
	} else {
		//ui.Render(style, "Could not connect to the internet.")
		fmt.Println("Could not connect to the internet.")
		//ui.Render(style, "Exiting...")
		fmt.Println("Exiting...")
	}
}
