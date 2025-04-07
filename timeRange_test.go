package bookk

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

var (
	testTime = time.Now()
)

func TestTimeRangeInitialization(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testBoundsCases := [4]TimeRangeBound{
			TimeRangeIlEu,
			TimeRangeElIu,
			TimeRangeBoundsExclusion,
			TimeRangeBoundsInclusion,
		}

		for _, testBound := range testBoundsCases {
			lowerBound := testTime
			upperBound := lowerBound.Add(time.Hour)
			_, err := NewTimeRange(lowerBound, upperBound, testBound)
			if err != nil {
				t.Errorf("Cannot initialize TimeRange. Throwed error: %s", err.Error())
			}
		}
	})

	t.Run("Failed due to bound data", func(t *testing.T) {
		lowerBound := testTime
		upperBound := lowerBound.Add(-time.Hour)

		_, error := NewTimeRange(lowerBound, upperBound, TimeRangeIlEu)
		if error == nil {
			t.Errorf("Should have failed to initialize TimeRange due to: %s",
				timeRangeInitializatioDataError.Error(),
			)
		} else if !errors.Is(error, timeRangeInitializatioDataError) {
			t.Errorf("Should have failed to initialize TimeRange due to: %s. \nInstead: %s",
				timeRangeInitializatioDataError.Error(),
				error.Error(),
			)
		}
	})

	t.Run("Failed due to not recognized bound", func(t *testing.T) {
		lowerBound := testTime
		upperBound := lowerBound.Add(time.Hour)

		_, error := NewTimeRange(lowerBound, upperBound, 0b100)
		if error == nil {
			t.Errorf("Should have failed to initialize TimeRange due to: %s",
				timeRangeInitializationNotRecognizedBoundError.Error(),
			)
		} else if !errors.Is(error, timeRangeInitializationNotRecognizedBoundError) {
			t.Errorf("Should have failed to initialize TimeRange due to: %s. \nInstead: %s",
				timeRangeInitializationNotRecognizedBoundError.Error(),
				error.Error(),
			)
		}
	})
}

func TestTimeRangeToPostgresRangeString(t *testing.T) {
	timeRange, _ := NewTimeRange(testTime, testTime.Add(time.Hour), TimeRangeIlEu)
	strRange := timeRange.ToPostgresRangeString()

	containsStartCharacters := strings.Contains("[(", string(strRange[0]))
	containsLastCharacters := strings.Contains("])", string(strRange[len(strRange)-1]))
	if !containsStartCharacters || !containsLastCharacters {
		t.Errorf(
			"TimeRange string conversion should starts with '[' or '(' and ends with with ']' or ')'. Instead: Starts with %v and ends with %v",
			string(strRange[0]),
			string(strRange[len(strRange)-1]),
		)
	}

	quotesQty := len(strings.Split(strRange, "\""))
	if quotesQty != 5 {
		t.Errorf(
			"TimeRange string conversion is missing %d \" characters",
			5-quotesQty,
		)
	}
}

func TestTimeRangeFromPostgresRangeString(t *testing.T) {
	testCases := [4]string{
		"[\"2024-10-13 10:00:00\",\"2024-10-13 15:00:00\")",
		"(\"2024-10-13 10:00:00\",\"2025-10-13 15:00:00\")",
		"[\"2024-10-13 10:00:00\",\"2026-10-13 15:00:00\"]",
		"(\"2024-10-13 10:00:00\",\"2027-10-13 15:00:00\"]",
	}

	for _, testCase := range testCases {
		timeRange, err := TimeRangeFromPostgresString(testCase)
		if err != nil {
			t.Errorf("There is an error in TimeRange string parser. Parser Error: %s", err.Error())
		}

		timeRangeStringRepr := timeRange.ToPostgresRangeString()
		if timeRangeStringRepr != testCase {
			t.Errorf(
				"No expected output:\nExpecting\t: %s\nRecieved\t: %s",
				testCase,
				timeRangeStringRepr,
			)
			return
		}
	}
}

