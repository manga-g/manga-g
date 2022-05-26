package extensions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"manga-g/app"
)

// IncrementImageUrl to increment the image url
func IncrementImageUrl(app *app.MG, url string) string {
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

func GetPageCount(app *app.MG, html string) int {
	reg, _ := regexp.Compile("<span class=\"num-pages\">[1-9]*</span>")
	PageCount := reg.FindString(html)
	PageCount = strings.Replace(PageCount, "<span class=\"num-pages\">", "", -1)
	PageCount = strings.Replace(PageCount, "</span>", "", -1)
	PageCountInt := app.StringToInt(PageCount)
	fmt.Println("Page Count: ", PageCountInt)
	return PageCountInt
}

func CycleImages(app *app.MG, ImageUrl []string, max int) {
	wg := new(sync.WaitGroup)
	for i := 1; i < max; i++ {
		wg.Add(1)
		go func(ImageUrl []string, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Println("Attempting to download page:", i)
			app.SaveImage(ImageUrl[1], app.GetImageNumber(ImageUrl[1]))
			ImageUrl[1] = IncrementImageUrl(app, ImageUrl[1])
		}(ImageUrl, i, wg)
		wg.Wait()
	}
	fmt.Println("Finished downloading all pages.")
}
