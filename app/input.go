package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	input = strings.Replace(input, "\n", "", -1)

	return input
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
