package main

import (
	"testing"
)

// False nondampened cases mean there is at least one level not decreasing/increasing
// from the previous level
var falseNonDampenedDecreasingCases = [][]int{
	{5, 4, 6, 2, 1},
	{10, 8, 9, 4, 2},
	{15, 12, 13, 6, 3},
}
var falseNonDampenedIncreasingCases = [][]int{
	{1, 2, 0, 4, 5},
	{2, 4, 3, 8, 10},
	{3, 6, 5, 12, 15},
}
var trueNonDampenedDecreasingCases = [][]int{
	{5, 4, 3, 2, 1},
	{10, 8, 6, 4, 2},
	{15, 12, 9, 6, 3},
}
var trueNonDampenedIncreasingCases = [][]int{
	{1, 2, 3, 4, 5},
	{2, 4, 6, 8, 10},
	{3, 6, 9, 12, 15},
}

// False dampened cases mean there are more than one problem levels
var falseDampenedDecreasingCases = [][]int{
	{5, 6, 3, 4, 1},
	{5, 4, 5, 2, 3},
	{15, 9, 12, 4, 6},
}
var falseDampenedIncreasingCases = [][]int{
	{1, 0, 3, 2, 5},
	{1, 2, 2, 4, 3},
	{1, 2, 3, 3, 3},
	{1, 2, 2, 2, 3},
}
var trueDampenedDecreasingCases = [][]int{
	{5, 6, 3, 2, 1},
	{10, 8, 11, 5, 4},
	{6, 12, 9, 6, 3},
}
var trueDampenedIncreasingCases = [][]int{
	{1, 3, 2, 4, 5},
	{2, 4, 6, 10, 8},
	{3, 6, 11, 9, 12},
	{0, 1, 4, 9, 7},
}

// func TestSnowballing(t *testing.T) {
// 	for _, val := range falseNonDampenedDecreasingCases {
// 		snowballing, _, _ := getReportSnowballing(val, "decreasing", false)
//
// 		if snowballing {
// 			t.Errorf("%d should not be snowballing", val)
// 		}
// 	}
//
// 	for _, val := range falseNonDampenedIncreasingCases {
// 		snowballing, _, _ := getReportSnowballing(val, "increasing", false)
//
// 		if snowballing {
// 			t.Errorf("%d should not be snowballing", val)
// 		}
// 	}
//
// 	for _, val := range trueNonDampenedDecreasingCases {
// 		snowballing, _, _ := getReportSnowballing(val, "decreasing", false)
//
// 		if !snowballing {
// 			t.Errorf("%d should be snowballing", val)
// 		}
// 	}
//
// 	for _, val := range trueNonDampenedIncreasingCases {
// 		snowballing, _, _ := getReportSnowballing(val, "increasing", false)
//
// 		if !snowballing {
// 			t.Errorf("%d should be snowballing", val)
// 		}
// 	}
//
// 	for _, val := range falseDampenedDecreasingCases {
// 		snowballing, dampened, dampenedReport := getReportSnowballing(val, "decreasing", true)
//
// 		if snowballing {
// 			t.Errorf("%d should not be snowballing", val)
// 		}
//
// 		if !dampened {
// 			t.Errorf("%d should be dampened", val)
// 		}
//
// 		if len(dampenedReport) >= len(val) {
// 			t.Errorf("Dampened report (%d, %d) should be shorter than report (%d, %d)", dampenedReport, len(dampenedReport), val, len(val))
// 		}
// 	}
//
// 	for _, val := range falseDampenedIncreasingCases {
// 		snowballing, dampened, dampenedReport := getReportSnowballing(val, "increasing", true)
//
// 		if snowballing {
// 			t.Errorf("%d should not be snowballing", val)
// 		}
//
// 		if !dampened {
// 			t.Errorf("%d should be dampened", val)
// 		}
//
// 		if len(dampenedReport) >= len(val) {
// 			t.Errorf("Dampened report (%d, %d) should be shorter than report (%d, %d)", dampenedReport, len(dampenedReport), val, len(val))
// 		}
// 	}
//
// 	for _, val := range trueDampenedDecreasingCases {
// 		snowballing, dampened, dampenedReport := getReportSnowballing(val, "decreasing", true)
//
// 		if !snowballing {
// 			t.Errorf("%d should be snowballing", val)
// 		}
//
// 		if !dampened {
// 			t.Errorf("%d should be dampened", val)
// 		}
//
// 		if len(dampenedReport) >= len(val) {
// 			t.Errorf("Dampened report (%d, %d) should be shorter than report (%d, %d)", dampenedReport, len(dampenedReport), val, len(val))
// 		}
// 	}
//
// 	for _, val := range trueDampenedIncreasingCases {
// 		snowballing, dampened, dampenedReport := getReportSnowballing(val, "increasing", true)
//
// 		if !snowballing {
// 			t.Errorf("%d should be snowballing", val)
// 		}
//
// 		if !dampened {
// 			t.Errorf("%d should be dampened", val)
// 		}
//
// 		if len(dampenedReport) >= len(val) {
// 			t.Errorf("Dampened report (%d, %d) should be shorter than report (%d, %d)", dampenedReport, len(dampenedReport), val, len(val))
// 		}
// 	}
// }

func TestAcceptableAdjacentLevels(t *testing.T) {
	// for _, val := range falseNonDampenedDecreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, true, false, false)
	//
	// 	if acceptable {
	// 		t.Errorf("%d should not be acceptable", val)
	// 	}
	// }
	//
	// for _, val := range falseNonDampenedIncreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, false, false, false)
	//
	// 	if acceptable {
	// 		t.Errorf("%d should not be acceptable", val)
	// 	}
	// }
	//
	// for _, val := range trueNonDampenedDecreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, true, false, false)
	//
	// 	if !acceptable {
	// 		t.Errorf("%d should be acceptable", val)
	// 	}
	// }
	//
	// for _, val := range trueNonDampenedIncreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, false, false, false)
	//
	// 	if !acceptable {
	// 		t.Errorf("%d should be acceptable", val)
	// 	}
	// }
	//
	// for _, val := range falseDampenedDecreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, true, true, false)
	//
	// 	if acceptable {
	// 		t.Errorf("%d should not be acceptable", val)
	// 	}
	// }
	//
	// for _, val := range falseDampenedIncreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, false, true, false)
	//
	// 	if acceptable {
	// 		t.Errorf("%d should not be acceptable", val)
	// 	}
	// }

	for _, val := range trueDampenedDecreasingCases {
		acceptable := getReportAdjacentLevelsAcceptable(val, true, true, false)

		if !acceptable {
			t.Errorf("%d should be acceptable", val)
		}
	}

	// for _, val := range trueDampenedIncreasingCases {
	// 	acceptable := getReportAdjacentLevelsAcceptable(val, false, true, false)
	//
	// 	if !acceptable {
	// 		t.Errorf("%d should be acceptable", val)
	// 	}
	// }
}
