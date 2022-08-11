package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/manga-g/manga-g/app"
)

// Entrypoint for the program.
func main() {
	//gotenv.Load()
	//basedApiUrl := os.Getenv("BASED_API_URL")
	//fmt.Println("BASED_API_URL:", basedApiUrl)
	//if basedApiUrl == "" {
	//    fmt.Println("BASED_API_URL is not set in env")
	//    os.Exit(1)
	//}

	//mangaSaveDir := os.Getenv("MANGA_SAVE_DIR")
	//fmt.Println("MANGA_SAVE_DIR:", mangaSaveDir)
	//if mangaSaveDir == "" {
	//    fmt.Println("MANGA_SAVE_DIR is not set in env")
	//    currentDirectory, err := os.Getwd()
	//    if err != nil {
	//        fmt.Println("Error getting current directory:", err)
	//        os.Exit(1)
	//    }
	//    fmt.Println("Using default" + currentDirectory)
	//    mangaSaveDir = "."
	//}

	//port := "3000"
	//basedApiUrl := "http://localhost:" + port + "/"
	basedApiUrl := "http://manga-api.bytecats.codes/"
	mangaSaveDir := "./"
	_, err := http.Get(basedApiUrl)
	if err == nil {
		fmt.Println("Online Manga is ready!")
	} else {
		fmt.Println("Online Manga is offline at the moment ;(\nBe back online ASAP =)")
		os.Exit(1)
	}
	fmt.Print("Search for manga: ")
	query := app.GetInput()
	//if there is no input, loop the request 3 times
	if query == "" {
		for n := 0; n < 3; n++ {
			fmt.Println("You should choose the number corresponding to the manga you want to read.\nTry again,please :)\n" + "Search for manga: ")
			query = app.GetInput()
			if query != "" {
				break
			}
		}
	}
	query = url.QueryEscape(query)
	fmt.Println("Searching for:", query)
	apiSearch := basedApiUrl + "mk/search?q=" + query
	res, _ := app.CustomRequest(apiSearch)
	// wait for the results to be ready
	var mangaList app.MangaList
	app.ParseMangaSearch(res, &mangaList)
	fmt.Println("Found", len(mangaList), "manga")
	var titles []string
	for i, manga := range mangaList {
		titles = append(titles, fmt.Sprintf("%d. %s", i+1, manga.Title))
	}

	number := len(titles)

	for _, title := range titles {
		fmt.Println(title)
	}

	SelectMessage := "Select a title: (1 - " + strconv.Itoa(number) + ") "
	fmt.Print(SelectMessage)
	mangaChoice := app.GetInput()
	//if there is no input, loop the request 3 times
	if mangaChoice == "" {
		for n := 0; n < 3; n++ {
			fmt.Println("You should choose the number corresponding to the manga you want to read.\nTry again,please :)" + SelectMessage)
			mangaChoice = app.GetInput()
			if mangaChoice != "" {
				break
			}
		}
	}
	mangaChoiceInt := app.StringToInt(mangaChoice)
	if mangaChoiceInt > number {
		fmt.Println("Invalid choice")
		return
	}
	mangaId := mangaList[mangaChoiceInt-1].ID
	chapterUrl := basedApiUrl + "mk/chapters?q=" + mangaId
	fmt.Println("Checking ID:" + mangaId)
	fmt.Println("Loading chapters...")
	res, _ = app.CustomRequest(chapterUrl)
	var mangaChapters app.MangaChapters
	app.ParseChapters(res, &mangaChapters)
	var chapterTitles []string
	n := 0
	for i := len(mangaChapters.Chapters); i >= 1; i-- {
		chapter := mangaChapters.Chapters[i-1]
		chapterTitles = append(chapterTitles, fmt.Sprintf("%d. %s", n+1, chapter))
		n++
	}

	for _, title := range chapterTitles {
		fmt.Println(title)
	}

	fmt.Print("Select a chapter: (1 - " + strconv.Itoa(len(chapterTitles)) + ") ")
	chapterChoice := app.GetInput()
	//if there is no input, loop the request 3 times
	if chapterChoice == "" {
		for n := 0; n < 3; n++ {
			fmt.Println("You should choose the number corresponding to the chapter you want to read.\nTry again,please :)" + "Select a chapter: (1 - " + strconv.Itoa(len(chapterTitles)) + ") ")
			chapterChoice = app.GetInput()
			if chapterChoice != "" {
				break
			}
		}
	}
	chapterChoiceInt := app.StringToInt(chapterChoice)
	if chapterChoiceInt > len(chapterTitles) {
		fmt.Println("Invalid choice")
		return
	}
	chapterId := mangaChapters.ChapterID[chapterChoiceInt-1]
	chapterNumber := strings.Replace(chapterId, "chapter-", "", -1)

	fmt.Println("Chapter number:", chapterNumber)

	fmt.Println("Trying to load images for " + chapterNumber)
	// keep only the number at the end of the string

	imagesUrl := basedApiUrl + "mk/images?id=" + mangaId + "&chapterNumber=" + chapterNumber
	fmt.Println("Loading images...")
	fmt.Println(imagesUrl)
	res, _ = app.CustomRequest(imagesUrl)
	var mangaImages app.MangaImages
	app.ParseImages(res, &mangaImages)
	var images []string
	for _, image := range mangaImages {
		images = append(images, image.ImageUrl)
	}

	app.NewDir(mangaSaveDir + "/" + "manga")

	mangaName := strings.Replace(mangaList[mangaChoiceInt-1].Title, " ", "_", -1)
	mangaName = strings.Replace(mangaName, ":", "", -1)
	mangaName = strings.Replace(mangaName, " ", "_", -1)

	app.NewDir(mangaSaveDir + "/" + "manga/" + mangaName)
	app.NewDir(mangaSaveDir + "/" + "manga/" + mangaName + "/" + chapterNumber)

	fmt.Println("Downloading", len(images), "pages")
	//fmt.Println(images)
	for _, image := range images {
		imageName := strings.Split(image, "/")
		imageName = strings.Split(imageName[len(imageName)-1], ".")
		imageName[0] = strings.Replace(imageName[0], " ", "_", -1)
		imageFullDir := mangaSaveDir + "manga/" + mangaName + "/" + chapterNumber + "/" + imageName[0] + "." + imageName[1]
		app.SaveImage(image, imageFullDir)
	}
	//pdfDir := mangaSaveDir + "manga/" + mangaName + "/" + chapterNumber + "/"

	fmt.Println("Done")
}
