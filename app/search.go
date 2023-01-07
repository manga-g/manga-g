package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/byte-cats/filechick"
)

type Api struct {
	ApiUrl  string
	SaveDir string
}

func NewApi() *Api {
	return &Api{
		ApiUrl:  os.Getenv("API_URL"),
		SaveDir: os.Getenv("DEFAULT_DOWNLOAD_PATH"),
	}
}

// Prolly is better to create a DTO for each function?

type Results struct {
	List          MangaList
	Titles        []string
	ChosenManga   int
	Chapters      MangaChapters
	ChapterTitles []string
	ChosenChapter int
	Images        MangaImages
	ImageNames    []string
}

func NewResults(api *Api, query string) *Results {
	mangaList, titles := api.GetMangas(query)
	chosenManga, chapters := api.GetChapters(titles, mangaList)
	chapterTitles := api.GetChapterTitles(chapters)
	chosenChapter, images, imageNames := api.GetChapterImages(chapterTitles, chapters)
	return &Results{
		List:          mangaList,
		Titles:        titles,
		ChosenManga:   chosenManga,
		Chapters:      chapters,
		ChapterTitles: chapterTitles,
		ChosenChapter: chosenChapter,
		Images:        images,
		ImageNames:    imageNames,
	}
}

// MkSearch makes a search request to the api at mk endpoint
func MkSearch(query string) *Results {
	api := NewApi()
	results := NewResults(api, query)
	return results
}

func (a *Api) GetMangas(query string) (MangaList, []string) {
	apiSearch := a.ApiUrl + "mk/search?q=" + query
	mangaResults, _ := CustomRequest(apiSearch)

	fmt.Println("Searching for:", query)
	var mangaList MangaList
	ParseMangaSearch(mangaResults, &mangaList)
	fmt.Println("Found ", len(mangaList), " manga")
	var titles []string
	for i, manga := range mangaList {
		titles = append(titles, fmt.Sprintf("%d. %s", i+1, manga.Title))
	}
	return mangaList, titles
}

func (a *Api) GetChapters(titles []string, mangaList MangaList) (int, MangaChapters) {
	number := len(titles)
	for _, title := range titles {
		fmt.Println(title)
	}

	SelectMessage := "Select a title: (1 - " + strconv.Itoa(number) + ") "
	fmt.Print(SelectMessage)
	mangaChoice := filechick.GetInput()
	QueryCheck(mangaChoice)

	chosenManga := filechick.StringToInt(mangaChoice)
	if chosenManga > number {
		fmt.Println("Invalid choice")
		os.Exit(1)
	}
	mangaId := mangaList[chosenManga-1].ID
	chapterUrl := a.ApiUrl + "mk/chapters?q=" + mangaId
	fmt.Println("Checking ID:" + mangaId)
	fmt.Println("Loading chapters...")

	var mangaChapters MangaChapters
	chapterResults, _ := CustomRequest(chapterUrl)
	ParseChapters(chapterResults, &mangaChapters)
	return chosenManga, mangaChapters
}

func (a *Api) GetChapterTitles(mangaChapters MangaChapters) []string {
	n := 0
	var chapterTitles []string
	for i := len(mangaChapters.Chapters); i >= 1; i-- {
		chapter := mangaChapters.Chapters[i-1]
		chapterTitles = append(chapterTitles, fmt.Sprintf("%d. %s", n+1, chapter))
		n++
	}

	for _, title := range chapterTitles {
		fmt.Println(title)
	}

	return chapterTitles
}

func (a *Api) GetChapterImages(chapterTitles []string, mangaChapters MangaChapters) (int, MangaImages, []string) {
	fmt.Print("Select a result: (1 - " + strconv.Itoa(len(chapterTitles)) + ") ")
	resultChoice := filechick.GetInput()
	QueryCheck(resultChoice)

	chapterChoiceInt := filechick.StringToInt(resultChoice)
	chapterChoiceInt = len(mangaChapters.ChapterID) - chapterChoiceInt
	if chapterChoiceInt > len(chapterTitles) || chapterChoiceInt < 1 {
		fmt.Println("Invalid choice")
		os.Exit(1)
	}
	chapterId := mangaChapters.ChapterID[chapterChoiceInt]
	chapterNumber := strings.Replace(chapterId, "chapter-", "", -1)
	fmt.Println("Chapter number:", chapterNumber)

	// keep only the number at the end of the string
	imagesUrl := a.ApiUrl + "mk/images?id=" + chapterId + "&chapterNumber=" + chapterNumber
	//	fmt.Println(imagesUrl)
	var images MangaImages
	var imageResults []string
	result, _ := CustomRequest(imagesUrl)
	ParseImages(result, &images)
	for _, image := range images {
		imageResults = append(imageResults, image.ImageUrl)
	}

	filechick.NewDir(a.SaveDir)

	// mangaName := strings.Replace(results.List[results.ChosenManga-1].Title, " ", "_", -1)
	// mangaName = strings.Replace(mangaName, ":", "", -1)
	// mangaName = strings.Replace(mangaName, " ", "_", -1)

	// filechick.NewDir(a.SaveDir + "/" + mangaName)
	// filechick.ExitIfExists(a.SaveDir + "/" + mangaName + "/" + chapterNumber)
	// filechick.NewDir(a.SaveDir + "/" + mangaName + "/" + chapterNumber)

	fmt.Println("Trying to load pages for Chapter " + chapterNumber)
	fmt.Println("Downloading", len(imageResults), "pages")
	for imageNumber, image := range imageResults {
		imageName := strings.Split(image, "/")
		imageName = strings.Split(imageName[len(imageName)-1], ".")
		imageName[0] = strings.Replace(imageName[0], " ", "_", -1)
		ProgressBar(imageNumber, len(imageResults))
		imageFullDir := a.SaveDir + "manga/" + "yeah" + "/" + chapterNumber + "/" + strconv.Itoa(imageNumber+1) + "." + imageName[1]
		filechick.SaveImage(image, imageFullDir)
	}
	return chapterChoiceInt, images, imageResults
}

// ProgressBar is a simple progress bar
func ProgressBar(imageNumber int, lenImages int) {
	// fancy rainbow multicolored progress bar with percentage at the end
	fmt.Printf("\r\033[38;5;%dm[%-50s]\033[0m %d%%", 1+imageNumber%255, strings.Repeat("=", imageNumber/2), imageNumber*100/lenImages)
}

// ComicSearch makes a search request to the api at comic endpoint
// func ComicSearch() {
// 	fmt.Println("Not implemented")
// 	os.Exit(1)
// }
