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

func handleReportConditionCheck(condition, dampened, dampening, safe, lookahead bool, input report, i int) (bool, bool, report, int, bool, bool) {
	// NOTE hrcc = handle report condition check
	hrccSafe := safe
	hrccDampened := dampened
	hrccI := i
	hrccReport := make(report, len(input))
	shouldBreak := false
	shouldContinue := false

	copy(hrccReport, input)

	if !condition {
		if dampened {
			// fmt.Println("Dampened")
			hrccSafe = false
			shouldBreak = true
		} else if dampening {
			// fmt.Println("Dampening")
			hrccDampened = true
			hrccReport = deleteSliceIndex(hrccReport, getIndexToDelete(i, lookahead))

			// By removing a value at index n, we are changing the values
			// checked in the condition above, so we want to reduce the index
			// by 1 to recheck the condition with the new value at index n
			if i != 0 {
				hrccI--
			}

			shouldContinue = true
		} else {
			// fmt.Println("Not dampening or dampened")

			hrccSafe = false
			shouldBreak = true
		}
	} else {
		// fmt.Println("Increasing index")

		hrccI++
	}

	return hrccSafe, hrccDampened, hrccReport, hrccI, shouldBreak, shouldContinue
}

func handleDuplicateReports(reports directionallySafeReports, report report) bool {
	// If identical report is already in safeReports, don't add it again
	// (consequence of looping through delete directions and have duplicate #'s at the last two indices)
	if len(reports) > 0 && slices.Equal(reports[0].report, report) {
		return true
	}

	return false
}

func getDirectionSafeReportUpdateSafeReports(safe bool, safeReports directionallySafeReports, report report, dampened bool) directionallySafeReports {
	reports := make(directionallySafeReports, len(safeReports))

	copy(reports, safeReports)

	if safe && !handleDuplicateReports(safeReports, report) {
		// fmt.Println("Safe", report, dampened)
		safeReport := directionallySafeReport{
			dampened: dampened,
			report:   report,
		}
		reports = append(safeReports, safeReport)
	}
	return reports
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

				hrccSafe, hrccDampened, hrccReport, hrccI, shouldBreak, shouldContinue := handleReportConditionCheck(condition, dampened, dampening, safe, lookahead, report, i)

				dampened = hrccDampened
				i = hrccI
				report = hrccReport
				safe = hrccSafe

				if shouldBreak {
					break
				}

				if shouldContinue {
					continue
				}

			}

			safeReports = getDirectionSafeReportUpdateSafeReports(safe, safeReports, report, dampened)
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
			// NOTE: getReportAdjacentLevelsAcceptable does not return modified reports
			// so we're adding reports that are directionally safe only because we only
			// care about the count of safe reports
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

func getReportAdjacentLevelsAcceptable(dsr directionallySafeReport, dampening bool) bool {
	dampened := dsr.dampened
	i := 0 // Artificially control loop index to psuedo recurse loop iteration
	report := make([]int, len(dsr.report))
	safe := true

	copy(report, dsr.report)

	for range len(report) - 1 {
		hrccSafe, hrccDampened, hrccReport, hrccI, shouldBreak, shouldContinue := handleReportConditionCheck(getAdjacentCondition(report, i), dampened, dampening, safe, false, report, i)

		dampened = hrccDampened
		i = hrccI
		report = hrccReport
		safe = hrccSafe

		if shouldBreak {
			break
		}

		if shouldContinue {
			continue
		}
	}

	return safe
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
