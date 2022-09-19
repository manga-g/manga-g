package app

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// GetInput Function to get user input from the command line
func GetInput() string {
	var input string
	in := bufio.NewScanner(os.Stdin)

	in.Scan()
	err := in.Err()
	input = in.Text()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	input = strings.Replace(input, "\n", "", -1)

	return input
}

// ValidateUrl checks if url is valid
func ValidateUrl(UrlToCheck string) bool {
	_, err := url.ParseRequestURI(UrlToCheck)
	if err != nil {
		return false
	}
	return true
}

// StringToInt to change string to int
func StringToInt(str string) int {
	var i int
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return 0
	}
	return i
}

// Retry if there is no input, loop the request 3 times
func Retry(query string) {
	for n := 0; n < 3; n++ {
		fmt.Print("\nYou should choose the number corresponding to the manga you want to read.\nTry again,please :)\n" + "Search for manga: ")
		query = GetInput()
		if query != "" {
			break
		}
		if n == 2 {
			os.Exit(0)
		}
	}
}
func QueryCheck(query string) {
	if query == "" || query == " " {
		Retry(query)
	}
}
