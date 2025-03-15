package bookk

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func createValidRange() (*TimeRange, error) {
	lowerBound := time.Now()
	upperBound := lowerBound.Add(time.Hour)
	return NewTimeRange(lowerBound, upperBound, TIME_RANGE_IL_EU)
}

func TestTimeRangeSuccesfulInitialization(t *testing.T) {
	_, error := createValidRange()
	if error != nil {
		t.Errorf("Cannot initialize TimeRange. Throwed error: %s", error.Error())
	}
}

func TestTimeRangeWrongBoundDataInitialization(t *testing.T) {
	lowerBound := time.Now()
	upperBound := lowerBound.Add(-time.Hour)

	_, error := NewTimeRange(lowerBound, upperBound, TIME_RANGE_IL_EU)
	if error == nil {
		t.Errorf("Should have failed to initialize TimeRange due to: %s",
			timeRangeInitializatioDataError.Error(),
		)
	} else if !errors.Is(error, timeRangeInitializatioDataError) {
		t.Errorf("Should have failed to initialize TimeRange due to: %s. Instead: %s",
			timeRangeInitializatioDataError.Error(),
			error.Error(),
		)
	}
}

func TestTimeRangeNotRecognizedgBoundInitialization(t *testing.T) {
	lowerBound := time.Now()
	upperBound := lowerBound.Add(time.Hour)

	_, error := NewTimeRange(lowerBound, upperBound, 4)
	if error == nil {
		t.Errorf("Should have failed to initialize TimeRange due to: %s",
			timeRangeInitializationNotRecognizedBoundError.Error(),
		)
	} else if !errors.Is(error, timeRangeInitializationNotRecognizedBoundError) {
		t.Errorf("Should have failed to initialize TimeRange due to: %s. Instead: %s",
			timeRangeInitializationNotRecognizedBoundError.Error(),
			error.Error(),
		)
	}
}

func TestTimeRangeToValidRangeString(t *testing.T) {
	timeRange, _ := createValidRange()
	strRange := timeRange.ToRangeString()
	startCharacter := string(strRange[0])
	lastCharacter := string(strRange[len(strRange)-1])

	containsStartCharacters := strings.Contains("[(", startCharacter)
	containsLastCharacters := strings.Contains("])", lastCharacter)
	if !containsStartCharacters || !containsLastCharacters {
		t.Errorf(
			"TimeRange string conversion should starts with '[' or '(' and ends with with ']' or ')'. Instead: Starts with %v and ends with %v",
			string(strRange[0]),
			string(strRange[len(strRange)-1]),
		)
	}
}
