package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	// Test searching for a manga
	mangaID := "a96676e5-8ae2-425e-b549-7f15dd34a6d8" // Example manga ID (One Piece)

	// Get chapters for this manga
	chapterURL := fmt.Sprintf("https://api.mangadex.org/manga/%s/feed?translatedLanguage[]=en&limit=10", mangaID)
	fmt.Println("Requesting:", chapterURL)

	results, err := customRequest(chapterURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Pretty print the JSON response
	var prettyJSON map[string]interface{}
	err = json.Unmarshal([]byte(results), &prettyJSON)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	prettyOutput, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println("API Response:")
	fmt.Println(string(prettyOutput))
}

// CustomRequest performs a custom HTTP GET request with timeout and error handling
func customRequest(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set a user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36")

	client := &http.Client{
		Timeout: 30 * time.Second, // Increased timeout for API requests
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}
	return string(body), nil
}
