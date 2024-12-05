// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"fmt"
	"os"

        "github.com/Johnsoct/adventofcode/get"
)

func challengeFunction(file *os.File) {
	// fmt.Printf("Total distance is: %d\n", distance)
        fmt.Printf("The correct answer is 2000468\n")
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

