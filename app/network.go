package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Connected tests if the connection is alive
func Connected() bool {
	_, err := http.Get("https://clients3.google.com/generate_204")
	return err == nil
}

// ValidateURL checks if the provided URL is valid
func ValidateURL(urlToCheck string) bool {
	_, err := url.ParseRequestURI(urlToCheck)
	return err == nil
}

// CheckAPI checks if the API is alive
func CheckAPI(baseAPIURL string) {
	resp, err := http.Get(baseAPIURL + "/hc")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Online Manga is offline at the moment ;(\nBe back online ASAP =)")
		os.Exit(1)
	}
	fmt.Println("Online Manga is ready!")
}

// RandomizeUserAgent returns a random user agent string
func RandomizeUserAgent() string {
	return "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36"
}

// ApplyUserAgent applies the user agent to an HTTP request
func ApplyUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", RandomizeUserAgent())
}

// CustomRequest performs a custom HTTP GET request
func CustomRequest(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "Something went wrong with request", err
	}
	ApplyUserAgent(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Something went wrong with the request", err
	}
	defer closeBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading body:", err
	}
	return string(body), nil
}

// closeBody closes the response body and handles any errors
func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		fmt.Println("Error closing body:", err)
	}
}
