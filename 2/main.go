// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Johnsoct/adventofcode/get"
)

func analyzeReports(reports [][]int) {
	problemDampenerSafeReportCount := 0
	safeReportCount := 0

	for _, val := range reports {
		if getReportSafety(val, false) {
			safeReportCount++
		}
                if (getReportSafety(val, true))
	}

	fmt.Printf("Total safe reports: %d\n", safeReportCount)
	fmt.Printf("The correct answer is 680\n")
}

func getReportAdjacentLevelsAcceptable(report []int, isDecreasing, isIncreasing bool) bool {
	acceptable := true
	min := 1
	max := 3

	for i := 0; i < len(report)-1; i++ {
		if isDecreasing {
			diff := report[i] - report[i+1]
			if diff < min || diff > max {
				acceptable = false
				break
			}
		}
		if isIncreasing {
			diff := report[i+1] - report[i]
			if diff < min || diff > max {
				acceptable = false
				break
			}
		}
	}

	return acceptable
}

func getReportIsDecreasing(report []int, problemDampener bool) (bool, bool) {
        dampened := false
	decreasing := true

	for i := 0; i < len(report)-1; i++ {
		if report[i+1] > report[i] {
                        if problemDampener {
                                if dampened == true {
                                        decreasing := false
                                        break
                                }
                                dampened = true
                                continue
                        } else {
                                decreasing = false
                                break
                        }
		}
	}

	return decreasing, dampened
}

func getReportIsIncreasing(report []int, problemDampener bool) (bool, bool) {
        dampened := false
	increasing := true

	for i := 0; i < len(report)-1; i++ {
		if report[i+1] < report[i] {
                        if problemDampener {
                                if dampened == true {
                                        increasing := false
                                        break
                                }
                                dampened = true
                                continue
                        } else {
                                increasing = false
                                break
                        }
		}
	}

	return increasing, dampened
}

func getReportSafety(report []int, problemDampener bool) bool {
        dampened := false
	isAdjacentLevelsAcceptable := false
	isDecreasing, damp := getReportIsDecreasing(report, problemDampener)
	isIncreasing, damp := getReportIsIncreasing(report, problemDampener)

	if isDecreasing || isIncreasing {
		isAdjacentLevelsAcceptable = getReportAdjacentLevelsAcceptable(report, isDecreasing, isIncreasing)
	}

	return (isDecreasing || isIncreasing) && isAdjacentLevelsAcceptable
}

func parseRawInput(file *os.File) [][]int {
	reports := make([][]int, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tuple := strings.Split(line, " ")
		convTuple := make([]int, len(tuple))

		for i, val := range tuple {
			v, err := strconv.Atoi(val)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error converting %s to integer", val)
			}

			convTuple[i] = v
		}

		reports = append(reports, convTuple)
	}

	return reports
}

func main() {
	fmt.Println("Day two! Les get it!")

	get.GetEnv()

	// If local file exists, do not make reqeust to AOC
	file, err := get.GetInputFile()
	if err != nil {
		get.GetPuzzleInput("2")
		file, err = get.GetInputFile()
	}

	parsedInput := parseRawInput(file)
	analyzeReports(parsedInput)
}
