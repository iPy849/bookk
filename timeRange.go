package bookk

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	TimeRangeBoundsExclusion = 0b00 // Excludes both bounds of the range: (lower, upper)
	TimeRangeBoundsInclusion = 0b11 // Includes both bounds of the range: [lower, upper]
	TimeRangeIlEu            = 0b10 // Includes lower bound and excludes upper bound of the range: [lower, upper)
	TimeRangeElIu            = 0b01 // Excludes lower bound and includes upper bound of the range: (lower, upper]
)

var (
	timeRangeInitializatioDataError                = errors.New("Lower bound is greater than Upper bound in range")
	timeRangeInitializationNotRecognizedBoundError = errors.New("Bounds not recognized in range")
	timeRangeParseError                            = errors.New("Error parsing time range")
)

type TimeRangeBound byte

type MultiTimeRange []*TimeRange

func (mr MultiTimeRange) Len() int {
	return len(mr)
}

func (mr MultiTimeRange) Less(i, j int) bool {
	return mr[i].lowerBound.Unix() < mr[j].lowerBound.Unix()
}

func (mr MultiTimeRange) Swap(i, j int) {
	mr[i], mr[j] = mr[j], mr[i]
}

type TimeRange struct {
	lowerBound time.Time
	upperBound time.Time
	boundsConf TimeRangeBound
}

/*
NewTimeRange creates a new TimeRange with the specified bounds and inclusion/exclusion configuration.

This function constructs a TimeRange object with lower and upper time bounds and a configuration
that specifies whether each bound is inclusive or exclusive.

Parameters:
  - lower: The lower time bound of the TimeRange
  - upper: The upper time bound of the TimeRange
  - bounds: Configuration for inclusion/exclusion of bounds (use the predefined constants:
    TimeRangeBoundsExclusion, TimeRangeBoundsInclusion, TimeRangeIlEu, or TimeRangeElIu)

Returns:
  - A pointer to a new TimeRange object
  - An error if lower > upper or if the bounds configuration is invalid
*/
func NewTimeRange(lower, upper time.Time, bounds TimeRangeBound) (*TimeRange, error) {
	if lower.After(upper) {
		return nil, timeRangeInitializatioDataError
	} else if bounds > TimeRangeBoundsInclusion {
		return nil, timeRangeInitializationNotRecognizedBoundError
	}
	return &TimeRange{lower, upper, bounds}, nil
}

/*
TimeRangeFromPostgresString creates a TimeRange from a PostgreSQL range string format.

This function parses a string in PostgreSQL range notation and creates a corresponding
TimeRange object. The format should follow PostgreSQL conventions where square brackets
indicate inclusive bounds and parentheses indicate exclusive bounds.

Parameters:
  - r: A string in PostgreSQL range format, e.g., "[2025-01-01 00:00:00,2025-01-02 00:00:00)"

Returns:
  - A pointer to a new TimeRange object
  - An error if the string cannot be properly parsed
*/
func TimeRangeFromPostgresString(r string) (*TimeRange, error) {
	bounds := TimeRangeBound(0)
	if r[0] == '[' {
		bounds |= 0b10
	}
	if r[len(r)-1] == ']' {
		bounds |= 0b01
	}

	r = strings.ReplaceAll(r, "\"", "")
	dates := strings.Split(r[1:len(r)-1], ",")[:2]
	lowerTime, err := time.Parse(time.DateTime, dates[0])
	if err != nil {
		return nil, timeRangeParseError
	}

	upperTime, err := time.Parse(time.DateTime, dates[1])
	if err != nil {
		return nil, timeRangeParseError
	}

	return NewTimeRange(lowerTime, upperTime, bounds)
}

/*
ToPostgresRangeString converts the TimeRange to a PostgreSQL range string format.

This function formats the TimeRange into a string following PostgreSQL range notation
where square brackets indicate inclusive bounds and parentheses indicate exclusive bounds.

Returns:
  - A string representation of the time range in PostgreSQL format
*/
func (t *TimeRange) ToPostgresRangeString() string {
	var formatString string

	switch t.boundsConf {
	case TimeRangeBoundsExclusion:
		formatString = "(\"%s\",\"%s\")"
	case TimeRangeBoundsInclusion:
		formatString = "[\"%s\",\"%s\"]"
	case TimeRangeElIu:
		formatString = "(\"%s\",\"%s\"]"
	default:
		formatString = "[\"%s\",\"%s\")"
	}

	return fmt.Sprintf(
		formatString,
		t.lowerBound.Format(time.DateTime),
		t.upperBound.Format(time.DateTime),
	)
}

/*
Verbose returns a human-readable string representation of the time range.

This function creates a natural language description of the time range that varies
based on the inclusion/exclusion configuration of the bounds.

Returns:
  - A string describing the time range in natural language
*/
func (t *TimeRange) Verbose() string {
	var formatString string
	switch t.boundsConf {
	case TimeRangeBoundsExclusion:
		formatString = "Past %s and before %s"
	case TimeRangeBoundsInclusion:
		formatString = "From %s to %s"
	case TimeRangeElIu:
		formatString = "Past %s to %s"
	default:
		formatString = "From %s and before %s"
	}

	return fmt.Sprintf(
		formatString,
		t.lowerBound.Format(time.DateTime),
		t.upperBound.Format(time.DateTime),
	)
}

/*
Clone creates an independent copy of the TimeRange.

This function returns a new TimeRange object with the same lower and upper bounds
and the same inclusion/exclusion configuration as the original.

Returns:
  - A pointer to a copy of the oroginal TimeRange object
*/
func (t *TimeRange) Clone() *TimeRange {
	clone, _ := NewTimeRange(t.lowerBound, t.upperBound, t.boundsConf)
	return clone
}

