package app

import (
	"fmt"
	"strings"

	"github.com/byte-cats/filechick"
)

// Search makes a search request to the MangaDex API
func MkSearch(basedApiUrl string, query string) {
	var mangaList MangaList
	var chapterTitles []string
	var mangaChapters MangaChapters
	var mangaImages MangaImages
	var images []string

	mangaSaveDir := "./"

	// Update to use the correct MangaDex API v5 endpoint for manga search
	apiSearch := fmt.Sprintf("%s/manga?title=%s&limit=20", basedApiUrl, query)
	fmt.Println("Searching using API endpoint:", apiSearch)

	results, err := CustomRequest(apiSearch)
	if err != nil {
		fmt.Println("Error fetching manga:", err)
		return
	}

	fmt.Println("Searching for:", query)
	ParseMangaSearch(results, &mangaList)
	fmt.Printf("Found %d manga\n", len(mangaList))

	titles := formatMangaTitles(mangaList) // Extracted formatting logic
	printTitles(titles)

	mangaChoiceInt := getUserChoice(len(titles)) // Extracted user choice logic
	if mangaChoiceInt < 1 {
		return
	}

	mangaId := mangaList[mangaChoiceInt-1].ID
	// Update to use the correct MangaDex API v5 endpoint for chapters
	chapterUrl := fmt.Sprintf("%s/manga/%s/feed?translatedLanguage[]=en&limit=100", basedApiUrl, mangaId)
	fmt.Println("Loading chapters...")

	results, err = CustomRequest(chapterUrl)
	if err != nil {
		fmt.Println("Error fetching chapters:", err)
		return
	}
	ParseChapters(results, &mangaChapters)

	chapterTitles = formatChapterTitles(mangaChapters) // Extracted formatting logic
	printTitles(chapterTitles)

	chapterChoiceInt := getUserChoice(len(chapterTitles)) // Extracted user choice logic
	if chapterChoiceInt < 1 {
		return
	}

	chapterId := mangaChapters.ChapterID[len(mangaChapters.ChapterID)-chapterChoiceInt]
	chapterNumber := strings.TrimPrefix(chapterId, "chapter-")
	fmt.Println("Chapter number:", chapterNumber)

	// Update to use the correct MangaDex API v5 endpoint for at-home server
	// This endpoint provides the actual image URLs
	imagesUrl := fmt.Sprintf("%s/at-home/server/%s", basedApiUrl, chapterId)
	results, err = CustomRequest(imagesUrl)
	if err != nil {
		fmt.Println("Error fetching images:", err)
		return
	}
	ParseImages(results, &mangaImages)

	images = extractImageUrls(mangaImages) // Extracted image URL logic

	// Get the title from the new manga structure
	var mangaTitle string
	if title, ok := mangaList[mangaChoiceInt-1].Attributes.Title["en"]; ok {
		mangaTitle = title
	} else {
		// Fallback to the first available title
		for _, title := range mangaList[mangaChoiceInt-1].Attributes.Title {
			mangaTitle = title
			break
		}
	}

	mangaName := sanitizeMangaName(mangaTitle) // Extracted sanitization logic

	prepareDirectories(mangaSaveDir, mangaName, chapterNumber) // Extracted directory preparation logic

	fmt.Println("Trying to load pages for Chapter " + chapterNumber)
	fmt.Printf("Downloading %d pages\n", len(images))
	downloadImages(images, mangaSaveDir, mangaName, chapterNumber) // Extracted download logic
}

// Helper functions

func formatMangaTitles(mangaList MangaList) []string {
	var titles []string
	for i, manga := range mangaList {
		// Get the title in the best available language (prefer English)
		var title string
		if t, ok := manga.Attributes.Title["en"]; ok {
			title = t
		} else {
			// Fallback to the first available title
			for _, t := range manga.Attributes.Title {
				title = t
				break
			}
		}
		titles = append(titles, fmt.Sprintf("%d. %s", i+1, title))
	}
	return titles
}

func printTitles(titles []string) {
	for _, title := range titles {
		fmt.Println(title)
	}
}

func getUserChoice(max int) int {
	fmt.Printf("Select a title: (1 - %d) ", max)
	mangaChoice := filechick.GetInput()
	QueryCheck(mangaChoice)

	mangaChoiceInt := filechick.StringToInt(mangaChoice)
	if mangaChoiceInt < 1 || mangaChoiceInt > max {
		fmt.Println("Invalid choice")
		return -1
	}
	return mangaChoiceInt
}

func extractImageUrls(mangaImages MangaImages) []string {
	var images []string
	for _, image := range mangaImages {
		images = append(images, image.ImageUrl)
	}
	return images
}

func sanitizeMangaName(title string) string {
	return strings.NewReplacer(" ", "_", ":", "").Replace(title) // More efficient sanitization
}

func prepareDirectories(baseDir, mangaName, chapterNumber string) {
	filechick.NewDir(baseDir + "/manga")
	filechick.NewDir(baseDir + "/manga/" + mangaName)
	filechick.ExitIfExists(baseDir + "/manga/" + mangaName + "/" + chapterNumber)
	filechick.NewDir(baseDir + "/manga/" + mangaName + "/" + chapterNumber)
}

func downloadImages(images []string, baseDir, mangaName, chapterNumber string) {
	for imageNumber, image := range images {
		imageName := strings.Split(image, "/")
		ext := strings.Split(imageName[len(imageName)-1], ".")[1]
		imageFullDir := fmt.Sprintf("%s/manga/%s/%s/%d.%s", baseDir, mangaName, chapterNumber, imageNumber+1, ext)
		ProgressBar(imageNumber, len(images))
		filechick.SaveImage(image, imageFullDir)
	}
}

// ProgressBar is a simple progress bar
func ProgressBar(imageNumber int, lenImages int) {
	fmt.Printf("\r\033[38;5;%dm[%-50s]\033[0m %d%%", 1+imageNumber%255, strings.Repeat("=", imageNumber/2), imageNumber*100/lenImages)
}

func formatChapterTitles(mangaChapters MangaChapters) []string {
	var titles []string
	for i, chapter := range mangaChapters.Chapters {
		titles = append(titles, fmt.Sprintf("%d. %s", i+1, chapter)) // Changed from chapter.Title to chapter
	}
	return titles
}
