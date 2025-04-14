package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	_ "image/jpeg" // Import for JPEG decoding
	_ "image/png"  // Import for PNG decoding
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/manga-g/manga-g/app"
	"github.com/manga-g/manga-g/internal/db" // New import for DB package
)

// Define download state tracking struct
type DownloadState struct {
	MangaID        string   `json:"manga_id"`
	MangaTitle     string   `json:"manga_title"`
	ChapterID      string   `json:"chapter_id"`
	ChapterHash    string   `json:"chapter_hash"`
	BaseURL        string   `json:"base_url"`
	Files          []string `json:"files"`
	Downloaded     []bool   `json:"downloaded"`
	LastUpdated    string   `json:"last_updated"`
	TotalFiles     int      `json:"total_files"`
	CompletedFiles int      `json:"completed_files"`
}

// Global variables for download management
var (
	downloadState  DownloadState
	downloadMutex  sync.Mutex
	downloadActive bool
	downloadChan   chan int
	workerWg       sync.WaitGroup
)

// saveDownloadState writes the current download state to a JSON file
// IMPORTANT: Caller MUST hold downloadMutex before calling this function.
func saveDownloadState() error {
	// downloadMutex.Lock() // Removed lock acquisition
	// defer downloadMutex.Unlock() // Removed lock release

	downloadState.LastUpdated = time.Now().Format(time.RFC3339)

	stateDir := filepath.Dir(getStateFilePath())
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return err
	}

	stateData, err := json.MarshalIndent(downloadState, "", "  ")
	if err != nil {
		return err
	}

	// Use a temporary file for atomic write
	tmpStateFile := getStateFilePath() + ".tmp"
	err = os.WriteFile(tmpStateFile, stateData, 0644)
	if err != nil {
		os.Remove(tmpStateFile) // Clean up tmp file on error
		return err
	}

	// Rename temporary file to final state file
	return os.Rename(tmpStateFile, getStateFilePath())
}

// loadDownloadState attempts to load a previous download state from disk
func loadDownloadState() (bool, error) {
	stateFile := getStateFilePath()

	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return false, nil // No state file exists
	}

	stateData, err := os.ReadFile(stateFile)
	if err != nil {
		return false, err
	}

	downloadMutex.Lock()
	defer downloadMutex.Unlock()

	if err := json.Unmarshal(stateData, &downloadState); err != nil {
		return false, err
	}

	return true, nil
}

// getStateFilePath returns the path to the download state file
func getStateFilePath() string {
	return "debug_downloads/download_state.json"
}

func main() {
	// Define command-line flags
	searchQuery := flag.String("search", "", "Search for manga by title")
	mangaID := flag.String("manga", "", "Get manga details by ID")
	chapterID := flag.String("chapter", "", "Get chapter details by ID")
	testMode := flag.Bool("test", false, "Run in test mode (no API calls)")
	mockMode := flag.Bool("mock", false, "Run with mock data instead of real API calls")
	verbose := flag.Bool("v", false, "Verbose output (print raw JSON)")
	apiURL := flag.String("api", "https://api.mangadex.org", "MangaDex API base URL")
	flag.Parse()

	err := db.InitDB()
	if err != nil {
		fmt.Printf("Error initializing database: %v. Falling back to API only.\n", err)
	}
	if flag.NFlag() == 0 {
		fmt.Print("Enter a manga search name: ")
		var searchQuery string
		fmt.Scanln(&searchQuery)
		if searchQuery != "" {
			searchManga(*apiURL, searchQuery, *verbose)
			return // Exit after interactive search
		} else {
			fmt.Println("No search query provided.")
			printUsage()
			return
		}
	}
	fmt.Printf("Using API: %s\n", *apiURL)

	// Execute the requested command
	switch {
	case *testMode:
		printApiEndpoints(*apiURL)
	case *mockMode && *searchQuery != "":
		app.MockSearchManga(*searchQuery, *verbose)
	case *mockMode && *mangaID != "":
		app.MockMangaDetails(*mangaID, *verbose)
	case *searchQuery != "":
		searchManga(*apiURL, *searchQuery, *verbose)
	case *mangaID != "":
		getMangaDetails(*apiURL, *mangaID, *verbose)
	case *chapterID != "":
		getChapterDetails(*apiURL, *chapterID, *verbose)
	default:
		fmt.Println("No valid command specified.")
		printUsage()
	}
}

