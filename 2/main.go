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

func deleteSliceIndex[S ~[]E, E any](s S, i int) S {
	oldLen := len(s)
	tempI := i
	tempS := make(S, len(s))

	copy(tempS, s)

	if i == 0 {
		tempI = 1
	}
	if oldLen == i {
		return s
	}

	fmt.Println("old:", tempS)
	fmt.Println("index:", tempI, oldLen)
	fmt.Println("end", tempS[:tempI])
	fmt.Println("start", tempS[tempI+1:])
	tempS = append(tempS[:tempI], tempS[tempI+1:]...)
	fmt.Println("new:", tempS)
	fmt.Println()
	clear(tempS[len(tempS):oldLen]) // zero/nil out obsolete elements (GC)

	return tempS
}

func getDirectionDecreasingComparison(report []int, i int) bool {
	return report[i] > report[i+1]
}

func getDirectionIncreasingComparison(report []int, i int) bool {
	return report[i] < report[i+1]
}

func getDirection(report []int, dampening bool) (string, bool, []int, bool) {
	dampened := false
	direction := "increasing" // default
	i := 0                    // Control index to psuedo recurse
	safe := true
	temp := make([]int, len(report))

	copy(temp, report)

	// Logic:
	// In default direction,
	// 1.   Compare index 0 to index 1
	// 1a.  IF true, compare index 1 to index 2
	// 1aa. IF true, is default direction
	// 1ab. IF false, remove index 2, and if dampened != true, try 1a again
	// 1b.  Remove index 1, and if dampened != true, try 1 again
	//
	// If #1 does not return the default direction, repeat with the opposite direction

	for range len(temp) - 1 {
		if !getDirectionIncreasingComparison(temp, i) {
			if dampened {
				safe = false
				break
			}

			if dampening {
				dampened = true
				temp = deleteSliceIndex(temp, i)
				continue
			} else {
				safe = false
				break
			}
		}

		i++
	}

	// If report wasn't increasing...
	if !safe {
		// reset loop state
		i = 0
		safe = true
		temp = make([]int, len(report))

		copy(temp, report)

		for range len(temp) - 1 {
			if !getDirectionDecreasingComparison(temp, i) {
				if dampened {
					safe = false
					break
				}

				if dampening {
					dampened = true
					temp = deleteSliceIndex(temp, i)
					continue
				} else {
					safe = false
					break
				}
			}

			i++
		}

		if safe {
			direction = "decreasing"
		}
	}

	return direction, dampened, temp, safe
}

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
	// fmt.Println("adjacent condition", i, report, "current", current, "next", next, "condition", diff < min || diff > max)

	return diff < min || diff > max
}

func getReportAdjacentLevelsAcceptable(report []int, isDecreasing, problemDampener, damp bool) bool {
	acceptable := true
	dampened := damp
	i := 0 // Artificially control loop index to psuedo recurse loop iteration
	temp := make([]int, len(report))

	copy(temp, report)

	for i < len(temp)-1 {
		if getAdjacentCondition(temp, i, isDecreasing) {
			if problemDampener {
				a, r, d := ifProblemDampener(temp, i, dampened)
				// fmt.Println(a, r, d)

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
	tempReport := make([]int, len(report))

	copy(tempReport, report)

	isDecreasing, damp, dampenedReport := getReportSnowballing(tempReport, "decreasing", problemDampener)
	if damp == true {
		dampened = true
		tempReport = dampenedReport
	}
	isIncreasing, damp, dampenedReport := getReportSnowballing(tempReport, "increasing", problemDampener)
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
