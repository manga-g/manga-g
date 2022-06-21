package app

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "regexp"
    "strings"
)

// CreateEmptyFile a function to create an empty file
func (app *MG) CreateEmptyFile(fileName string) (*os.File, error) {
    file, err := os.Create(fileName)
    if err != nil {
        return nil, err
    }
    return file, nil
}

// DeleteFile function to delete os file
func (app *MG) DeleteFile(file string) {
    err := os.RemoveAll(file)
    if err != nil {
        fmt.Println(err)
    }
}

func (app *MG) SaveHtml(url string, fileName string) {
    // get the html from the url

    // save it to a file
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Error creating file", err)
        return
    }
    // turn result into a string
    // write the string to the file
    res, _ := app.CustomRequest(url)
    _, err2 := file.WriteString(res)
    if err2 != nil {
        fmt.Println("Error writing to file", err2)
        return
    }
}

// LoadHtml load html from file to string
func (app *MG) LoadHtml(file string) string {
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

func (app *MG) NewDir(dir string) {

    // if directory doesn't exist, create it
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        err := os.Mkdir(dir, 0755)
        if err != nil {
            println("Error creating directory: " + err.Error())
        }
    } else if err != nil {
        fmt.Printf("Error creating directory: %s\n", err)
    }
}

// SaveImage save image to file
func (app *MG) SaveImage(url string, filename string) {

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

    emptyFile, _ := app.CreateEmptyFile(filename)
    _, copyErr := io.Copy(emptyFile, results.Body)
    if copyErr != nil {
        fmt.Println("error copying file", copyErr)
    }
}

func (app *MG) TitleToDirName(title string) string {
    reg, _ := regexp.Compile("[^a-zA-Z\\d]+")
    return reg.ReplaceAllString(title, "")
}
