package app

import (
	"fmt"
	"strconv"
	"strings"
)

func MkSearch(basedApiUrl string, query string) {
	mangaSaveDir := "./"
	fmt.Println("Searching for:", query)
	apiSearch := basedApiUrl + "mk/search?q=" + query
	res, _ := CustomRequest(apiSearch)
	// wait for the results to be ready
	var mangaList MangaList
	ParseMangaSearch(res, &mangaList)
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
	mangaChoice := GetInput()
	//if there is no input, loop the request 3 times
	if mangaChoice == "" {
		Retry(mangaChoice)
	}

	mangaChoiceInt := StringToInt(mangaChoice)
	if mangaChoiceInt > number {
		fmt.Println("Invalid choice")
		return
	}
	mangaId := mangaList[mangaChoiceInt-1].ID
	chapterUrl := basedApiUrl + "mk/chapters?q=" + mangaId
	fmt.Println("Checking ID:" + mangaId)
	fmt.Println("Loading chapters...")
	res, _ = CustomRequest(chapterUrl)
	var mangaChapters MangaChapters
	ParseChapters(res, &mangaChapters)
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
	chapterChoice := GetInput()
	//if there is no input, loop the request 3 times
	if chapterChoice == "" {
		for n := 0; n < 3; n++ {
			fmt.Println("You should choose the number corresponding to the chapter you want to read.\nTry again,please :)" + "Select a chapter: (1 - " + strconv.Itoa(len(chapterTitles)) + ") ")
			chapterChoice = GetInput()
			if chapterChoice != "" {
				break
			}
		}
	}
	chapterChoiceInt := StringToInt(chapterChoice)
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
	res, _ = CustomRequest(imagesUrl)
	var mangaImages MangaImages
	ParseImages(res, &mangaImages)
	var images []string
	for _, image := range mangaImages {
		images = append(images, image.ImageUrl)
	}

	NewDir(mangaSaveDir + "/" + "manga")

	mangaName := strings.Replace(mangaList[mangaChoiceInt-1].Title, " ", "_", -1)
	mangaName = strings.Replace(mangaName, ":", "", -1)
	mangaName = strings.Replace(mangaName, " ", "_", -1)

	NewDir(mangaSaveDir + "/" + "manga/" + mangaName)
	NewDir(mangaSaveDir + "/" + "manga/" + mangaName + "/" + chapterNumber)

	fmt.Println("Downloading", len(images), "pages")
	//fmt.Println(images)
	for _, image := range images {
		imageName := strings.Split(image, "/")
		imageName = strings.Split(imageName[len(imageName)-1], ".")
		imageName[0] = strings.Replace(imageName[0], " ", "_", -1)
		imageFullDir := mangaSaveDir + "manga/" + mangaName + "/" + chapterNumber + "/" + imageName[0] + "." + imageName[1]
		SaveImage(image, imageFullDir)
	}
}
