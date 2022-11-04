package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/spf13/afero"
)

var Afero = afero.Afero{Fs: afero.NewOsFs()}

// CreateEmptyFile a function to create an empty file
func CreateEmptyFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// DeleteFile function to delete os file
func DeleteFile(file string) {
	err := os.RemoveAll(file)
	if err != nil {
		fmt.Println(err)
	}
}

// SaveHtml save html to file
func SaveHtml(url string, fileName string) {
	// get the html from the url

	// save it to a file
	file := CreateFile(fileName)
	// turn result into a string
	// write the string to the file
	res, reqErr := CustomRequest(url)
	if reqErr != nil {
		fmt.Println("Error getting html from url", reqErr)
		return
	}
	_, fileErr := file.WriteString(res)

	if fileErr != nil {
		fmt.Println("Error writing to file", fileErr)
		return
	}
}

// CreateFile creates a file
func CreateFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file", err)
		return nil
	}
	return file
}

// LoadHtml load html from file to string
func LoadHtml(file string) string {
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

// NewDir creates a new directory
func NewDir(dir string) {
	// if directory doesn't exist, create it
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory: " + err.Error())
		}
	} else if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
	}
}

// ExitIfExists exits the program if the file already exists
func ExitIfExists(dir string) {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		fmt.Println("Manga-G detected that you already have this manga downloaded. Exiting...")
		os.Exit(0)
	}
}

// SaveImage save image to file
func SaveImage(url string, filename string) {

	// fmt.Println("got page for filename:", filename)
	// filename = strings.Replace(filename, ".png", ".jpg", -1)
	// fmt.Println("Image being written to file location: " + filename)

	results, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(results.Body)

	emptyFile, _ := CreateEmptyFile(filename)
	_, copyErr := io.Copy(emptyFile, results.Body)
	if copyErr != nil {
		fmt.Println("error copying file", copyErr)
	}
}

func TitleToDirName(title string) string {
	reg, _ := regexp.Compile("[^a-zA-Z\\d]+")
	return reg.ReplaceAllString(title, "")
}
func RemoveIfExists(path string) error {
	exists, err := Afero.Exists(path)

	if err != nil {
		return err
	}

	if exists {
		err = Afero.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFileNames returns a slice of strings containing all the file names in a directory
func GetFileNames(dir string) ([]string, error) {
	files, err := Afero.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}
