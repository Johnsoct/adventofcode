// Solution for Advent of Code 2024 - Day one; puzzle one
package get

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/joho/godotenv"
)

func GetInput(day string) (*os.File, error) {
	// If local file exists, do not make reqeust to AOC
	file, err := GetInputFile()
	if err != nil {
		GetPuzzleInput(day)
		file, err = GetInputFile()
	}

        return file, err
}

func GetNewHTTPClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	return client
}

func GetNewRequest(url string) *http.Request {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "formRequest(): error creating new request: %v", err)
	}

	r.AddCookie(GetSessionCookie())

	return r
}

func GetSessionCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: os.Getenv("AOC_TOKEN"),
	}

	return cookie
}

func GetPuzzleInput(day string) {
        url := "https://adventofcode.com/2024/day/" + day + "/input"
	client := GetNewHTTPClient()
	req := GetNewRequest(url)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not GET %s", url)
	}
	defer resp.Body.Close()

	plainText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getPuzzleInput(): error reading response body: %v", err)
	}

	WriteInputFile(plainText)
}

func GetInputFile() (*os.File, error) {
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

func GetEnv () {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ENV's not loaded")
	}
}

func WriteInputFile(input []byte) {
	err := os.WriteFile("./input.txt", input, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writeInputFile(): error creating file: %v", err)
	}

	fmt.Println("Local input file created")
}