func (t *TimeRange) lowerInclusion() bool {
	return t.boundsConf>>1 == 1
}

func (t *TimeRange) upperInclusion() bool {
	return t.boundsConf&0b01 == 1
}

/*
Equal determines whether two time ranges are identical.

This function checks if two TimeRange objects have the same lower bound, upper bound,
and inclusion/exclusion configuration. It correctly handles time comparison by using
time.Equal() for comparing Time objects.

Parameters:
  - r: The time range to compare with the current range
*/
func (t *TimeRange) Equal(r *TimeRange) bool {
	if t == r {
		return true
	}
	return t.lowerBound.Equal(r.lowerBound) && t.upperBound.Equal(r.upperBound) && t.boundsConf == r.boundsConf
}

/*
Contains checks if this TimeRange completely contains another TimeRange.

A TimeRange contains another when both the lower and upper bounds of the contained
TimeRange lies within the containing range. This function takes into account the
inclusion/exclusion configuration of both ranges' boundaries.

Parameters:
  - r: The TimeRange to check if it's contained within the current range
*/
func (t *TimeRange) Contains(r *TimeRange) bool {
	// Trivial cases
	if t.lowerBound.After(r.lowerBound) || t.upperBound.Before(r.upperBound) {
		return false
	}

	// Lower bound cases
	lowerBoundInclusion := (t.boundsConf>>1 >= r.boundsConf>>1)
	if t.lowerBound.Equal(r.lowerBound) && !lowerBoundInclusion {
		return false
	}

	// Upper bound cases
	upperBoundInclusion := (t.boundsConf&0b01 >= r.boundsConf&0b01)
	if t.upperBound.Equal(r.upperBound) && !upperBoundInclusion {
		return false
	}

	return true
}

/*
Overlaps determines how two TimeRange overlap with each other.

This function checks if the current TimeRange overlaps with another TimeRange
and indicates how they overlap. It does not consider exact boundary touches as overlaps;
for those cases, use the Union function instead.

Parameters:
  - r: The TimeRange to check for overlap with the current range

Returns:
  - 1: if this range overlaps on the upper bound of the other range
  - -1: if this range overlaps on the lower bound of the other range
  - 0: if there is no overlap or if the overlap is exactly at the boundaries
*/
func (t *TimeRange) Overlaps(r *TimeRange) int8 {
	if t.upperBound.After(r.lowerBound) && t.upperBound.Before(r.upperBound) {
		return 1
	} else if t.lowerBound.After(r.lowerBound) && t.upperBound.Before(r.upperBound) {
		return -1
	}
	return 0
}

/*
Union try combining two TimeRange into a single TimeRange if possible.

A union to is possible when ranges overlaps or exactly touch at their boundaries
with compatible inclusion/exclusion configurations. The boundary configuration of the
resulting TimeRange uses an OR combination of the original configurations, always producing the
most inclusive range possible.

Parameters:
  - r: The TimeRange to union with the current range

Returns:
  - A new TimeRange representing the union if the ranges can be merged
  - nil if the ranges are completely separate and cannot be merged
*/
func (t *TimeRange) Union(r *TimeRange) *TimeRange {
	// Equals
	if t.Equal(r) {
		return t.Clone()
	}

	// Contains itself
	if t.Contains(r) {
		return t.Clone()
	} else if r.Contains(t) {
		return r.Clone()
	}

	// Overlaps
	overlap := t.Overlaps(r)
	if overlap != 0 {
		var newRange *TimeRange
		var err error
		if overlap == 1 {
			newRange, err = NewTimeRange(t.lowerBound, r.upperBound, t.boundsConf|r.boundsConf)
		} else {
			newRange, err = NewTimeRange(r.lowerBound, t.upperBound, t.boundsConf|r.boundsConf)
		}
		if err != nil {
			return nil
		} else {
			return newRange
		}
	}

	// Boundaries overlaps
	if t.lowerBound.Equal(r.upperBound) && (t.lowerInclusion() || r.upperInclusion()) {
		newRange, err := NewTimeRange(r.lowerBound, t.upperBound, t.boundsConf|r.boundsConf)
		if err != nil {
			return nil
		}
		return newRange
	} else if t.upperBound.Equal(r.lowerBound) && (t.upperInclusion() || r.lowerInclusion()) {
		newRange, err := NewTimeRange(t.lowerBound, r.upperBound, t.boundsConf|r.boundsConf)
		if err != nil {
			return nil
		}
		return newRange
	}

	return nil
}

/*
Merge combines the current TimeRange with another TimeRange.

This function attempts to merge two TimeRange objects and returns a MultiTimeRange
representing the result. If the ranges can be unified (they overlap or touch at their
boundaries with compatible inclusion/exclusion configurations), the result will be
a MultiTimeRange with a single element. If they cannot be unified (they are completely
separate), the result will be a MultiTimeRange containing both ranges sorted by their
lower bounds.

Parameters:
  - r: The TimeRange to merge with the current range

Returns:
  - A MultiTimeRange containing either the unified range (if possible) or both
    ranges in chronological order
*/
func (t *TimeRange) Merge(r *TimeRange) MultiTimeRange {
	unionRange := t.Union(r)
	// No tienen puntos comunes
	if unionRange == nil {
		multiRanges := MultiTimeRange{t.Clone(), r.Clone()}
		sort.Sort(multiRanges)
		return multiRanges
	}
	return MultiTimeRange{unionRange.Clone()}
}
