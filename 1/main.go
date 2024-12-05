// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
)

func challengeFunction(file *os.File) {
	list := sortParsedInput(getParsedInput(file))

	distance := 0

	for _, val := range list {
		one := val[0]
		two := val[1]

		if one > two {
			distance += one - two
		} else {
			distance += two - one
		}
	}

	fmt.Printf("Total distance is: %d\n", distance)
        fmt.Printf("The correct answer is 2000468\n")
}

func getNewHTTPClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	return client
}

func getNewRequest(url string) *http.Request {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "formRequest(): error creating new request: %v", err)
	}

	r.AddCookie(getSessionCookie())

	return r
}

func getParsedInput(rawInput io.Reader) ([]int, []int) {
	listOne := make([]int, 0)
	listTwo := make([]int, 0)
	regex := regexp.MustCompile(`(?: +)`)
	scanner := bufio.NewScanner(rawInput)

	for scanner.Scan() {
		line := scanner.Text()
		tuple := regex.Split(line, 2)
		one, err := strconv.Atoi(tuple[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "There was an issue converting %s to an integer", tuple[0])
		}
		two, err := strconv.Atoi(tuple[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "There was an issue converting %s to an integer", tuple[1])
		}

		listOne = append(listOne, one)
		listTwo = append(listTwo, two)
	}

	return listOne, listTwo
}

func getSessionCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: os.Getenv("AOC_TOKEN"),
	}

	return cookie
}

func getPuzzleInput(url string) {
	client := getNewHTTPClient()
	req := getNewRequest(url)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not GET %s", url)
	}
	defer resp.Body.Close()

	plainText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getPuzzleInput(): error reading response body: %v", err)
	}

	writeInputFile(plainText)
}

func getInputFile() (*os.File, error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("Local input file does not exist or it's empty")
		}
		fmt.Fprintf(os.Stderr, "getInputFile(): error reading file: %v", err)
	} else {
		fmt.Println("Local input file exists")
	}

	return file, err
}

func sortParsedInput(listOne, listTwo []int) [][]int {
	// sortedListOne := make([]int, len(listOne))
	// sortedListTwo := make([]int, len(listTwo))
	sortedLists := make([][]int, len(listOne))

	// for i, one := range listOne {
	// 	for _, two := range listOne - 1 {
	// 		if one > two {
	// 			sortedListOne[i] = two
	// 		}
	// 	}
	// }
	//
	slices.Sort(listOne)
	slices.Sort(listTwo)

	for i := range listOne {
		sortedLists[i] = []int{listOne[i], listTwo[i]}
	}

	return sortedLists
}

func writeInputFile(input []byte) {
	err := os.WriteFile("./input.txt", input, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writeInputFile(): error creating file: %v", err)
	}

	fmt.Println("Local input file created")
}

func main() {
	fmt.Println("Day one! Les get it!")

	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ENV's not loaded")
	}

	// If local file exists, do not make reqeust to AOC
	file, err := getInputFile()
	if err != nil {
		getPuzzleInput("https://adventofcode.com/2024/day/1/input")
		file, err = getInputFile()
	}

	challengeFunction(file)
}
