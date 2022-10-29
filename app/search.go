package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func MkSearch(basedApiUrl string, query string) {
	var mangaList MangaList
	var titles []string
	var chapterTitles []string
	var mangaChapters MangaChapters
	var mangaImages MangaImages
	var images []string

	mangaSaveDir := "./"
	apiSearch := basedApiUrl + "mk/search?q=" + query
	results, _ := CustomRequest(apiSearch)

	fmt.Println("Searching for:", query)
	ParseMangaSearch(results, &mangaList)
	fmt.Println("Found", len(mangaList), "manga")
	for i, manga := range mangaList {
		titles = append(titles, fmt.Sprintf("%d. %s", i+1, manga.Title))
	}
	number := len(titles)
	for _, title := range titles {
		fmt.Println(title)
	}

	SelectMessage := "Select a title: (1 - " + strconv.Itoa(number) + ") "
	fmt.Print(SelectMessage)
	mangaChoice := GetInput()
	QueryCheck(mangaChoice)

	mangaChoiceInt := StringToInt(mangaChoice)
	if mangaChoiceInt > number {
		fmt.Println("Invalid choice")
		return
	}
	mangaId := mangaList[mangaChoiceInt-1].ID
	chapterUrl := basedApiUrl + "mk/chapters?q=" + mangaId
	fmt.Println("Checking ID:" + mangaId)
	fmt.Println("Loading chapters...")

	results, _ = CustomRequest(chapterUrl)
	ParseChapters(results, &mangaChapters)

	n := 0
	for i := len(mangaChapters.Chapters); i >= 1; i-- {
		chapter := mangaChapters.Chapters[i-1]
		chapterTitles = append(chapterTitles, fmt.Sprintf("%d. %s", n+1, chapter))
		n++
	}

	for _, title := range chapterTitles {
		fmt.Println(title)
	}

	fmt.Print("Select a result: (1 - " + strconv.Itoa(len(chapterTitles)) + ") ")
	resultChoice := GetInput()
	QueryCheck(resultChoice)

	chapterChoiceInt := StringToInt(resultChoice)
	chapterChoiceInt = len(mangaChapters.ChapterID) - chapterChoiceInt
	if chapterChoiceInt > len(chapterTitles) || chapterChoiceInt < 1 {
		fmt.Println("Invalid choice")
		return
	}

	chapterId := mangaChapters.ChapterID[chapterChoiceInt]
	chapterNumber := strings.Replace(chapterId, "chapter-", "", -1)
	fmt.Println("Chapter number:", chapterNumber)

	// keep only the number at the end of the string
	imagesUrl := basedApiUrl + "mk/images?id=" + mangaId + "&chapterNumber=" + chapterNumber
	//	fmt.Println(imagesUrl)
	results, _ = CustomRequest(imagesUrl)
	ParseImages(results, &mangaImages)
	for _, image := range mangaImages {
		images = append(images, image.ImageUrl)
	}

	NewDir(mangaSaveDir + "/" + "manga")

	mangaName := strings.Replace(mangaList[mangaChoiceInt-1].Title, " ", "_", -1)
	mangaName = strings.Replace(mangaName, ":", "", -1)
	mangaName = strings.Replace(mangaName, " ", "_", -1)

	NewDir(mangaSaveDir + "/" + "manga/" + mangaName)
	ExitIfExists(mangaSaveDir + "/" + "manga/" + mangaName + "/" + chapterNumber)
	NewDir(mangaSaveDir + "/" + "manga/" + mangaName + "/" + chapterNumber)

	fmt.Println("Trying to load images for Chapter" + chapterNumber)
	fmt.Println("Downloading", len(images), "pages")
	for imageNumber, image := range images {
		imageName := strings.Split(image, "/")
		imageName = strings.Split(imageName[len(imageName)-1], ".")
		imageName[0] = strings.Replace(imageName[0], " ", "_", -1)
		imageFullDir := mangaSaveDir + "manga/" + mangaName + "/" + chapterNumber + "/" + strconv.Itoa(imageNumber+1) + "." + imageName[1]
		SaveImage(image, imageFullDir)
	}
}

func ComicSearch() {
	fmt.Println("Not implemented")
	os.Exit(1)
}
