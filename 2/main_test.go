package main

import (
	"fmt"
	"slices"
	"testing"
)

type adjacentExpect struct {
	safe bool
}
type adjacentTest struct {
	dampened  bool
	dampening bool
	direction string
	expect    adjacentExpect
	reports   [][]int
}
type test struct {
	input  report
	output reports
}

func styleBAD(s string) string {
	return "\033[37;41m" + s + "\033[0m"
}

func TestDeleteSliceIndex(t *testing.T) {
	s := deleteSliceIndex([]int{1, 2, 3, 4, 5}, 0)
	if !slices.Equal(s, []int{2, 3, 4, 5}) {
		t.Errorf("Deleting slice index 0 (%d)", s)
	}

	s = deleteSliceIndex([]int{1, 2, 3, 4, 5}, 1)
	if !slices.Equal(s, []int{1, 3, 4, 5}) {
		t.Errorf("Deleting slice index 1 (%d)", s)
	}

	s = deleteSliceIndex([]int{1, 2, 3, 4, 5}, 2)
	if !slices.Equal(s, []int{1, 2, 4, 5}) {
		t.Errorf("Deleting slice index 2 (%d)", s)
	}

	s = deleteSliceIndex([]int{1, 2, 3, 4, 5}, 3)
	if !slices.Equal(s, []int{1, 2, 3, 5}) {
		t.Errorf("Deleting slice index 3 (%d)", s)
	}

	s = deleteSliceIndex([]int{1, 2, 3, 4, 5}, 4)
	if !slices.Equal(s, []int{1, 2, 3, 4}) {
		t.Errorf("Deleting slice index last (%d)", s)
	}

	s = deleteSliceIndex([]int{1, 2, 3, 4, 5}, len([]int{1, 2, 3, 4, 5})-1)
	if !slices.Equal(s, []int{1, 2, 3, 4}) {
		t.Errorf("Deleting slice index last (%d)", s)
	}
}

func getAdjacentLevelsTestResult(t *testing.T, cases adjacentTest) bool {
	pass := true

	for _, val := range cases.reports {
		safe := getReportAdjacentLevelsAcceptable(val, cases.direction, cases.dampening, cases.dampened)

		if safe != cases.expect.safe {
			pass = false
			t.Errorf("Safe (%t) is wrong", safe)
		}

		if !pass {
			t.Errorf("Incorrect report/case: %d", val)
		}
	}

	return pass
}

