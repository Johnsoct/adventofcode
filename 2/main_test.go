package main

import (
	"fmt"
	"slices"
	"testing"
)

type expectation struct {
	direction string
	dampened  bool
	safe      bool
}
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
type directionTest struct {
	dampening bool
	expect    expectation
	reports   [][]int
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

func getDirectionTestResult(t *testing.T, cases directionTest) bool {
	pass := true

	for _, val := range cases.reports {
		direction, dampened, temp, safe := getDirection(val, cases.dampening)

		if direction != cases.expect.direction {
			pass = false
			t.Errorf("Direction (%s) is wrong", direction)
		}

		if dampened != cases.expect.dampened {
			pass = false
			t.Errorf("Dampened (%t) is wrong", dampened)
		}

		// fmt.Println(slices.Equal(temp, val), cases.expect.reportToMatch)
		if (safe && !dampened) || !safe {
			reportsShouldMatch := true
			if slices.Equal(temp, val) != reportsShouldMatch {
				pass = false
				t.Errorf("Updated report slice (%d) is wrong", temp)
			}
		}

		if safe != cases.expect.safe {
			pass = false
			t.Errorf("Safe (%t) is wrong", safe)
		}

		if !pass {
			t.Errorf("Incorrect report/case: %d, temp: %d", val, temp)
		}
	}

	return pass
}

func TestGetDirection(t *testing.T) {
	reports := [][]int{
		{5, 4, 3, 2, 1},
		{10, 8, 6, 4, 2},
		{15, 12, 9, 6, 3},
		{3, 2, 1},
		{3, 2},
		{7123, 66, 33, 11, 1},
	}
	cases := directionTest{
		dampening: false,
		expect: expectation{
			direction: "decreasing",
			dampened:  false,
			safe:      true,
		},
		reports: reports,
	}

	pass := getDirectionTestResult(t, cases)
	if !pass {
		fmt.Printf("Decreasing, non-dampened, safe: " + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{5, 5, 4, 3, 2, 1},
		{10, 8, 6, 8, 4, 2},
		{985, 733, 1013, 44, 3},
		{5, 4, 3, 2, 1, 1, 0},
	}
	cases = directionTest{
		dampening: true,
		expect: expectation{
			direction: "decreasing",
			dampened:  true,
			safe:      true,
		},
		reports: reports,
	}

	pass = getDirectionTestResult(t, cases)
	if !pass {
		fmt.Printf("Decreasing, dampened, safe: " + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{5, 5, 4, 3, 3, 2, 1},
		{10, 8, 6, 8, 4, 4, 2},
		{985, 733, 1013, 44, 3, 33},
		{5, 4, 3, 2, 1, 1, 0, 1},
	}
	cases = directionTest{
		dampening: true,
		expect: expectation{
			direction: "",
			dampened:  true,
			safe:      false,
		},
		reports: reports,
	}

	pass = getDirectionTestResult(t, cases)
	if !pass {
		fmt.Printf("Decreasing, dampened, unsafe: " + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{1, 2, 3, 4, 5},
		{2, 4, 6, 8, 10},
		{3, 6, 9, 12, 15},
		{100, 300, 700, 900, 300000},
	}
	cases = directionTest{
		dampening: false,
		expect: expectation{
			direction: "increasing",
			dampened:  false,
			safe:      true,
		},
		reports: reports,
	}

	pass = getDirectionTestResult(t, cases)
	if !pass {
		fmt.Printf("Increasing, non-dampened, safe: " + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{1, 2, 1, 4, 5},
		{2, 4, 6, 6, 10},
		{3, 6, 9, 32, 15},
		{100, 300, 200, 900, 300000},
		{1, 2, 2, 3, 4, 5},
		{14, 18, 21, 24, 21},
		{46, 44, 46, 49, 51, 53, 60, 64},
	}
	cases = directionTest{
		dampening: true,
		expect: expectation{
			direction: "increasing",
			dampened:  true,
			safe:      true,
		},
		reports: reports,
	}

	pass = getDirectionTestResult(t, cases)
	if !pass {
		fmt.Printf("Increasing, dampened, safe: " + styleBAD("bad") + "\n")
	}

	reports = [][]int{
		{1, 2, 1, 0, 5},
		{2, 4, 6, 6, 10, 3},
		{3, 6, 9, 32, 15, 2},
		{100, 100, 300, 200, 900, 300000},
		{14, 18, 21, 24, 21, 21},
	}
	cases = directionTest{
		dampening: true,
		expect: expectation{
			direction: "",
			dampened:  true,
			safe:      false,
		},
		reports: reports,
	}

	pass = getDirectionTestResult(t, cases)
	if !pass {
		fmt.Printf("Increasing, dampened, unsafe: " + styleBAD("bad") + "\n")
	}
}

func TestGetDirectionallySafeReports(t *testing.T) {
	type test struct {
		input  report
		output reports
	}
	rawNonDampeningSafeIncreasingReports := []test{
		{input: report{1, 2, 3, 4, 5}, output: reports{{1, 2, 3, 4, 5}}},
		{input: report{2, 4, 6, 8, 10}, output: reports{{2, 4, 6, 8, 10}}},
		{input: report{3, 6, 9, 12, 15}, output: reports{{3, 6, 9, 12, 15}}},
		{input: report{100, 200, 300, 400, 500}, output: reports{{100, 200, 300, 400, 500}}},
	}
	for _, val := range rawNonDampeningSafeIncreasingReports {
		reports := getDirectionallySafeReports(val.input, false)

		for i, r := range reports {
			// fmt.Printf("Report (%d) of type (%T) compared to output (%d)\n", r, r, val.output[i])
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output[i], r)
			}
		}

		// fmt.Println(reports)
	}

	rawDampeningSafeIncreasingReports := []test{
		{input: report{10, 2, 3, 4, 5}, output: reports{{2, 3, 4, 5}}},
		{input: report{1, 10, 3, 4, 5}, output: reports{{1, 3, 4, 5}}},
		{input: report{1, 2, 10, 4, 5}, output: reports{{1, 2, 4, 5}}},
		{input: report{1, 2, 3, 10, 5}, output: reports{{1, 2, 3, 5}, {1, 2, 3, 10}}},
		{input: report{1, 2, 3, 4, 1}, output: reports{{1, 2, 3, 4}}},
	}
	for _, val := range rawDampeningSafeIncreasingReports {
		reports := getDirectionallySafeReports(val.input, true)

		fmt.Printf("Reports based on test input (%d): %d\n", val.input, reports)

		for i, r := range reports {
			fmt.Printf("Report (%d) of type (%T) compared to output (%d)\n\n", r, r, val.output[i])
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
		reports := getDirectionallySafeReports(val.input, false)

		fmt.Printf("Reports based on test input (%d): %d\n", val.input, reports)

		for i, r := range reports {
			fmt.Printf("Report (%d) of type (%T) compared to output (%d)\n\n", r, r, val.output[i])
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output, r)
			}
		}
	}

	rawDampeningSafeDecreasingReports := []test{
		{input: report{5, 6, 3, 2, 1}, output: reports{{5, 3, 2, 1}}},
		{input: report{5, 4, 6, 2, 1}, output: reports{{5, 4, 2, 1}}},
		{input: report{5, 4, 3, 6, 1}, output: reports{{5, 4, 3, 1}}},
		{input: report{5, 4, 3, 2, 6}, output: reports{{5, 4, 3, 2}}},
	}
	for _, val := range rawDampeningSafeDecreasingReports {
		reports := getDirectionallySafeReports(val.input, false)

		fmt.Printf("Reports based on test input (%d): %d\n", val.input, reports)

		for i, r := range reports {
			fmt.Printf("Report (%d) of type (%T) compared to output (%d)\n\n", r, r, val.output[i])
			if !slices.Equal(r, val.output[i]) {
				t.Errorf("Input (%d) did not result in expected output (%d); resulted in (%d)", val.input, val.output, r)
			}
		}
	}
}
