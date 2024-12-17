// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Johnsoct/adventofcode/get"
)

func getParsedInput(rawInput io.Reader) []string {
	list := make([]string, 0)
	scanner := bufio.NewScanner(rawInput)

	for scanner.Scan() {
		line := scanner.Text()
		list = append(list, line)
	}

	return list
}

func main() {
	fmt.Println("Day three! Les get it!")

	get.GetEnv()

	// If local file exists, do not make reqeust to AOC
	file, err := get.GetInputFile()
	if err != nil {
		get.GetPuzzleInput("3")
		file, err = get.GetInputFile()
	}

	commands := getParsedInput(file)

	fmt.Println(commands)
}