func printUsage() {
	fmt.Println("MangaDex API Debug Tool")
	fmt.Println("----------------------")
	fmt.Println("Usage:")
	fmt.Println("  -search <title>   Search for manga by title")
	fmt.Println("  -manga <id>       Get manga details by ID")
	fmt.Println("  -chapter <id>     Get chapter details by ID")
	fmt.Println("  -test             Run in test mode (no API calls)")
	fmt.Println("  -mock             Use mock data instead of real API calls")
	fmt.Println("  -v                Verbose output (print raw JSON)")
	fmt.Println("  -api <url>        Use custom API URL (default: https://api.mangadex.org)")
	fmt.Println("\nExamples:")
	fmt.Println("  go run cmd/debug/main.go -search \"one piece\"")
	fmt.Println("  go run cmd/debug/main.go -manga \"some-manga-id\"")
	fmt.Println("  go run cmd/debug/main.go -chapter \"some-chapter-id\" -v")
	fmt.Println("  go run cmd/debug/main.go -mock -search \"any title\"")
}

func searchManga(apiURL, query string, verbose bool) {
	// First, check internal DB
	resultsFromDB := db.GetMangaFromDB(query) // Assume this returns cached results
	if len(resultsFromDB) > 0 {
		fmt.Println("Using data from internal database:")
		// Process and display DB results similarly to API results
		for i, manga := range resultsFromDB {
			fmt.Printf("[%d] ID: %s - Title: %s\n", i+1, manga.ID, manga.Title)
		}
	} else {
		// Fall back to API if not in DB
		fmt.Println("No data in internal DB, querying API...")
		encodedQuery := url.QueryEscape(query)
		apiSearch := fmt.Sprintf("%s/manga?title=%s&limit=5", apiURL, encodedQuery)
		fmt.Printf("Searching for: %s\n", query)
		fmt.Printf("API Request: %s\n", apiSearch)

		results, err := app.CustomRequest(apiSearch)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			printFormattedJSON(results)
		}

		var mangaList app.MangaList
		app.ParseMangaSearch(results, &mangaList)

		fmt.Printf("Found %d manga results:\n", len(mangaList))
		for i, manga := range mangaList {
			title := mangaList.GetTitle(i)
			fmt.Printf("[%d] ID: %s - Title: %s\n", i+1, manga.ID, title)

			// Print some tags if available
			if len(manga.Attributes.Tags) > 0 {
				var tags []string
				for _, tag := range manga.Attributes.Tags[:min(3, len(manga.Attributes.Tags))] {
					if name, ok := tag.Attributes.Name["en"]; ok {
						tags = append(tags, name)
					}
				}
				fmt.Printf("    Tags: %s\n", strings.Join(tags, ", "))
			}

			// Print description snippet if available
			if desc, ok := manga.Attributes.Description["en"]; ok && desc != "" {
				// Truncate description to first 100 chars
				if len(desc) > 100 {
					desc = desc[:100] + "..."
				}
				fmt.Printf("    Desc: %s\n", desc)
			}
			fmt.Println()
		}

		// Add selection prompt if there are results
		if len(mangaList) > 0 {
			fmt.Print("Select a manga to view details (enter number, or 0 to exit): ")
			var choice int
			_, err := fmt.Scanf("%d", &choice)
			if err != nil || choice < 0 || choice > len(mangaList) {
				fmt.Println("Invalid selection")
				return
			}

			if choice > 0 {
				selectedID := mangaList[choice-1].ID
				fmt.Println("\n--- Fetching details for selected manga ---")
				getMangaDetails(apiURL, selectedID, verbose)
			}
		}
	}
}

