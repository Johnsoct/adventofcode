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
		if getReportSafety(val, true) {
			problemDampenerSafeReportCount++
		}
	}

	fmt.Printf("Total safe reports: %d\n", safeReportCount)
	fmt.Printf("Total problem dampened reports: %d\n", problemDampenerSafeReportCount)
	fmt.Printf("The correct number of safe reports is 680\n")
	fmt.Printf("The correct number of problem dampened reports is 680\n")
}

func getReportAdjacentLevelsAcceptable(report []int, isDecreasing, isIncreasing, problemDampener, dampened bool) bool {
	acceptable := true
	min := 1
	max := 3

	for i := 0; i < len(report)-1; i++ {
		if isDecreasing {
			diff := report[i] - report[i+1]
			if diff < min || diff > max {
				if problemDampener {
					if dampened {
						acceptable = false
						break
					}
					acceptable = false
					continue
				} else {
					acceptable = false
					break
				}
			}
		}
		if isIncreasing {
			diff := report[i+1] - report[i]
			if diff < min || diff > max {
				if problemDampener {
					if dampened {
						acceptable = false
						break
					}
					acceptable = false
					continue
				} else {
					acceptable = false
					break
				}
			}
		}
	}

	return acceptable
}

func getReportSnowballing(report []int, direction string, problemDampener bool) (bool, bool, []int) {
	dampened := false
	dampenedReport := make([]int, len(report))
	snowballing := true

	copy(dampenedReport, report)

	for i := 0; i < len(dampenedReport); i++ {
		if i == 0 {
			continue
		}

		current := dampenedReport[i]
		previous := dampenedReport[i-1]

		fmt.Println("dampenedReport", i, dampenedReport, "current", current, "previous", previous)

		snowballCheck := previous < current
		if direction == "decreasing" {
			snowballCheck = previous > current
		}

		if !snowballCheck {
			fmt.Println("check failed")
			if problemDampener {
				fmt.Println("problem dampener true")
				// If previous iteration set dampened to true
				if dampened == true {
					fmt.Println("dampened is already true")
					// After "removing" one problem level, still not snowballing
					snowballing = false
					break
				}
				fmt.Println("setting dampened to true")
				// "Remove" this problem level; continue on
				fmt.Println("current report", dampenedReport)
				dampenedReport = append(dampenedReport[:i], dampenedReport[i+1:]...)
				fmt.Println("updated report", dampenedReport)
				dampened = true
				continue
			}

			snowballing = false
			break
		}
	}

	return snowballing, dampened, dampenedReport
}

func getReportSafety(report []int, problemDampener bool) bool {
	dampened := false
	isAdjacentLevelsAcceptable := false
	tempReport := report

	isDecreasing, damp, dampenedReport := getReportSnowballing(report, "decreasing", problemDampener)
	if damp == true {
		dampened = true
		tempReport = dampenedReport
	}
	isIncreasing, damp, dampenedReport := getReportSnowballing(report, "increasing", problemDampener)
	if damp == true {
		dampened = true
		tempReport = dampenedReport
	}

	if isDecreasing || isIncreasing {
		isAdjacentLevelsAcceptable = getReportAdjacentLevelsAcceptable(tempReport, isDecreasing, isIncreasing, problemDampener, dampened)
		fmt.Printf("Dampened? %t\tOriginal: %d\tUpdated: %d\tAcceptable: %t\n", dampened, report, tempReport, isAdjacentLevelsAcceptable)
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
