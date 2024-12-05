// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"

        "github.com/Johnsoct/adventofcode/get"
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

func sortParsedInput(listOne, listTwo []int) [][]int {
	sortedLists := make([][]int, len(listOne))

	slices.Sort(listOne)
	slices.Sort(listTwo)

	for i := range listOne {
		sortedLists[i] = []int{listOne[i], listTwo[i]}
	}

	return sortedLists
}

func main() {
	fmt.Println("Day one! Les get it!")

        get.GetEnv()

	// If local file exists, do not make reqeust to AOC
	file, err := get.GetInputFile()
        if err != nil {
		get.GetPuzzleInput("1")
		file, err = get.GetInputFile()
	}

	challengeFunction(file)
}
