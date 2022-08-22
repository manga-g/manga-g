package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Connected tests if the connection is alive
func Connected() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

// RandomizeUserAgent randomize the user agent
func RandomizeUserAgent() string {
	return "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36"
}

// ApplyUserAgent function to apply the user agent to a http request
func ApplyUserAgent(req *http.Request) {
	req.Header.Add("User-Agent", RandomizeUserAgent())
}

func CustomRequest(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "Something went wrong with request", err
	}
	ApplyUserAgent(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Something went wrong with the request", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error reading body:", err
	}
	return string(body), nil
}

func CheckApi(basedApiUrl string) {
	_, err := http.Get(basedApiUrl)
	if err == nil {
		fmt.Println("Online Manga is ready!")
	} else {
		fmt.Println("Online Manga is offline at the moment ;(\nBe back online ASAP =)")
		os.Exit(1)
	}
}
