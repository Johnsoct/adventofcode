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

type report []int
type reports [][]int

func deleteSliceIndex[S ~[]E, E any](s S, i int) S {
	oldLen := len(s)
	tempI := i
	tempS := make(S, len(s))

	copy(tempS, s)

	if oldLen == i {
		return s
	}

	// fmt.Println(tempS[:tempI], tempS[tempI+1:])

	tempS = append(tempS[:tempI], tempS[tempI+1:]...)
	clear(tempS[len(tempS):oldLen]) // zero/nil out obsolete elements (GC)

	return tempS
}

func getDirectionComparison(direction string, report report, i int) bool {
	var condition bool

	if direction == "decreasing" {
		condition = getDirectionDecreasingComparison(report, i)
	} else {
		condition = getDirectionIncreasingComparison(report, i)
	}

	return condition
}

func getDirectionDecreasingComparison(report []int, i int) bool {
	return report[i] > report[i+1]
}

func getDirectionIncreasingComparison(report []int, i int) bool {
	return report[i] < report[i+1]
}

func getDirectionState(report []int) (bool, string, int, bool, []int) {
	dampened := false
	direction := "increasing" // default
	i := 0                    // Control index to psuedo recurse
	safe := true
	temp := make([]int, len(report))

	copy(temp, report)

	return dampened, direction, i, safe, temp
}

func getIndexToDelete(i int, lookahead bool) int {
	if lookahead {
		return i + 1
	}

	return i
}

func getDirectionallySafeReport(report report, dampening bool, direction string, lookahead bool) (bool, report, string, bool) {
	dampened, _, i, safe, r := getDirectionState(report)

	for i < len(r)-1 {
		condition := getDirectionComparison(direction, r, i)

		if !condition {
			if dampened {
				safe = false
				break
			}

			if dampening {
				dampened = true
				r = deleteSliceIndex(r, getIndexToDelete(i, lookahead))

				// By removing a value at index n, we are changing the values
				// checked in the condition above, so we want to reduce the index
				// by 1 to recheck the condition with the new value at index n
				if i != 0 {
					i--
				}

				continue
			}

			safe = false
			break
		}

		i++
	}

	return safe, r, direction, dampened
}

func getSafeReports(report report, dampening bool) (reports, reports) {
	dampeningReports := make(reports, 0)
	nondampeningReports := make(reports, 0)

	// Increasing looking behind
	safe, r, direction, dampened := getDirectionallySafeReport(report, dampening, "increasing", false)
	if safe {
		acceptable := getReportAdjacentLevelsAcceptable(r, direction, false, dampened)
		if acceptable {
			nondampeningReports = append(nondampeningReports, r)
		}
	}

	// Decreasing looking behind
	safe, r, direction, dampened = getDirectionallySafeReport(report, dampening, "decreasing", false)
	if safe {
		acceptable := getReportAdjacentLevelsAcceptable(r, direction, false, dampened)
		if acceptable {
			nondampeningReports = append(nondampeningReports, r)
		}
	}

	// If not dampening, checking lookahead and lookbehind results in duplicates
	if dampening {
		// Increasing looking ahead
		safe, r, direction, dampened = getDirectionallySafeReport(report, dampening, "increasing", true)
		if safe {
			acceptable := getReportAdjacentLevelsAcceptable(r, direction, true, dampened)
			if acceptable {
				dampeningReports = append(dampeningReports, r)
			}
		}

		// Decreasing looking ahead
		safe, r, direction, dampened = getDirectionallySafeReport(report, dampening, "decreasing", true)
		if safe {
			acceptable := getReportAdjacentLevelsAcceptable(r, direction, true, dampened)
			if acceptable {
				dampeningReports = append(dampeningReports, r)
			}
		}
	}

	return nondampeningReports, dampeningReports
}

func analyzeReports(reports reports) {
	problemDampenerSafeReportCount := 0
	safeReportCount := 0

	for _, val := range reports {
		nondampeningReports, dampeningReports := getSafeReports(val, true)
		problemDampenerSafeReportCount += len(dampeningReports)
		safeReportCount += len(nondampeningReports)
	}

	fmt.Printf("Total safe reports: %d\n", safeReportCount)
	fmt.Printf("The correct number of safe reports is 680\n")
	fmt.Printf("Total problem dampened reports: %d\n", problemDampenerSafeReportCount)
	fmt.Printf("The correct number of problem dampened reports is 680\n")
}

func getAdjacentCondition(report []int, i int, isDecreasing bool) bool {
	current := report[i]
	max := 3
	min := 1
	next := report[i+1]

	diff := next - current
	if isDecreasing {
		diff = current - next
	}

	// fmt.Println("adjacent condition", i, report, "current", current, "next", next, "condition", diff < min || diff > max)

	return diff >= min && diff <= max
}

func getReportAdjacentLevelsAcceptable(report report, direction string, dampening, damp bool) bool {
	// NOTE: All reports passed as report are assumed directionally safe
	acceptable := true
	dampened := damp
	i := 0 // Artificially control loop index to psuedo recurse loop iteration
	isDecreasing := true
	temp := make([]int, len(report))

	if direction == "increasing" {
		isDecreasing = false
	}

	copy(temp, report)

	for range len(temp) - 1 {
		if !getAdjacentCondition(temp, i, isDecreasing) {
			if dampening && !dampened {
				dampened = true
				indexToRemove := i

				if !isDecreasing {
					indexToRemove = i + 1
				}

				temp = deleteSliceIndex(temp, indexToRemove)

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