func TestAdjacentLevels(t *testing.T) {
	reports := [][]int{
		{5, 4, 3, 2, 1},
		{10, 8, 6, 4, 2},
		{15, 12, 9, 6, 3},
	}
	cases := adjacentTest{
		dampened:  false,
		dampening: false,
		direction: "decreasing",
		expect: adjacentExpect{
			safe: true,
		},
		reports: reports,
	}
	pass := getAdjacentLevelsTestResult(t, cases)
	if !pass {
		fmt.Printf("Decreasing, non-dampening, safe:" + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{5, 4, 3, 2, 1},
		{10, 8, 6, 4, 2},
		{15, 12, 9, 6, 3},
		{8, 4, 3, 2, 1},
	}
	cases = adjacentTest{
		dampened:  false,
		dampening: true,
		direction: "decreasing",
		expect: adjacentExpect{
			safe: true,
		},
		reports: reports,
	}
	pass = getAdjacentLevelsTestResult(t, cases)
	if !pass {
		fmt.Printf("Decreasing, dampening, dry, safe:" + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{1, 2, 3, 4, 5},
		{2, 4, 6, 8, 10},
		{3, 6, 9, 12, 15},
		{1, 2, 3, 4, 8},
	}
	cases = adjacentTest{
		dampened:  false,
		dampening: true,
		direction: "increasing",
		expect: adjacentExpect{
			safe: true,
		},
		reports: reports,
	}
	pass = getAdjacentLevelsTestResult(t, cases)
	if !pass {
		fmt.Printf("Decreasing, dampening, dry, safe:" + styleBAD("bad") + "\n")
	}
}

func TestGetDirectionallySafeReport(t *testing.T) {
	rawNonDampeningSafeIncreasingReports := []test{
		{input: report{1, 2, 3, 4, 5}, output: reports{{1, 2, 3, 4, 5}}},
		{input: report{2, 4, 6, 8, 10}, output: reports{{2, 4, 6, 8, 10}}},
		{input: report{3, 6, 9, 12, 15}, output: reports{{3, 6, 9, 12, 15}}},
		{input: report{100, 200, 300, 400, 500}, output: reports{{100, 200, 300, 400, 500}}},
	}
	for i, test := range rawNonDampeningSafeIncreasingReports {
		safe, report, _, dampened := getDirectionallySafeReport(test.input, false, "increasing")

		if !safe {
			t.Errorf("Input (%d) should be safe", test.input)
		}

		if dampened == true {
			t.Errorf("Dampened should not be true")
		}

		fmt.Println(i, report, test.output)
		if !slices.Equal(report, test.output[i]) {
			t.Errorf("Input's (%d) output does not match (%d); resulted in (%d)", test.input, test.output[i], report)
		}
	}
}

func TestGetSafeReports(t *testing.T) {
	rawNonDampeningSafeIncreasingReports := []test{
		{input: report{1, 2, 3, 4, 5}, output: reports{{1, 2, 3, 4, 5}}},
		{input: report{2, 4, 6, 8, 10}, output: reports{{2, 4, 6, 8, 10}}},
		{input: report{3, 6, 9, 12, 15}, output: reports{{3, 6, 9, 12, 15}}},
		// TODO: either update the output to be [] or swap getSafeReports call with
		// getDirectionallySafeReports becuase getSafeReports checks for adjacent acceptability
		{input: report{100, 200, 300, 400, 500}, output: reports{{100, 200, 300, 400, 500}}},
	}
	for _, val := range rawNonDampeningSafeIncreasingReports {
		nondampeningReports, dampeningReports := getSafeReports(val.input, false)

		fmt.Printf("Reports based on test input (%d): %d\n", val.input, nondampeningReports)

		if len(nondampeningReports) == 0 {
			t.Errorf("getSafeReports output 0 nondampening reports")
		}

		if len(dampeningReports) != 0 {
			t.Errorf("dampeningReports output reports; it shouldn't have")
		}

		for i, r := range nondampeningReports {
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output[i], r)
			}
		}
	}

	rawDampeningSafeIncreasingReports := []test{
		{input: report{10, 2, 3, 4, 5}, output: reports{{2, 3, 4, 5}}},
		{input: report{1, 10, 3, 4, 5}, output: reports{{1, 3, 4, 5}}},
		{input: report{1, 2, 10, 4, 5}, output: reports{{1, 2, 4, 5}}},
		{input: report{1, 2, 3, 10, 5}, output: reports{{1, 2, 3, 5}, {1, 2, 3, 10}}},
		{input: report{1, 2, 3, 4, 1}, output: reports{{1, 2, 3, 4}}},
	}
	for _, val := range rawDampeningSafeIncreasingReports {
		_, dampeningReports := getSafeReports(val.input, true)

		// fmt.Printf("Reports based on test input (%d): %d\n", val.input, reports)

		if len(dampeningReports) == 0 {
			t.Errorf("getSafeReports output 0 dampening reports")
		}

		for i, r := range dampeningReports {
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output, r)
			}
		}
	}

	rawNonDampeningSafeDecreasingReports := []test{
		{input: report{5, 4, 3, 2, 1}, output: reports{{5, 4, 3, 2, 1}}},
		{input: report{10, 8, 6, 4, 2}, output: reports{{10, 8, 6, 4, 2}}},
		{input: report{15, 12, 9, 6, 3}, output: reports{{15, 12, 9, 6, 3}}},
		{input: report{500, 400, 300, 200, 100}, output: reports{{500, 400, 300, 200, 100}}},
	}
	for _, val := range rawNonDampeningSafeDecreasingReports {
		nondampeningReports, dampeningReports := getSafeReports(val.input, false)

		// fmt.Printf("Reports based on test input (%d): %d\n", val.input, reports)

		if len(nondampeningReports) == 0 {
			t.Errorf("getSafeReports output 0 nondampening reports")
		}

		if len(dampeningReports) != 0 {
			t.Errorf("dampeningReports output reports; it shouldn't have")
		}

		for i, r := range nondampeningReports {
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output, r)
			}
		}
	}

	rawDampeningSafeDecreasingReports := []test{
		{input: report{5, 6, 3, 2, 1}, output: reports{{6, 3, 2, 1}, {5, 3, 2, 1}}},
		{input: report{5, 4, 6, 2, 1}, output: reports{{5, 4, 2, 1}}},
		{input: report{5, 4, 3, 6, 1}, output: reports{{5, 4, 3, 1}}},
		{input: report{5, 4, 3, 2, 6}, output: reports{{5, 4, 3, 2}}},
	}
	for _, val := range rawDampeningSafeDecreasingReports {
		_, dampeningReports := getSafeReports(val.input, true)

		// fmt.Printf("Reports based on test input (%d): %d\n", val.input, dampeningReports)

		if len(dampeningReports) == 0 {
			t.Errorf("getSafeReports output 0 reports")
		}

		for i, r := range dampeningReports {
			fmt.Println(r, val.output[i])
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output, r)
			}
		}
	}
}
