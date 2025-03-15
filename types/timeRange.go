package types

import (
	"errors"
	"fmt"
	"time"
)

const (
	TIME_RANGE_BOUNDS_EXCLUSION = iota // Excludes both bounds of the range: (lower, upper)
	TIME_RANGE_BOUNDS_INCLUSION        // Includes both bounds of the range: [lower, upper]
	TIME_RANGE_IL_EU                   // Includes lower bound and excludes upper bound of the range: [lower, upper)
	TIME_RANGE_EL_IU                   // Excludes lower bound and includes upper bound of the range: (lower, upper]
)

var (
	timeRangeInitializatioDataError                = errors.New("Lower bound is greater than Upper bound in range")
	timeRangeInitializationNotRecognizedBoundError = errors.New("Bounds not recognized in range")
)

type TimeRange struct {
	lowerBound time.Time
	upperBound time.Time
	boundsConf int
}

func NewTimeRange(lower, upper time.Time, bounds int) (*TimeRange, error) {
	if lower.After(upper) {
		return nil, timeRangeInitializatioDataError
	} else if bounds > TIME_RANGE_EL_IU {
		return nil, timeRangeInitializationNotRecognizedBoundError
	}
	return &TimeRange{lower, upper, TIME_RANGE_IL_EU}, nil
}

func (t *TimeRange) ToRangeString() string {
	var formatString string

	switch t.boundsConf {
	case TIME_RANGE_BOUNDS_EXCLUSION:
		formatString = "(%s, %s)"
	case TIME_RANGE_BOUNDS_INCLUSION:
		formatString = "[%s, %s]"
	case TIME_RANGE_EL_IU:
		formatString = "(%s, %s]"
	default:
		formatString = "[%s, %s)"
	}

	return fmt.Sprintf(formatString, t.lowerBound.String(), t.lowerBound.String())
}
