// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/Johnsoct/adventofcode/get"
)

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

func main() {
	fmt.Println("Day [ENTER #]! Les get it!")

	get.GetEnv()

	// If local file exists, do not make reqeust to AOC
	file, err := get.GetInputFile()
	if err != nil {
		get.GetPuzzleInput("1")
		file, err = get.GetInputFile()
	}
}