func getMangaDetails(apiURL, mangaID string, verbose bool) {
	// First, check internal DB
	resultFromDB := db.GetMangaDetailFromDB(mangaID)
	if resultFromDB.ID != "" {
		fmt.Println("Using data from internal database:")
		fmt.Printf("ID: %s\nTitle: %s\n", resultFromDB.ID, resultFromDB.Title)
		// Add more fields as needed
	} else {
		fmt.Println("No data in internal DB, querying API...")
		apiEndpoint := fmt.Sprintf("%s/manga/%s", apiURL, mangaID)
		fmt.Printf("Getting manga details for ID: %s\n", mangaID)
		fmt.Printf("API Request: %s\n", apiEndpoint)

		results, err := app.CustomRequest(apiEndpoint)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			printFormattedJSON(results)
		} else {
			// Parse just enough to show useful info
			var response struct {
				Data struct {
					ID         string `json:"id"`
					Attributes struct {
						Title       map[string]string `json:"title"`
						Description map[string]string `json:"description"`
						Status      string            `json:"status"`
						Year        int               `json:"year"`
					} `json:"attributes"`
				} `json:"data"`
			}

			if err := json.Unmarshal([]byte(results), &response); err != nil {
				fmt.Printf("Error parsing JSON: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("ID: %s\n", response.Data.ID)
			fmt.Printf("Title: %s\n", getFirstValue(response.Data.Attributes.Title))
			fmt.Printf("Status: %s\n", response.Data.Attributes.Status)
			fmt.Printf("Year: %d\n", response.Data.Attributes.Year)
			fmt.Printf("Description: %s\n", getFirstValue(response.Data.Attributes.Description))
		}

		// Get chapters for this manga
		apiEndpoint = fmt.Sprintf("%s/manga/%s/feed?translatedLanguage[]=en&limit=10", apiURL, mangaID)
		fmt.Printf("\nGetting latest chapters for manga ID: %s\n", mangaID)
		fmt.Printf("API Request: %s\n", apiEndpoint)

		chapterResults, err := app.CustomRequest(apiEndpoint)
		if err != nil {
			fmt.Printf("Error getting chapters: %v\n", err)
			// Decide whether to exit or continue
			return
		}

		var mangaChapters app.MangaChapters
		app.ParseChapters(chapterResults, &mangaChapters)

		fmt.Printf("\nFound %d chapters:\n", len(mangaChapters.Chapters))
		for i, chapter := range mangaChapters.Chapters[:min(5, len(mangaChapters.Chapters))] {
			fmt.Printf("[%d] %s (ID: %s)\n", i+1, chapter, mangaChapters.ChapterID[i])
		}

		// Add chapter selection prompt
		if len(mangaChapters.Chapters) > 0 {
			fmt.Print("\nSelect a chapter to view details (enter number, or 0 to exit): ")
			var choice int
			_, scanErr := fmt.Scanf("%d", &choice)
			if scanErr != nil || choice < 0 || choice > len(mangaChapters.Chapters) {
				fmt.Println("Invalid selection")
				return
			}

			if choice > 0 {
				selectedIndex := choice - 1
				if selectedIndex >= len(mangaChapters.ChapterID) {
					fmt.Println("Chapter index out of range")
					return
				}

				selectedID := mangaChapters.ChapterID[selectedIndex]
				fmt.Println("\n--- Fetching details for selected chapter ---")
				getChapterDetails(apiURL, selectedID, verbose)
			}
		}
	}
}

func getChapterDetails(apiURL, chapterID string, verbose bool) {
	var mangaID string
	var mangaTitle string

	// First, get chapter info to find the manga ID
	apiEndpoint := fmt.Sprintf("%s/chapter/%s", apiURL, chapterID)
	fmt.Printf("Getting chapter details for ID: %s\n", chapterID)
	fmt.Printf("API Request: %s\n", apiEndpoint)

	results, err := app.CustomRequest(apiEndpoint)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		printFormattedJSON(results)
	}

	// Extract manga ID from chapter
	var chapterResponse struct {
		Data struct {
			Relationships []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(results), &chapterResponse); err == nil {
		for _, rel := range chapterResponse.Data.Relationships {
			if rel.Type == "manga" {
				mangaID = rel.ID
				break
			}
		}
	}

	// Get manga title if we found an ID
	if mangaID != "" {
		apiEndpoint = fmt.Sprintf("%s/manga/%s", apiURL, mangaID)
		results, err := app.CustomRequest(apiEndpoint)
		if err == nil {
			var mangaResponse struct {
				Data struct {
					Attributes struct {
						Title map[string]string `json:"title"`
					} `json:"attributes"`
				} `json:"data"`
			}

			if err := json.Unmarshal([]byte(results), &mangaResponse); err == nil {
				mangaTitle = getFirstValue(mangaResponse.Data.Attributes.Title)
			}
		}
	}

	// If we couldn't get the title, use a default
	if mangaTitle == "" {
		mangaTitle = "Unknown_Manga"
	}

	// Then get the images (at-home server)
	fmt.Println("\nGetting image server information:")
	apiAtHomeEndpoint := fmt.Sprintf("%s/at-home/server/%s", apiURL, chapterID)
	fmt.Printf("API Request: %s\n", apiAtHomeEndpoint)

	atHomeResults, err := app.CustomRequest(apiAtHomeEndpoint)
	if err != nil {
		fmt.Printf("Error getting at-home server: %v\n", err)
		os.Exit(1)
	}

	var atHome app.AtHomeResponse
	if err := json.Unmarshal([]byte(atHomeResults), &atHome); err != nil {
		fmt.Printf("Error parsing at-home JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Base URL: %s\n", atHome.BaseURL)
	fmt.Printf("Chapter Hash: %s\n", atHome.Chapter.Hash)
	fmt.Printf("Image Count: %d\n", len(atHome.Chapter.Data))

	if len(atHome.Chapter.Data) == 0 {
		fmt.Println("This manga source is not available on the MangaDex API at this time. More sources will be added later (TBD WIP).")
		return // Cannot download if no images
	}

	fmt.Printf("First image: %s/data/%s/%s\n",
		atHome.BaseURL, atHome.Chapter.Hash, atHome.Chapter.Data[0])

	if verbose {
		printFormattedJSON(atHomeResults)
	}

	// Prompt to download images
	fmt.Print("\nDownload images? (y/n): ")
	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		fmt.Print("Number of download threads (1-10): ")
		var threads int
		fmt.Scanf("%d", &threads)
		if threads < 1 || threads > 10 {
			threads = 3
			fmt.Println("Using default of 3 threads")
		}

		// Create and use the Downloader instance
		// Pass nil for the progress channel as debug doesn't use it
		downloader := app.NewDownloader(threads, "", "", nil)

		fmt.Println("Starting download...")
		downloadErr := downloader.DownloadChapter(&atHome, mangaTitle, chapterID)

		if downloadErr != nil {
			fmt.Printf("Download encountered an error: %v\n", downloadErr)
		} else {
			// Download succeeded, prompt for PDF
			sanitizedTitle := app.SanitizeFilename(mangaTitle) // Use app version
			// Construct the expected output dir based on downloader defaults/logic
			// TODO: Ideally, DownloadChapter would return the output path
			outputDir := filepath.Join(downloader.OutputDirBase, sanitizedTitle, fmt.Sprintf("chapter_%s", atHome.Chapter.Hash[:8]))

			fmt.Print("Package images into a PDF? (y/n): ")
			var pdfAnswer string
			fmt.Scanln(&pdfAnswer)
			if strings.ToLower(pdfAnswer) == "y" || strings.ToLower(pdfAnswer) == "yes" {
				fmt.Println("Starting PDF packaging...")
				pdfErr := app.PackageChapterToPDF(outputDir, sanitizedTitle, atHome.Chapter.Hash[:8])
				if pdfErr != nil {
					fmt.Printf("Error creating PDF: %v\n", pdfErr)
				} // Success message handled in PackageChapterToPDF
			}
		}
	}
}

func printApiEndpoints(apiURL string) {
	fmt.Println("MangaDex API v5 Endpoints:")
	fmt.Println("==========================")

	endpoints := []struct {
		Name        string
		Endpoint    string
		Description string
	}{
		{"Search Manga", "/manga?title={title}&limit={limit}", "Search for manga by title"},
		{"Manga Details", "/manga/{id}", "Get detailed information about a manga"},
		{"Manga Chapters", "/manga/{id}/feed?translatedLanguage[]={lang}&limit={limit}", "Get chapters for a manga"},
		{"Chapter Details", "/chapter/{id}", "Get detailed information about a chapter"},
		{"At-Home Server", "/at-home/server/{id}", "Get image server information for a chapter"},
		{"Cover Art", "/cover/{coverId}", "Get cover art for a manga"},
		{"Author", "/author/{id}", "Get information about an author"},
		{"Scanlation Group", "/group/{id}", "Get information about a scanlation group"},
	}

	for _, e := range endpoints {
		fmt.Printf("- %s\n", e.Name)
		fmt.Printf("  %s%s\n", apiURL, e.Endpoint)
		fmt.Printf("  %s\n\n", e.Description)
	}
}

// Helper function to pretty-print JSON
func printFormattedJSON(jsonStr string) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(jsonStr), "", "  "); err != nil {
		fmt.Printf("Error formatting JSON: %v\nRaw: %s\n", err, jsonStr)
		return
	}
	fmt.Println(prettyJSON.String())
}

// Helper function to get the first value from a map
func getFirstValue(m map[string]string) string {
	if val, ok := m["en"]; ok {
		return val
	}
	for _, v := range m {
		return v
	}
	return "N/A"
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
