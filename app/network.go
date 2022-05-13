package app

import (
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
