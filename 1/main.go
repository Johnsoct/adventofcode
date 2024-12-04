// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/joho/godotenv"
)

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

func getSessionCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: os.Getenv("AOC_TOKEN"),
	}

	return cookie
}

func getPuzzleInput(url string) []byte {
	client := getNewHTTPClient()
	req := getNewRequest(url)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not GET %s", url)
	}
	defer resp.Body.Close()

	input, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getPuzzleInput(): error reading response body, %v", err)
	}

	return input
}

func writeInputFile(input []byte) {

}

func main() {
	fmt.Println("Day one! Les get it!")

	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ENV's not loaded")
	}

	getPuzzleInput("https://adventofcode.com/2024/day/1/input")
	// writeInputFile(input)
}
