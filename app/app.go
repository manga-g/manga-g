package app

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// GetInput Function to get user input from the command line
func GetInput() string {
	var input string
	in := bufio.NewReader(os.Stdin)

	input, err := in.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}

	// serialize the input
	input = url.QueryEscape(input)
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
