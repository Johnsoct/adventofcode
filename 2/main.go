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

func getAdjacentCondition(report []int, i int, isDecreasing bool) bool {
	current := report[i]
	min := 1
	max := 3
	next := report[i+1]

	diff := next - current
	if isDecreasing {
		diff = current - next
	}
	fmt.Println("adjacent condition", i, report, "current", current, "next", next, "condition", diff < min || diff > max)

	return diff < min || diff > max
}

func getReportAdjacentLevelsAcceptable(report []int, isDecreasing, problemDampener, damp bool) bool {
	acceptable := true
	dampened := damp
	i := 0 // Artificially control loop index to psuedo recurse loop iteration
	temp := make([]int, len(report))

	copy(temp, report)

	for i < len(report)-1 {
		if getAdjacentCondition(report, i, isDecreasing) {
			if problemDampener {
				a, r, d := ifProblemDampener(temp, i, dampened)
				fmt.Println(a, r, d)

				acceptable = a
				dampened = d
				temp = r

				if !acceptable {
					break
				}

				// Do not increase index; "recursion" with updated temp
				continue
			} else {
				acceptable = false
				break
			}
		}

		i++
	}

	return acceptable
}

func getSnowBallCheck(report []int, index int, direction string) bool {
	current := report[index]
	next := report[index+1]

	// fmt.Println("dampenedReport", index, report, "current", current, "previous", previous)

	snowballCheck := next > current
	if direction == "decreasing" {
		snowballCheck = next < current
	}

	return snowballCheck
}

func ifProblemDampener(report []int, i int, damp bool) (bool, []int, bool) {
	acceptable := true
	dampened := damp
	temp := make([]int, len(report))

	copy(temp, report)

	if dampened == true {
		// If temp has already had one level removed, it is no longer acceptable
		acceptable = false
	} else {
		dampened = true
		temp = append(temp[:i], temp[i+1:]...) // "Remove" this problem level; continue on
	}

	return acceptable, temp, dampened
}

func getReportSnowballing(report []int, direction string, problemDampener bool) (bool, bool, []int) {
	dampened := false
	dampenedReport := make([]int, len(report))
	i := 0 // Control the index for when we remove a level from dampenedReport (psuedo recursion)
	snowballing := true

	copy(dampenedReport, report)

	for i < len(dampenedReport)-1 {
		if !getSnowBallCheck(dampenedReport, i, direction) {
			if problemDampener {
				a, r, d := ifProblemDampener(dampenedReport, i, dampened)

				dampened = d
				dampenedReport = r
				snowballing = a

				if !snowballing {
					break
				}

				// Do not increase index; "recursion" with updated dampenedReport
				continue
			} else {
				snowballing = false
				break
			}
		}

		i++
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
		isAdjacentLevelsAcceptable = getReportAdjacentLevelsAcceptable(tempReport, isDecreasing, problemDampener, dampened)
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
