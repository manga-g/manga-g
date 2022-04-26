package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	results, _ := http.Get(url)
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

func NewDir(dir string) {

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
