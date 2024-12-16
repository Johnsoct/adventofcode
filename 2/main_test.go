package main

import (
	"fmt"
	"slices"
	"testing"
)

type adjacentTest struct {
	dampening bool
	reports   directionallySafeReports
}
type test struct {
	dampening bool
	input     report
	output    reports
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

func TestAdjacentLevels(t *testing.T) {
	reports := directionallySafeReports{
		{dampened: false, report: report{5, 4, 3, 2, 1}},
		{dampened: false, report: report{10, 8, 6, 4, 2}},
		{dampened: false, report: report{15, 12, 9, 6, 3}},
	}
	cases := adjacentTest{
		dampening: false,
		reports:   reports,
	}
	for _, test := range cases.reports {
		pass := getReportAdjacentLevelsAcceptable(test, cases.dampening)
		if !pass {
			fmt.Printf("Decreasing, non-dampening, safe:" + styleBAD("bad") + "\n")
		}
	}

	reports = directionallySafeReports{
		{dampened: false, report: report{5, 4, 3, 2, 1}},
		{dampened: false, report: report{10, 8, 6, 4, 2}},
		{dampened: false, report: report{15, 12, 9, 6, 3}},
		{dampened: false, report: report{8, 4, 3, 2, 1}},
	}
	cases = adjacentTest{
		dampening: true,
		reports:   reports,
	}
	for _, test := range cases.reports {
		pass := getReportAdjacentLevelsAcceptable(test, cases.dampening)
		if !pass {
			fmt.Printf("Decreasing, dampening, dry, safe:" + styleBAD("bad") + "\n")
		}
	}

	reports = directionallySafeReports{
		{dampened: false, report: report{1, 2, 3, 4, 5}},
		{dampened: false, report: report{2, 4, 6, 8, 10}},
		{dampened: false, report: report{3, 6, 9, 12, 15}},
		{dampened: false, report: report{1, 2, 3, 4, 8}},
		{dampened: false, report: report{19, 23, 26, 29, 31, 32}},
	}
	cases = adjacentTest{
		dampening: true,
		reports:   reports,
	}
	for _, test := range cases.reports {
		pass := getReportAdjacentLevelsAcceptable(test, cases.dampening)
		if !pass {
			fmt.Printf("Increasing, dampening, dry, safe:" + styleBAD("bad") + "\n")
		}
	}
}

func TestGetDirectionallySafeReport(t *testing.T) {
	tests := []test{
		{dampening: false, input: report{1, 2, 3, 4, 5}, output: reports{{1, 2, 3, 4, 5}}},
		{dampening: false, input: report{2, 4, 6, 8, 10}, output: reports{{2, 4, 6, 8, 10}}},
		{dampening: false, input: report{3, 6, 9, 12, 15}, output: reports{{3, 6, 9, 12, 15}}},
		{dampening: false, input: report{87, 90, 93, 94, 98}, output: reports{{87, 90, 93, 94, 98}}},
		{dampening: false, input: report{100, 200, 300, 400, 500}, output: reports{{100, 200, 300, 400, 500}}},
		{dampening: true, input: report{60, 64, 65, 69, 70, 73, 74, 74}, output: reports{{60, 64, 65, 69, 70, 73, 74}}},
		{dampening: true, input: report{37, 41, 43, 50, 52, 54, 52}, output: reports{{37, 41, 43, 50, 52, 54}}},
		{dampening: true, input: report{61, 59, 62, 63, 65, 66, 69, 72}, output: reports{{59, 62, 63, 65, 66, 69, 72}, {61, 62, 63, 65, 66, 69, 72}}},
		{dampening: true, input: report{53, 58, 65, 66, 63}, output: reports{{53, 58, 65, 66}}},
		{dampening: true, input: report{15, 19, 21, 26, 28, 28}, output: reports{{15, 19, 21, 26, 28}}},
		{dampening: true, input: report{8, 7, 11, 12, 17}, output: reports{{7, 11, 12, 17}, {8, 11, 12, 17}}},
		{dampening: true, input: report{51, 53, 54, 55, 57, 60, 63, 63}, output: reports{{51, 53, 54, 55, 57, 60, 63}}},
		{dampening: true, input: report{27, 29, 30, 33, 34, 35, 37, 35}, output: reports{{27, 29, 30, 33, 34, 35, 37}}},
		{dampening: true, input: report{10, 2, 3, 4, 5}, output: reports{{2, 3, 4, 5}}},
		{dampening: true, input: report{1, 10, 3, 4, 5}, output: reports{{1, 3, 4, 5}}},
		{dampening: true, input: report{1, 2, 10, 4, 5}, output: reports{{1, 2, 4, 5}}},
		{dampening: true, input: report{1, 2, 3, 10, 5}, output: reports{{1, 2, 3, 5}, {1, 2, 3, 10}}},
		{dampening: true, input: report{1, 2, 3, 4, 1}, output: reports{{1, 2, 3, 4}}},
		{dampening: false, input: report{93, 92, 91, 90, 88}, output: reports{{93, 92, 91, 90, 88}}},
		{dampening: false, input: report{5, 4, 3, 2, 1}, output: reports{{5, 4, 3, 2, 1}}},
		{dampening: false, input: report{10, 8, 6, 4, 2}, output: reports{{10, 8, 6, 4, 2}}},
		{dampening: false, input: report{15, 12, 9, 6, 3}, output: reports{{15, 12, 9, 6, 3}}},
		{dampening: false, input: report{500, 400, 300, 200, 100}, output: reports{{500, 400, 300, 200, 100}}},
		{dampening: true, input: report{76, 74, 73, 70, 64, 64}, output: reports{{76, 74, 73, 70, 64}}},
		{dampening: true, input: report{44, 44, 41, 40, 37, 31, 27}, output: reports{{44, 41, 40, 37, 31, 27}}},
		{dampening: true, input: report{96, 96, 93, 91, 88, 87, 81, 74}, output: reports{{96, 93, 91, 88, 87, 81, 74}}},
		{dampening: true, input: report{8, 6, 4, 4, 1}, output: reports{{8, 6, 4, 1}}},
		{dampening: true, input: report{5, 6, 3, 2, 1}, output: reports{{6, 3, 2, 1}, {5, 3, 2, 1}}},
		{dampening: true, input: report{5, 4, 6, 2, 1}, output: reports{{5, 4, 2, 1}}},
		{dampening: true, input: report{5, 4, 3, 6, 1}, output: reports{{5, 4, 3, 1}}},
		{dampening: true, input: report{5, 4, 3, 2, 6}, output: reports{{5, 4, 3, 2}}},
	}
	for _, test := range tests {
		reports := getDirectionallySafeReport(test.input, test.dampening)

		fmt.Printf("Reports based on test input (%d): %v\n", test.input, reports)

		if len(reports) == 0 {
			t.Errorf("getSafeReports output 0 nondampening reports")
		}

		for i, report := range reports {
			if report.dampened && !test.dampening {
				t.Errorf("Input (%d) should be not be dampened", test.input)
			}

			if !slices.Equal(report.report, test.output[i]) {
				t.Errorf("Input's (%d) output does not match (%d); resulted in (%d)", test.input, test.output[i], report.report)
			}
		}
	}
}