func TestTimeRangeContains(t *testing.T) {
	containerIlEu, _ := NewTimeRange(
		testTime,
		testTime.Add(time.Hour),
		TimeRangeIlEu,
	)

	containerElIu, _ := NewTimeRange(
		testTime,
		testTime.Add(time.Hour),
		TimeRangeElIu,
	)

	containerBe, _ := NewTimeRange(
		testTime,
		testTime.Add(time.Hour),
		TimeRangeBoundsExclusion,
	)

	containerBi, _ := NewTimeRange(
		testTime,
		testTime.Add(time.Hour),
		TimeRangeBoundsInclusion,
	)

	t.Cleanup(func() {
		containerIlEu = nil
		containerElIu = nil
		containerBe = nil
		containerBi = nil
	})

	containsErrorMessage := "It must be contained\nContainer: %s\nContained: %s"
	notContainsErrorMessage := "It must NOT be contained\nContainer: %s\nContained: %s"

	t.Run("Other TimeRange", func(t *testing.T) {
		testerRange, _ := NewTimeRange(
			testTime.Add(time.Minute*10),
			testTime.Add(time.Minute*50),
			TimeRangeIlEu,
		)

		if !containerIlEu.Contains(testerRange) {
			t.Errorf(
				containsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				testerRange.ToPostgresRangeString(),
			)
		}

		if !containerElIu.Contains(testerRange) {
			t.Errorf(
				containsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				testerRange.ToPostgresRangeString(),
			)
		}
	})

	t.Run("Not contains other TimeRange", func(t *testing.T) {
		testerWrongLowerBound, _ := NewTimeRange(
			testTime.Add(-time.Minute*10),
			testTime.Add(time.Minute*50),
			TimeRangeIlEu,
		)

		testerWrongUpperBound, _ := NewTimeRange(
			testTime.Add(time.Minute*10),
			testTime.Add(time.Minute*70),
			TimeRangeIlEu,
		)

		if containerIlEu.Contains(testerWrongLowerBound) {
			t.Errorf(
				notContainsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				testerWrongLowerBound.ToPostgresRangeString(),
			)
		}

		if containerElIu.Contains(testerWrongLowerBound) {
			t.Errorf(
				notContainsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				testerWrongLowerBound.ToPostgresRangeString(),
			)
		}

		if containerIlEu.Contains(testerWrongUpperBound) {
			t.Errorf(
				notContainsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				testerWrongUpperBound.ToPostgresRangeString(),
			)
		}

		if containerIlEu.Contains(testerWrongUpperBound) {
			t.Errorf(
				notContainsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				testerWrongUpperBound.ToPostgresRangeString(),
			)
		}
	})

	t.Run("Not contains other TimeRange outside of range", func(t *testing.T) {
		testerBeforeRange, _ := NewTimeRange(
			testTime.Add(-time.Minute*20),
			testTime.Add(-time.Minute*10),
			TimeRangeIlEu,
		)

		testerAfterRange, _ := NewTimeRange(
			testTime.Add(time.Minute*70),
			testTime.Add(time.Minute*80),
			TimeRangeIlEu,
		)

		if containerIlEu.Contains(testerBeforeRange) {
			t.Errorf(
				notContainsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				testerBeforeRange.ToPostgresRangeString(),
			)
		}

		if containerElIu.Contains(testerAfterRange) {
			t.Errorf(
				notContainsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				testerAfterRange.ToPostgresRangeString(),
			)
		}

		if containerIlEu.Contains(testerBeforeRange) {
			t.Errorf(
				notContainsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				testerBeforeRange.ToPostgresRangeString(),
			)
		}

		if containerIlEu.Contains(testerAfterRange) {
			t.Errorf(
				notContainsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				testerAfterRange.ToPostgresRangeString(),
			)
		}
	})

	t.Run("Itself", func(t *testing.T) {

		// Testing lower bound
		if !containerIlEu.Contains(containerIlEu) {
			t.Errorf(
				containsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				containerIlEu.ToPostgresRangeString(),
			)
		}

		if !containerElIu.Contains(containerElIu) {
			t.Errorf(
				containsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				containerElIu.ToPostgresRangeString(),
			)
		}

		if !containerBe.Contains(containerBe) {
			t.Errorf(
				containsErrorMessage,
				containerBe.ToPostgresRangeString(),
				containerBe.ToPostgresRangeString(),
			)
		}

		if !containerBi.Contains(containerBi) {
			t.Errorf(
				containsErrorMessage,
				containerBi.ToPostgresRangeString(),
				containerBi.ToPostgresRangeString(),
			)
		}
	})

	t.Run("Not contains itself with different bounds", func(t *testing.T) {
		if containerElIu.Contains(containerIlEu) {
			t.Errorf(
				notContainsErrorMessage,
				containerElIu.ToPostgresRangeString(),
				containerIlEu.ToPostgresRangeString(),
			)
		}

		if containerIlEu.Contains(containerElIu) {
			t.Errorf(
				notContainsErrorMessage,
				containerIlEu.ToPostgresRangeString(),
				containerElIu.ToPostgresRangeString(),
			)
		}
	})
}

func TestTimeRangeClone(t *testing.T) {

	type TestCase struct {
		name   string
		bounds TimeRangeBound
	}

	testCases := [4]TestCase{
		{"Inclusive-Exclusive", TimeRangeIlEu},
		{"Exclusive-Inclusive", TimeRangeElIu},
		{"Exclusive-Exclusive", TimeRangeBoundsExclusion},
		{"Inclusive-Inclusive", TimeRangeBoundsInclusion},
	}

	for _, testCase := range testCases {
		testName := fmt.Sprintf("Preserves all properties: %s", testCase.name)
		t.Run(testName, func(t *testing.T) {
			original, _ := NewTimeRange(testTime, testTime.Add(time.Hour), testCase.bounds)
			clone := original.Clone()

			if original == clone {
				t.Errorf("Clone returned the same object instance, not a copy")
			}

			if !original.lowerBound.Equal(clone.lowerBound) {
				t.Errorf("Clone lower bound mismatch:\nOriginal %v,\nclone %v",
					original.lowerBound, clone.lowerBound)
			}

			if !original.upperBound.Equal(clone.upperBound) {
				t.Errorf("Clone upper bound mismatch:\nOriginal %v,\bClone %v",
					original.upperBound, clone.upperBound)
			}

			if original.boundsConf != clone.boundsConf {
				t.Errorf("Clone bounds configuration mismatch:\nOriginal %v,\bClone %v",
					original.boundsConf, clone.boundsConf)
			}
		})
	}
}
