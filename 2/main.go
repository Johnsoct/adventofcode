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

func getDirectionDecreasingComparison(report []int, i int) bool {
	return report[i] > report[i+1]
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

func getDirectionIncreasingComparison(report []int, i int) bool {
	return report[i] < report[i+1]
}

func getDirection(report []int, dampening bool) (string, bool, []int, bool) {
	// Logic:
	// In default direction,
	// 1.   Compare index 0 to index 1
	// 1a.  IF true, compare index 1 to index 2
	// 1aa. IF true, is default direction
	// 1ab. IF false, remove index 2, and if dampened != true, try 1a again
	// 1b.  Remove index 1, and if dampened != true, try 1 again
	//
	// If #1 does not return the default direction, repeat with the opposite direction

	dampened, direction, i, safe, temp := getDirectionState(report)

	// TODO: refactor the two loops, which are identical except the getcomparison function
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
		dampened, direction, i, safe, temp = getDirectionState(report)

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

	if !safe {
		direction = ""
		temp = report
	}

	return direction, dampened, temp, safe
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

func getIndexToDelete(i int, lookahead bool) int {
	if lookahead {
		return i + 1
	}

	return i
}

func getDirectionallySafeReport(report report, dampening bool, direction string, lookahead bool) (bool, report) {
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

				if i == len(r)-1 && !lookahead {
					i--
				}

				continue
			}

			safe = false
			break
		}

		i++
	}

	return safe, r
}

func getDirectionallySafeReports(report report, dampening bool) reports {
	reports := make(reports, 0)

	// Increasing looking behind
	safe, r := getDirectionallySafeReport(report, dampening, "increasing", false)
	if safe {
		reports = append(reports, r)
	}

	// Decreasing looking behind
	safe, r = getDirectionallySafeReport(report, dampening, "decreasing", false)
	if safe {
		reports = append(reports, r)
	}

	// If not dampening, checking lookahead and lookbehind results in duplicates
	if dampening {
		// Increasing looking ahead
		safe, r = getDirectionallySafeReport(report, dampening, "increasing", true)
		if safe {
			reports = append(reports, r)
		}

		// Decreasing looking ahead
		safe, r = getDirectionallySafeReport(report, dampening, "decreasing", true)
		if safe {
			reports = append(reports, r)
		}
	}

	return reports
}

func analyzeReportSafety(report []int, dampening bool) bool {
	direction, dampened, r, safe := getDirection(report, dampening)
	if safe {
		safe = getReportAdjacentLevelsAcceptable(r, direction, dampening, dampened)
	} else {
		// TODO: remove
		fmt.Println(r, dampening, dampened, report)
	}

	return safe
}

func analyzeReports(reports [][]int) {
	problemDampenerSafeReportCount := 0
	safeReportCount := 0
	unsafe := 0

	for _, val := range reports {
		// if analyzeReportSafety(val, false) {
		// 	safeReportCount++
		// }
		if analyzeReportSafety(val, true) {
			problemDampenerSafeReportCount++
		} else {
			unsafe++
		}
	}

	fmt.Println("unsafe reports", unsafe)
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

func getReportAdjacentLevelsAcceptable(report []int, direction string, dampening, damp bool) bool {
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
