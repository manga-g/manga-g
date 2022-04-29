package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// parse jpg url and hash
// write a regex to extract the hash and the url

func (app *App) FindImageUrls(html string) []string {
	var urls []string
	rx := regexp.MustCompile("https://urasunday.com/secure/\\d+/webp/manga_page_high/\\d+/\\d+.webp\\?hash=.*?&expires=\\d+")
	for _, match := range rx.FindAllStringSubmatch(html, -1) {
		match[0] = strings.ReplaceAll(match[0], "'", "")
		match[0] = strings.ReplaceAll(match[0], "\"", "")
		match[0] = strings.ReplaceAll(match[0], "\\", "")
		match[0] = strings.ReplaceAll(match[0], " ", "")
		match[0] = strings.ReplaceAll(match[0], "\n", "")
		match[0] = strings.ReplaceAll(match[0], "\t", "")
		match[0] = strings.ReplaceAll(match[0], "'", "")
		match[0] = strings.ReplaceAll(match[0], "alt=", "")
		match[0] = strings.ReplaceAll(match[0], ">", "")
		//match[0] = match[0][:len(match[0])-2]

		urls = append(urls, match[0])
	}
	return urls
}

// function to download the webp from all urls in the array
// write a function to download the webp from the url
// write a function to write the webp to the disk
func (app *App) DownloadWebp(url string, filename string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	results, err2 := client.Do(req)
	if err2 != nil {
		log.Fatal(err2)
	}

	emptyFile, _ := CreateEmptyFile(filename)
	_, copyErr := io.Copy(emptyFile, results.Body)
	if copyErr != nil {
		fmt.Println("error copying file", copyErr)
	}
}

// get title from url header and use as filename but remove spaces and special characters and use smart numbers to avoid overwriting
func GetWebpTitle(title string) string {
	title = strings.ReplaceAll(title, " ", "")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, "?", "")
	title = strings.ReplaceAll(title, "!", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "'", "")
	title = strings.ReplaceAll(title, "\\", "")
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\t", "")
	title = strings.ReplaceAll(title, "`", "")
	title = strings.ReplaceAll(title, "~", "")

	return title
}

// function to make a for loop to download all webp files from the array of urls from FindImageUrls
func (app *App) DownloadAllWebp(urls []string) {
	for i, url := range urls {
		fmt.Println(i, url)
		number := fmt.Sprintf("%03d", i)
		filename := "images/0" + number + ".webp"
		app.DownloadWebp(url, filename)
	}
}
