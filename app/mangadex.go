package app

import (
	"fmt"
	"regexp"
)

func (app *App) GetMangaDex(html string) ([]string, error) {
	var manga []string
	reg := regexp.MustCompile("blob.*?\"")
	for i, match := range reg.FindAllStringSubmatch(html, -1) {
		fmt.Println(i)
		fmt.Println("match:", match[0])

	}
	return manga, nil
}
