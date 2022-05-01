package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// CreateEmptyFile a function to create an empty file
func CreateEmptyFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// DeleteFile function to delete os file
func (app *App) DeleteFile(file string) {
	err := os.RemoveAll(file)
	if err != nil {
		fmt.Println(err)
	}
}

func (app *App) StringifyHtml(url string) string {

	// set user-agent
	client := &http.Client{}
	html, _ := http.Get(url)
	body := html.Body // get body
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_2 like Mac OS X) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0 Mobile/14C92 Safari/602.1")

	results, err := client.Do(req)

	bytes, err := io.ReadAll(results.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)

}

func (app *App) SaveHtml(url string, fileName string) {
	// get the html from the url

	// save it to a file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}
	// turn result into a string
	// write the string to the file
	_, err2 := file.WriteString(app.StringifyHtml(url))
	if err2 != nil {
		fmt.Println("Error writing to file", err2)
		return
	}
}

// LoadHtml load html from file to string
func (app *App) LoadHtml(file string) string {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file", err)
		return ""
	}
	// read the file
	bytes, err2 := io.ReadAll(f)
	if err2 != nil {
		fmt.Println("Error reading file", err2)
		return ""
	}
	// turn result into a string
	return string(bytes)
}

func (app *App) NewDir(dir string) {

	// if directory doesn't exist, create it
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		err := os.Mkdir("images", 0755)
		if err != nil {
			println("Error creating directory: " + err.Error())
		}
	} else if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
	}
}

// SaveImage save image to file
func (app *App) SaveImage(url string, filename string) {

	app.NewDir("images")
	//filename := app.GetImageNumber(url)

	fmt.Println("got page for filename:", filename)
	filename = strings.Replace(filename, ".png", ".jpg", -1)
	fmt.Println("tweaked filename", filename)
	filename = "images/" + filename
	fmt.Println("Image being written to file location: " + filename)

	results, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(results.Body)

	emptyFile, _ := CreateEmptyFile(filename)
	_, copyErr := io.Copy(emptyFile, results.Body)
	if copyErr != nil {
		fmt.Println("error copying file", copyErr)
	}
}
