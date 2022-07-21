package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manga-g/manga-g/app"
)

// Entrypoint for the program.
func main() {
	basedApiUrl := os.Getenv("BASED_API_URL")
	if basedApiUrl == "" {
		fmt.Println("BASED_API_URL is not set in env")
		os.Exit(1)
	}

	mangaSaveDir := os.Getenv("MANGA_SAVE_DIR")
	if mangaSaveDir == "" {
		fmt.Println("MANGA_SAVE_DIR is not set in env")
		currentDirectory, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
		fmt.Println("Using default" + currentDirectory)
		mangaSaveDir = "."
	}
	fmt.Print("Search for manga: ")
	query := app.GetInput()
	apiSearch := basedApiUrl + "mk/search?q=" + query
	res, _ := app.CustomRequest(apiSearch)
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
	for i, chapter := range mangaChapters.Data {
		chapterTitles = append(chapterTitles, fmt.Sprintf("%d. %s", i+1, chapter))
	}
	for _, title := range chapterTitles {
		fmt.Println(title)
	}

	fmt.Print("Select a chapter: (1 - " + strconv.Itoa(len(chapterTitles)) + ") ")
	chapterChoice := app.GetInput()
	chapterChoiceInt := app.StringToInt(chapterChoice)
	if chapterChoiceInt > len(chapterTitles) {
		fmt.Println("Invalid choice")
		return
	}
	chapterNumberChoice := mangaChapters.Data[chapterChoiceInt-1]
	// use replace to get rid of any non-numeric characters in the ChapterNumberChoice string
	chapterNumberChoice = strings.Replace(chapterNumberChoice, "Chapter ", "", -1)

	fmt.Println("Trying to load images for chapter: " + chapterNumberChoice)

	imagesUrl := basedApiUrl + "mk/images?id=" + mangaId + "&chapterNumber=" + chapterNumberChoice
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
	app.NewDir(mangaSaveDir + "/" + "manga/" + mangaName)
	app.NewDir(mangaSaveDir + "/" + "manga/" + mangaName + "/" + chapterNumberChoice)

	fmt.Println("Downloading", len(images), "pages")
	fmt.Println(images)
	for _, image := range images {
		imageName := strings.Split(image, "/")
		imageName = strings.Split(imageName[len(imageName)-1], ".")
		app.SaveImage(image, "manga/"+mangaName+"/"+chapterNumberChoice+"/"+imageName[0]+"."+imageName[1])

	}
	fmt.Println("Done")
}
