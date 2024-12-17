// Solution for Advent of Code 2024 - Day one; puzzle one
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Johnsoct/adventofcode/get"
)

type directionallySafeReport struct {
	dampened bool
	report   report
}
type directionallySafeReports []directionallySafeReport
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

func getDirectionState(input report) (bool, int, bool, []int) {
	dampened := false
	i := 0 // Control index to psuedo recurse
	report := make(report, len(input))
	safe := true

	copy(report, input)

	return dampened, i, safe, report
}

func getIndexToDelete(i int, lookahead bool) int {
	if lookahead {
		return i + 1
	}

	return i
}

func getDirectionallySafeReport(input report, dampening bool) directionallySafeReports {
	deleteDirections := []string{"lookbehind", "lookahead"}
	directions := []string{"increasing", "decreasing"}
	safeReports := make(directionallySafeReports, 0)

	for _, direction := range directions {
		// fmt.Println("direction", direction)

		for _, deleteDirection := range deleteDirections {
			dampened, i, safe, report := getDirectionState(input)
			lookahead := false

			if deleteDirection == "lookahead" {
				lookahead = true
			}

			// fmt.Println("Delete direction", deleteDirection)

			for i < len(report)-1 {
				condition := getDirectionComparison(direction, report, i)

				// fmt.Println("Condition", condition)

				if !condition {
					if dampened {
						// fmt.Println("Dampened")
						safe = false
						break
					}

					if dampening {
						// fmt.Println("Dampening")
						dampened = true

						report = deleteSliceIndex(report, getIndexToDelete(i, lookahead))

						// By removing a value at index n, we are changing the values
						// checked in the condition above, so we want to reduce the index
						// by 1 to recheck the condition with the new value at index n
						if i != 0 {
							i--
						}

						continue
					}

					// fmt.Println("Not dampening or dampened")

					safe = false
					break
				}

				// fmt.Println("Increasing index")

				i++
			}

			if safe {
				// fmt.Println("Safe", report, dampened)

				// If identical report is already in safeReports, don't add it again
				// (consequence of looping through delete directions and have duplicate #'s at the last two indices)
				if len(safeReports) > 0 && slices.Equal(safeReports[0].report, report) {
					continue
				}

				safeReport := directionallySafeReport{
					dampened: dampened,
					report:   report,
				}
				safeReports = append(safeReports, safeReport)
			}
		}
	}

	return safeReports
}

func getSafeReports(report report, dampening bool) reports {
	safeReports := make(reports, 0)

	reports := getDirectionallySafeReport(report, dampening)

	for _, r := range reports {
		acceptable := getReportAdjacentLevelsAcceptable(r, dampening)
		if acceptable {
			// These aren't the dampened adjacent level report values as much
			// as a recognition that this report passed the adjacent levels
			// analysis with dampening
			safeReports = append(safeReports, r.report)
		}
	}

	return safeReports
}

func analyzeReports(input reports) {
	safeReports := make(reports, 0)
	unsafeReports := make(reports, 0)

	for _, val := range input {
		reports := getSafeReports(val, false)

		// By storing the unsafe reports, we can perform more efficient
		// dampening analysis by not running running against all inputs
		if len(reports) == 0 {
			unsafeReports = append(unsafeReports, val)
		} else {
			for _, r := range reports {
				safeReports = append(safeReports, r)
			}
		}
	}

	fmt.Printf("\nTotal unsafe reports: %d\n\n", len(unsafeReports))
	fmt.Printf("Total safe reports: %d\n", len(safeReports))
	fmt.Printf("The correct number of safe reports is 680\n\n")

	for _, val := range unsafeReports {
		reports := getSafeReports(val, true)

		for _, r := range reports {
			fmt.Println("original:", val, "dampened:", r)
			safeReports = append(safeReports, r)
		}
	}

	fmt.Printf("Total problem dampened reports: %d\n", len(safeReports))
	fmt.Printf("The correct number of problem dampened reports is 710\n\n")
}

func getAdjacentCondition(report []int, i int) bool {
	// NOTE: report is expected to be directionally safe
	var diff int
	current := report[i]
	max := 3
	min := 1
	next := report[i+1]

	if report[i] > report[i+1] {
		diff = current - next
	} else {
		diff = next - current
	}

	// fmt.Println("adjacent condition", i, report, "current", current, "next", next, "condition", diff < min || diff > max)

	return diff >= min && diff <= max
}

func getReportAdjacentLevelsAcceptable(report directionallySafeReport, dampening bool) bool {
	// NOTE: All reports passed as report are assumed directionally safe
	acceptable := true
	dampened := report.dampened
	i := 0 // Artificially control loop index to psuedo recurse loop iteration
	temp := make([]int, len(report.report))

	copy(temp, report.report)

	for range len(temp) - 1 {
		if !getAdjacentCondition(temp, i) {
			if dampened {
				// fmt.Println("Dampened")
				acceptable = false
				break
			}

			if dampening {
				// fmt.Println("Dampening")
				dampened = true
				temp = deleteSliceIndex(temp, getIndexToDelete(i, false))
				fmt.Println(temp)

				// By removing a value at index n, we are changing the values
				// checked in the condition above, so we want to reduce the index
				// by 1 to recheck the condition with the new value at index n
				if i != 0 {
					i--
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
