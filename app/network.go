package app

import (
    "io/ioutil"
    "net/http"
)

// Connected tests if the connection is alive
func (app *MG) Connected() (ok bool) {
    _, err := http.Get("http://clients3.google.com/generate_204")
    if err != nil {
        return false
    }
    return true
}

// RandomizeUserAgent randomize the user agent
func (app *MG) RandomizeUserAgent() string {
    return "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36"
}

// ApplyUserAgent function to apply the user agent to a http request
func (app *MG) ApplyUserAgent(req *http.Request) {
    req.Header.Add("User-Agent", app.RandomizeUserAgent())
}

func (app *MG) CustomRequest(url string) (string, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }
    app.ApplyUserAgent(req)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}
