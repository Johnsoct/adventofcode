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

func calculateDistance(list [][]int) {
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

func calculateSimilarity(list [][]int) {
	similarityScore := 0

	for _, val := range list {
		left := val[0]
		occurences := 0

		for _, val2 := range list {
			right := val2[1]
			previousMatch := false

			if right == left {
				occurences += 1
				previousMatch = true
			}

			// The list is sorted, so if there was a match and no longer is
			// then there won't be any more matches
			if previousMatch && right != left {
				break
			}
		}

		similarityScore += left * occurences
	}

	fmt.Printf("Similarity score is: %d\n", similarityScore)
	fmt.Printf("The correct answer is 18567089\n")
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

	sortedLists := sortParsedInput(getParsedInput(file))

	calculateDistance(sortedLists)
	calculateSimilarity(sortedLists)
}
